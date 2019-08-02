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

package parsers

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func EnvVar(str string) corev1.EnvVar {
	parts := strings.SplitN(str, "=", 2)

	return corev1.EnvVar{
		Name:  parts[0],
		Value: parts[1],
	}
}

func EnvVarValueFrom(str string) corev1.EnvVar {
	parts := strings.SplitN(str, "=", 2)
	source := strings.SplitN(parts[1], ":", 3)

	envVar := corev1.EnvVar{
		Name: parts[0],
	}
	if source[0] == "secretKeyRef" {
		envVar.ValueFrom = &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: source[1],
				},
				Key: source[2],
			},
		}
	}
	if source[0] == "configMapKeyRef" {
		envVar.ValueFrom = &corev1.EnvVarSource{
			ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: source[1],
				},
				Key: source[2],
			},
		}
	}

	return envVar
}

func EnvFromSource(str string) corev1.EnvFromSource {
	refType := ""
	refName := ""
	prefix := ""
	parts := strings.SplitN(str, ":", 3)
	switch len(parts) {
	case 2:
		prefix = ""
		refType = parts[0]
		refName = parts[1]
	case 3:
		prefix = parts[0]
		refType = parts[1]
		refName = parts[2]
	}

	envFromSource := corev1.EnvFromSource{}
	if refType == "secretRef" {
		envFromSource.SecretRef = &corev1.SecretEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: refName,
			},
		}
	}
	if refType == "configMapRef" {
		envFromSource.ConfigMapRef = &corev1.ConfigMapEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: refName,
			},
		}
	}

	if prefix != "" {
		envFromSource.Prefix = prefix
	}

	return envFromSource
}
