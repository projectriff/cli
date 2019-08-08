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

func NewConfigurerCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configurer",
		Short: "configurers map HTTP requests to a workload",
		Long: strings.TrimSpace(`
Configurers can be created for a build reference or image. Build based configurers
continuously watch for the latest built image and will deploy new images. If the
underlying build resource is deleted, the configurer will continue to run, but will
no longer self update. Image based configurers must be manually updated to trigger
roll out of an updated image.

Users wishing to perform checks on built images before deploying them can
provide their own external process to watch the build resource for new images
and only update the configurer image once those checks pass.

The hostname to access the configurer is available in the configurer listing.
`),
		Aliases: []string{"configurers"},
	}

	cmd.AddCommand(NewConfigurerListCommand(ctx, c))
	cmd.AddCommand(NewConfigurerCreateCommand(ctx, c))
	cmd.AddCommand(NewConfigurerDeleteCommand(ctx, c))
	cmd.AddCommand(NewConfigurerStatusCommand(ctx, c))
	cmd.AddCommand(NewConfigurerTailCommand(ctx, c))
	cmd.AddCommand(NewConfigurerInvokeCommand(ctx, c))

	return cmd
}
