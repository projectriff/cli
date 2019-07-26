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

func NewContainerCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "container",
		Short: "containers built from source using container buildpacks",
		Long: strings.TrimSpace(`
Containers are a mechanism to convert web container source code into
container images that can be invoked over HTTP. Cloud Native Buildpacks are
provided to detect the language, provide a language runtime, install build and
runtime dependencies, compile the container, and packaging everything as a
container.

The container resource is only responsible for converting source code into a
container. The container container image may then be deployed as a request
handler. See ` + "`" + c.Name + " handler --help" + "`" + ` for detail.
`),
		Aliases: []string{"containers"},
	}

	cmd.AddCommand(NewContainerListCommand(ctx, c))
	cmd.AddCommand(NewContainerCreateCommand(ctx, c))
	cmd.AddCommand(NewContainerDeleteCommand(ctx, c))
	cmd.AddCommand(NewContainerStatusCommand(ctx, c))

	return cmd
}
