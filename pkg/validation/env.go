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

package validation

import (
	"strings"

	"github.com/knative/pkg/apis"
)

func EnvVar(env, field string) *apis.FieldError {
	errs := &apis.FieldError{}

	if strings.HasPrefix(env, "=") || !strings.Contains(env, "=") {
		errs = errs.Also(apis.ErrInvalidValue(env, field))
	}

	return errs
}

func EnvVars(envs []string, field string) *apis.FieldError {
	errs := &apis.FieldError{}

	for i, env := range envs {
		errs = errs.Also(EnvVar(env, apis.CurrentField).ViaFieldIndex(field, i))
	}

	return errs
}

func EnvVarValueFrom(env, field string) *apis.FieldError {
	errs := &apis.FieldError{}

	parts := strings.SplitN(env, "=", 2)
	if len(parts) != 2 || parts[0] == "" {
		errs = errs.Also(apis.ErrInvalidValue(env, field))
	} else {
		value := strings.SplitN(parts[1], ":", 3)
		if len(value) != 3 {
			errs = errs.Also(apis.ErrInvalidValue(env, field))
		} else if value[0] != "configMapKeyRef" && value[0] != "secretKeyRef" {
			errs = errs.Also(apis.ErrInvalidValue(env, field))
		} else if value[1] == "" {
			errs = errs.Also(apis.ErrInvalidValue(env, field))
		}
	}

	return errs
}

func EnvVarValueFroms(envs []string, field string) *apis.FieldError {
	errs := &apis.FieldError{}

	for i, env := range envs {
		errs = errs.Also(EnvVarValueFrom(env, apis.CurrentField).ViaFieldIndex(field, i))
	}

	return errs
}

func EnvFromSource(env, field string) *apis.FieldError {
	errs := &apis.FieldError{}

	parts := strings.SplitN(env, ":", 3)
	if len(parts) < 2 || len(parts) > 3 {
		errs = errs.Also(apis.ErrInvalidValue(env, field))
	} else {
		if len(parts) == 3 {
			if parts[1] != "configMapRef" && parts[1] != "secretRef" {
				errs = errs.Also(apis.ErrInvalidValue(env, field))
			} else if len(parts[2]) < 1 {
				errs = errs.Also(apis.ErrInvalidValue(env, field))
			}
		} else if len(parts) == 2 {
			if parts[0] != "configMapRef" && parts[0] != "secretRef" {
				errs = errs.Also(apis.ErrInvalidValue(env, field))
			} else if len(parts[1]) < 1 {
				errs = errs.Also(apis.ErrInvalidValue(env, field))
			}
		}
	}

	return errs
}

func EnvFromSources(envs []string, field string) *apis.FieldError {
	errs := &apis.FieldError{}

	for i, env := range envs {
		errs = errs.Also(EnvFromSource(env, apis.CurrentField).ViaFieldIndex(field, i))
	}

	return errs
}
