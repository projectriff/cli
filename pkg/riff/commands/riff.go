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

func NewRiffCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "riff",
		Short: "riff is for functions",
		Long: strings.TrimSpace(`
The ` + c.Name + ` CLI combines with the projectriff system CRDs to build, run and wire
workloads (applications and functions). The CRDs provide the riff API of which
this CLI is a client.

Before running ` + c.Name + `, please install the projectriff system and its dependencies.
See https://projectriff.io/docs/getting-started/

This CLI contains commands that fit into three themes:

Builds - the application and function commands to define build plans and the
credential commands to authenticate builds to container registries.

Requests - the handler command to map HTTP requests to a built application,
function or container image.

Streams - the stream and processor commands to define streams of messages and
map those streams to function inputs and outputs with processors.
`),
		Args: cli.Args(),
	}

	cmd.AddCommand(NewCredentialCommand(ctx, c))
	cmd.AddCommand(NewApplicationCommand(ctx, c))
	cmd.AddCommand(NewFunctionCommand(ctx, c))
	cmd.AddCommand(NewHandlerCommand(ctx, c))
	cmd.AddCommand(NewStreamCommand(ctx, c))
	cmd.AddCommand(NewProcessorCommand(ctx, c))

	return cmd
}
