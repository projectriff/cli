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

package commands_test

import (
	"fmt"
	"os"
	"testing"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/riff/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
)

func TestHandlerInvokeOptions(t *testing.T) {
	table := rifftesting.OptionsTable{
		{
			Name: "invalid resource",
			Options: &commands.HandlerInvokeOptions{
				ResourceOptions: rifftesting.InvalidResourceOptions,
			},
			ExpectFieldError: rifftesting.InvalidResourceOptionsFieldError,
		},
		{
			Name: "valid resource",
			Options: &commands.HandlerInvokeOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
			},
			ShouldValidate: true,
		},
		{
			Name: "json content type",
			Options: &commands.HandlerInvokeOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ContentTypeJSON: true,
			},
			ShouldValidate: true,
		},
		{
			Name: "text content type",
			Options: &commands.HandlerInvokeOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ContentTypeText: true,
			},
			ShouldValidate: true,
		},
		{
			Name: "multiple content types",
			Options: &commands.HandlerInvokeOptions{
				ResourceOptions: rifftesting.ValidResourceOptions,
				ContentTypeJSON: true,
				ContentTypeText: true,
			},
			ExpectFieldError: cli.ErrMultipleOneOf(cli.JSONFlagName, cli.TextFlagName),
		},
	}

	table.Run(t)
}

func TestHandlerInvokeCommand(t *testing.T) {
	table := rifftesting.CommandTable{
		// TODO write tests
	}

	table.Run(t, commands.NewHandlerInvokeCommand)
}

func TestHelperProcess_HandlerInvoke(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stderr, "Command executed: %s\n", shellquote.Join(argsAfterBareDoubleDash(os.Args)...))
	os.Exit(0)
}

func TestHelperProcess_HandlerInvokeError(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stderr, "Command executed: %s\n", shellquote.Join(argsAfterBareDoubleDash(os.Args)...))
	os.Exit(1)
}

func argsAfterBareDoubleDash(args []string) []string {
	for i, arg := range args {
		if arg == "--" {
			return args[i+1:]
		}
	}
	return []string{}
}
