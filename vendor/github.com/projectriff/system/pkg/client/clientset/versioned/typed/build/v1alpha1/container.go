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
	v1alpha1 "github.com/projectriff/system/pkg/apis/build/v1alpha1"
	scheme "github.com/projectriff/system/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ContainersGetter has a method to return a ContainerInterface.
// A group's client should implement this interface.
type ContainersGetter interface {
	Containers(namespace string) ContainerInterface
}

// ContainerInterface has methods to work with Container resources.
type ContainerInterface interface {
	Create(*v1alpha1.Container) (*v1alpha1.Container, error)
	Update(*v1alpha1.Container) (*v1alpha1.Container, error)
	UpdateStatus(*v1alpha1.Container) (*v1alpha1.Container, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Container, error)
	List(opts v1.ListOptions) (*v1alpha1.ContainerList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Container, err error)
	ContainerExpansion
}

// containers implements ContainerInterface
type containers struct {
	client rest.Interface
	ns     string
}

// newContainers returns a Containers
func newContainers(c *BuildV1alpha1Client, namespace string) *containers {
	return &containers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the container, and returns the corresponding container object, and an error if there is any.
func (c *containers) Get(name string, options v1.GetOptions) (result *v1alpha1.Container, err error) {
	result = &v1alpha1.Container{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("containers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Containers that match those selectors.
func (c *containers) List(opts v1.ListOptions) (result *v1alpha1.ContainerList, err error) {
	result = &v1alpha1.ContainerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("containers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested containers.
func (c *containers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("containers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a container and creates it.  Returns the server's representation of the container, and an error, if there is any.
func (c *containers) Create(container *v1alpha1.Container) (result *v1alpha1.Container, err error) {
	result = &v1alpha1.Container{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("containers").
		Body(container).
		Do().
		Into(result)
	return
}

// Update takes the representation of a container and updates it. Returns the server's representation of the container, and an error, if there is any.
func (c *containers) Update(container *v1alpha1.Container) (result *v1alpha1.Container, err error) {
	result = &v1alpha1.Container{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("containers").
		Name(container.Name).
		Body(container).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *containers) UpdateStatus(container *v1alpha1.Container) (result *v1alpha1.Container, err error) {
	result = &v1alpha1.Container{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("containers").
		Name(container.Name).
		SubResource("status").
		Body(container).
		Do().
		Into(result)
	return
}

// Delete takes name of the container and deletes it. Returns an error if one occurs.
func (c *containers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("containers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *containers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("containers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched container.
func (c *containers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Container, err error) {
	result = &v1alpha1.Container{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("containers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
