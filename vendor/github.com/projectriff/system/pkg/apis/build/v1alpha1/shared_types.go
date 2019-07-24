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
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ImageResource interface {
	metav1.ObjectMetaAccessor
	GetImage() string
}

// ResolveDefaultImage applies the default image prefix as needed to an image.
//
// The default image prefix may apply to either a repository whose value is '_'
// or a repository with a leading '_/'.
//
// For a leading '_/', the underscore is replaced with the default image prefix.
// For a repository of '_', the default image prefix is combined with the name
// of the build resource.
func ResolveDefaultImage(resource ImageResource, defaultImagePrefix string) (string, error) {
	if defaultImagePrefix == "" {
		return "", fmt.Errorf("invalid default image prefix %q", defaultImagePrefix)
	}
	image := resource.GetImage()
	if image == "_" {
		// combine registry prefix and application name
		image = fmt.Sprintf("%s/%s", defaultImagePrefix, resource.GetObjectMeta().GetName())
	} else if strings.HasPrefix(image, "_/") {
		// add the prefix to the specified image name
		image = strings.Replace(image, "_", defaultImagePrefix, 1)
	} else {
		return "", fmt.Errorf("unable to default registry")
	}
	return image, nil
}

type Source struct {
	// Git source location to clone and checkout for the build.
	Git *GitSource `json:"git"`

	// SubPath within the source to mount. Files outside of the sub path will
	// not be available to tbe build.
	SubPath string `json:"subPath,omitempty"`
}

type GitSource struct {
	// Revision in the git repository to checkout. May be any valid refspec
	// including commit sha, tags or branches.
	Revision string `json:"revision"`

	// URL to a cloneable git repository.
	URL string `json:"url"`
}

type BuildStatus struct {
	// BuildCacheName is the name of the PersistentVolumeClaim used as a cache
	// for intermediate build resources.
	BuildCacheName string `json:"buildCacheName,omitempty"`

	// BuildName is the name of the Knative Build backing this build.
	BuildName string `json:"buildName,omitempty"`

	// LatestImage is the most recent image for this build.
	LatestImage string `json:"latestImage,omitempty"`

	// TargetImage is the resolved image repository where built images are
	// pushed.
	TargetImage string `json:"targetImage,omitempty"`
}
