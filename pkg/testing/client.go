/*
 * Copyright 2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package testing

import (
	kntesting "github.com/knative/pkg/reconciler/testing"
	projectriffclientset "github.com/projectriff/system/pkg/client/clientset/versioned/fake"
	buildv1alpha1clientset "github.com/projectriff/system/pkg/client/clientset/versioned/typed/build/v1alpha1"
	requestv1alpha1clientset "github.com/projectriff/system/pkg/client/clientset/versioned/typed/request/v1alpha1"
	streamv1alpha1clientset "github.com/projectriff/system/pkg/client/clientset/versioned/typed/stream/v1alpha1"
	apiextensionsv1beta1clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	kubernetes "k8s.io/client-go/kubernetes/fake"
	authv1client "k8s.io/client-go/kubernetes/typed/authorization/v1"
	corev1clientset "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type FakeClient struct {
	Namespace               string
	FakeKubeRestConfig      *rest.Config
	FakeKubeClient          *kubernetes.Clientset
	FakeRiffClient          *projectriffclientset.Clientset
	FakeApiExtensionsClient *apiextensionsv1beta1clientset.Clientset
	ActionRecorderList      kntesting.ActionRecorderList
}

func (c *FakeClient) DefaultNamespace() string {
	return c.Namespace
}

func (c *FakeClient) KubeRestConfig() *rest.Config {
	return c.FakeKubeRestConfig
}

func (c *FakeClient) Core() corev1clientset.CoreV1Interface {
	return c.FakeKubeClient.CoreV1()
}

func (c *FakeClient) Auth() authv1client.AuthorizationV1Interface {
	return c.FakeKubeClient.AuthorizationV1()
}

func (c *FakeClient) Build() buildv1alpha1clientset.BuildV1alpha1Interface {
	return c.FakeRiffClient.BuildV1alpha1()
}

func (c *FakeClient) Request() requestv1alpha1clientset.RequestV1alpha1Interface {
	return c.FakeRiffClient.RequestV1alpha1()
}

func (c *FakeClient) Stream() streamv1alpha1clientset.StreamV1alpha1Interface {
	return c.FakeRiffClient.StreamV1alpha1()
}
func (c *FakeClient) Discovery() discovery.DiscoveryInterface {
	return c.FakeKubeClient.Discovery()
}

func (c *FakeClient) ApiExtensions() apiextensionsv1beta1.ApiextensionsV1beta1Interface {
	return c.FakeApiExtensionsClient.ApiextensionsV1beta1()
}

func (c *FakeClient) PrependReactor(verb, resource string, reaction ReactionFunc) {
	c.FakeKubeClient.PrependReactor(verb, resource, reaction)
	c.FakeRiffClient.PrependReactor(verb, resource, reaction)
}

func NewClient(objects ...runtime.Object) *FakeClient {
	lister := NewListers(objects)

	kubeRestConfig := &rest.Config{Host: "https://localhost:8443"}
	kubeClient := kubernetes.NewSimpleClientset(lister.GetKubeObjects()...)
	riffClient := projectriffclientset.NewSimpleClientset(lister.GetProjectriffObjects()...)
	apiExtensionsObjects := lister.GetApiExtensionsObjects()
	apiExtensionsClient := apiextensionsv1beta1clientset.NewSimpleClientset(apiExtensionsObjects...)

	actionRecorderList := kntesting.ActionRecorderList{kubeClient, riffClient}

	return &FakeClient{
		Namespace:               "default",
		FakeKubeRestConfig:      kubeRestConfig,
		FakeKubeClient:          kubeClient,
		FakeRiffClient:          riffClient,
		FakeApiExtensionsClient: apiExtensionsClient,
		ActionRecorderList:      actionRecorderList,
	}
}
