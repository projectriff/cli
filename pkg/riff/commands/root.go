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
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/projectriff/cli/pkg/cli"
	"github.com/spf13/cobra"
)

// NewRootCommand wraps the riff command with flags and commands that should only be defined on
// the CLI root.
func NewRootCommand(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := NewRiffCommand(ctx, c)

	cmd.Use = c.Name
	if !c.GitDirty {
		cmd.Version = fmt.Sprintf("%s (%s)", c.Version, c.GitSha)
	} else {
		cmd.Version = fmt.Sprintf("%s (%s, with local modifications)", c.Version, c.GitSha)
	}
	cmd.DisableAutoGenTag = true

	// add root persistent flags
	cmd.PersistentFlags().StringVar(&c.ViperConfigFile, cli.StripDash(cli.ConfigFlagName), "", fmt.Sprintf("config `file` (default is $HOME/.%s.yaml)", c.Name))
	cmd.PersistentFlags().StringVar(&c.KubeConfigFile, cli.StripDash(cli.KubeConfigFlagName), "", "kubectl config `file` (default is $HOME/.kube/config)")
	cmd.PersistentFlags().BoolVar(&color.NoColor, cli.StripDash(cli.NoColorFlagName), color.NoColor, "disable color output in terminals")

	// add root-only commands
	cmd.AddCommand(NewCompletionCommand(ctx, c))
	cmd.AddCommand(NewDocsCommand(ctx, c))
	cmd.AddCommand(NewDoctorCommand(ctx, c))

	// override usage template to add arguments
	cmd.SetUsageTemplate(strings.ReplaceAll(cmd.UsageTemplate(), "{{.UseLine}}", "{{useLine .}}"))
	cobra.AddTemplateFunc("useLine", func(cmd *cobra.Command) string {
		result := cmd.UseLine()
		flags := ""
		if strings.HasSuffix(result, " [flags]") {
			flags = " [flags]"
			result = result[0 : len(result)-len(flags)]
		}
		return result + cli.FormatArgs(cmd) + flags
	})

	return cmd
}
