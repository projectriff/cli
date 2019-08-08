/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package v1alpha1

import (
	v1alpha1 "github.com/projectriff/system/pkg/apis/knative/v1alpha1"
	"github.com/projectriff/system/pkg/client/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type KnativeV1alpha1Interface interface {
	RESTClient() rest.Interface
	AdaptersGetter
	ConfigurersGetter
}

// KnativeV1alpha1Client is used to interact with features provided by the knative.projectriff.io group.
type KnativeV1alpha1Client struct {
	restClient rest.Interface
}

func (c *KnativeV1alpha1Client) Adapters(namespace string) AdapterInterface {
	return newAdapters(c, namespace)
}

func (c *KnativeV1alpha1Client) Configurers(namespace string) ConfigurerInterface {
	return newConfigurers(c, namespace)
}

// NewForConfig creates a new KnativeV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*KnativeV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &KnativeV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new KnativeV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *KnativeV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new KnativeV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *KnativeV1alpha1Client {
	return &KnativeV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *KnativeV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
