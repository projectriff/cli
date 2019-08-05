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

func NewHandlerCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler",
		Short: "handlers map HTTP requests to a workload",
		Long: strings.TrimSpace(`
Handlers can be created for a build reference or image. Build based handlers
continuously watch for the latest built image and will deploy new images. If the
underlying build resource is deleted, the handler will continue to run, but will
no longer self update. Image based handlers must be manually updated to trigger
roll out of an updated image.

Users wishing to perform checks on built images before deploying them can
provide their own external process to watch the build resource for new images
and only update the handler image once those checks pass.

The hostname to access the handler is available in the handler listing.
`),
		Aliases: []string{"handlers"},
	}

	cmd.AddCommand(NewHandlerListCommand(ctx, c))
	cmd.AddCommand(NewHandlerCreateCommand(ctx, c))
	cmd.AddCommand(NewHandlerDeleteCommand(ctx, c))
	cmd.AddCommand(NewHandlerStatusCommand(ctx, c))
	cmd.AddCommand(NewHandlerTailCommand(ctx, c))
	cmd.AddCommand(NewHandlerInvokeCommand(ctx, c))

	return cmd
}
