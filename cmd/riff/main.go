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

package main

import (
	"context"
	"errors"
	"os"
	"strings"

	// load credential helpers
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/riff/commands"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func main() {
	ctx := context.Background()
	c := cli.Initialize()
	cmd := commands.NewRootCommand(ctx, c)

	cmd.SilenceErrors = true
	if err := cmd.Execute(); err != nil {
		// silent errors should not log, but still exit with an error code
		// typically the command has already been logged with more detail
		if !errors.Is(err, cli.SilentError) {
			c.Errorf("Error executing command:\n")

			if aggregate, ok := err.(utilerrors.Aggregate); ok {
				for _, err := range aggregate.Errors() {
					c.Errorf("  %s\n", err.Error())
				}
			} else {
				// errors can be multiple lines, indent each line
				for _, line := range strings.Split(err.Error(), "\n") {
					c.Errorf("  %s\n", line)
				}
			}
		}
		os.Exit(1)
	}
}
