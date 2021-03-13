/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/Markus-Schwer/eventing-amqp/pkg/apis/sources/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// AmqpSourceLister helps list AmqpSources.
// All objects returned here must be treated as read-only.
type AmqpSourceLister interface {
	// List lists all AmqpSources in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.AmqpSource, err error)
	// AmqpSources returns an object that can list and get AmqpSources.
	AmqpSources(namespace string) AmqpSourceNamespaceLister
	AmqpSourceListerExpansion
}

// amqpSourceLister implements the AmqpSourceLister interface.
type amqpSourceLister struct {
	indexer cache.Indexer
}

// NewAmqpSourceLister returns a new AmqpSourceLister.
func NewAmqpSourceLister(indexer cache.Indexer) AmqpSourceLister {
	return &amqpSourceLister{indexer: indexer}
}

// List lists all AmqpSources in the indexer.
func (s *amqpSourceLister) List(selector labels.Selector) (ret []*v1alpha1.AmqpSource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.AmqpSource))
	})
	return ret, err
}

// AmqpSources returns an object that can list and get AmqpSources.
func (s *amqpSourceLister) AmqpSources(namespace string) AmqpSourceNamespaceLister {
	return amqpSourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// AmqpSourceNamespaceLister helps list and get AmqpSources.
// All objects returned here must be treated as read-only.
type AmqpSourceNamespaceLister interface {
	// List lists all AmqpSources in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.AmqpSource, err error)
	// Get retrieves the AmqpSource from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.AmqpSource, error)
	AmqpSourceNamespaceListerExpansion
}

// amqpSourceNamespaceLister implements the AmqpSourceNamespaceLister
// interface.
type amqpSourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all AmqpSources in the indexer for a given namespace.
func (s amqpSourceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.AmqpSource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.AmqpSource))
	})
	return ret, err
}

// Get retrieves the AmqpSource from the indexer for a given namespace and name.
func (s amqpSourceNamespaceLister) Get(name string) (*v1alpha1.AmqpSource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("amqpsource"), name)
	}
	return obj.(*v1alpha1.AmqpSource), nil
}
