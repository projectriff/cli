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

package commands

import (
	"context"
	"strings"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/spf13/cobra"
)

func NewFunctionCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "function",
		Short: "functions built from source using function buildpacks",
		Long: strings.TrimSpace(`
Functions are a mechanism for converting language idiomatic units of logic into
container images that can be HTTP invoked or process streams of messages. Cloud
Native Buildpacks are provided to detect the language, provide a language
runtime, install build and runtime dependencies, compile the function, and
package it into a container.

The function resource is only responsible for converting source code into a
container. The function container image may then be deployed as a request
handler, or as a stream processor.

Functions are distinct from applications in the scope of the source code. Unlike
applications, functions:

- practice Inversion of Control (we'll call you)
- are decoupled from networking protocols, no HTTP specifics
- limited to a single responsibility

`),
		Args:    cli.Args(),
		Aliases: []string{"functions", "func", "fn"},
	}

	cmd.AddCommand(NewFunctionListCommand(ctx, c))
	cmd.AddCommand(NewFunctionCreateCommand(ctx, c))
	cmd.AddCommand(NewFunctionDeleteCommand(ctx, c))
	cmd.AddCommand(NewFunctionStatusCommand(ctx, c))
	cmd.AddCommand(NewFunctionTailCommand(ctx, c))

	return cmd
}
