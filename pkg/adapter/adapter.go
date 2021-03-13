package adapter

import (
	"context"
	"fmt"
	nethttp "net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/uuid"

	"github.com/Azure/go-amqp"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"go.uber.org/zap"

	sourcesv1alpha1 "github.com/Markus-Schwer/eventing-amqp/pkg/apis/sources/v1alpha1"

	"knative.dev/eventing/pkg/adapter/v2"

	"knative.dev/eventing/pkg/kncloudevents"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/source"
)

const (
	resourceGroup = "amqpsources.sources.knative.dev"
)

type adapterConfig struct {
	adapter.EnvConfig

	Broker   string `envconfig:"AMQP_BROKER" required:"true"`
	Topic    string `envconfig:"AMQP_TOPIC" required:"true"`
	User     string `envconfig:"AMQP_USER" required:"false"`
	Password string `envconfig:"AMQP_PASSWORD" required:"false"`
}

func NewEnvConfig() adapter.EnvConfigAccessor {
	return &adapterConfig{}
}

// Adapter generates events at a regular interval.
type Adapter struct {
	config            *adapterConfig
	httpMessageSender *kncloudevents.HTTPMessageSender
	reporter          source.StatsReporter
	logger            *zap.Logger
	context           context.Context
}

var _ adapter.MessageAdapter = (*Adapter)(nil)
var _ adapter.MessageAdapterConstructor = NewAdapter

func NewAdapter(ctx context.Context, processed adapter.EnvConfigAccessor, httpMessageSender *kncloudevents.HTTPMessageSender, reporter source.StatsReporter) adapter.MessageAdapter {
	logger := logging.FromContext(ctx).Desugar()
	config := processed.(*adapterConfig)

	return &Adapter{
		config:            config,
		httpMessageSender: httpMessageSender,
		reporter:          reporter,
		logger:            logger,
		context:           ctx,
	}
}

func (a *Adapter) CreateClient(User string, Password string, logger *zap.Logger) (*amqp.Client, error) {
	addr := fmt.Sprintf("amqp://%s", a.config.Broker)

	var auth amqp.ConnOption
	if User != "" && Password != "" {
		auth = amqp.ConnSASLPlain(a.config.User, a.config.Password)
	} else {
		auth = amqp.ConnSASLAnonymous()
	}

	// TODO: support multiple auth methods
	client, err := amqp.Dial(addr, auth)
	if err != nil {
		logger.Error(err.Error())
	}

	return client, err
}

func (a *Adapter) CreateReceiver(session *amqp.Session, logger *zap.Logger) (*amqp.Receiver, error) {
	receiver, err := session.NewReceiver(amqp.LinkSourceAddress("/" + a.config.Topic))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return receiver, err
}

func (a *Adapter) Start(ctx context.Context) error {
	logger := a.logger

	logger.Info("Starting with config: ",
		zap.String("SinkURI", a.config.Sink),
		zap.String("Name", a.config.Name),
		zap.String("Namespace", a.config.Namespace))

	client, err := a.CreateClient(a.config.User, a.config.Password, logger)
	if err == nil {
		defer client.Close()
	}

	session, err := client.NewSession()
	if err != nil {
		logger.Error(err.Error())
	}

	receiver, err := a.CreateReceiver(session, logger)
	if err == nil {
		defer func() {
			closeTimetoutCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			receiver.Close(closeTimetoutCtx)
			cancel()
		}()
	}

	// queue, err := a.StartAmqpClient(&ch)
	// if err != nil {
	// 	logger.Error(err.Error())
	// }

	return a.PollForMessages(ctx, receiver, logger, ctx.Done())
}

func (a *Adapter) ConsumeMessage(msg *amqp.Message) error {
	return nil
}

func (a *Adapter) PollForMessages(ctx context.Context, receiver *amqp.Receiver, logger *zap.Logger, stopCh <-chan struct{}) error {
	for {
		err := receiver.HandleMessage(ctx, func(msg *amqp.Message) error {
			logger.Info("Received message: %s", zap.Any("value", msg.GetData()))

			if err := a.postMessage(msg); err == nil {
				logger.Info("Successfully sent event to sink")
				err = msg.Accept(ctx)
				if err != nil {
					logger.Error("Sending Accept failed")
				}
			} else {
				logger.Error("Sending event to sink failed: ", zap.Error(err))
				err = msg.Release(ctx)
				if err != nil {
					logger.Error("Sending Release failed")
				}
			}

			return nil
		})
		if err != nil {
			logger.Error("Reading message from AMQP:", zap.Error(err))
		}

		select {
		case <-stopCh:
			logger.Info("Shutting down...")
			return nil
		}
	}
}

func (a *Adapter) postMessage(msg *amqp.Message) error {
	a.logger.Info("url ->" + a.httpMessageSender.Target)
	req, err := a.httpMessageSender.NewCloudEventRequest(a.context)
	if err != nil {
		return err
	}

	// TODO: See if the message is already a Cloud Event and if so, do not wrap, just use as is.
	event := cloudevents.NewEvent()
	if msg.Properties.MessageID != nil {
		event.SetID(string(msg.Properties.MessageID.(string)))
	} else {
		event.SetID(string(uuid.NewUUID()))
	}
	event.SetTime(msg.Properties.CreationTime)
	event.SetType(sourcesv1alpha1.AmqpEventType)
	event.SetSource(sourcesv1alpha1.AmqpEventSource(a.config.Namespace, a.config.Name, a.config.Topic))
	event.SetSubject(msg.Properties.Subject)
	event.SetExtension("key", msg.Properties.MessageID)

	// TODO: Check the content type and use it instead.
	err = event.SetData(msg.Properties.ContentType, msg.GetData())
	if err != nil {
		return err
	}

	err = http.WriteRequest(a.context, binding.ToMessage(&event), req)
	if err != nil {
		return err
	}

	res, err := a.httpMessageSender.Send(req)

	if err != nil {
		a.logger.Debug("Error while sending the message", zap.Error(err))
		return err
	}

	if res.StatusCode/100 != 2 {
		a.logger.Debug("Unexpected status code", zap.Int("status code", res.StatusCode))
		return fmt.Errorf("%d %s", res.StatusCode, nethttp.StatusText(res.StatusCode))
	}

	reportArgs := &source.ReportArgs{
		Namespace:     a.config.Namespace,
		Name:          a.config.Name,
		ResourceGroup: resourceGroup,
	}

	_ = a.reporter.ReportEventCount(reportArgs, res.StatusCode)
	return nil
}
