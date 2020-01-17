/*
 * Copyright 2020 the original author or authors.
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
	"fmt"

	"github.com/projectriff/cli/pkg/cli"
)

func ContainerConcurrency(value int64, field string) cli.FieldErrors {
	errs := cli.FieldErrors{}

	if value < 0 || value > 1000 {
		errs = errs.Also(cli.ErrInvalidValue(fmt.Sprint(value), field))
	}

	return errs
}