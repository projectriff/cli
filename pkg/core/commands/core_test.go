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
	"context"
	"testing"

	"github.com/projectriff/cli/pkg/cli"
	"github.com/projectriff/cli/pkg/core/commands"
	rifftesting "github.com/projectriff/cli/pkg/testing"
)

func TestCoreCommand(t *testing.T) {
	table := rifftesting.CommandTable{
		{
			Name:        "empty",
			Args:        []string{},
			ShouldError: true,
		},
	}

	table.Run(t, commands.NewCoreCommand)
}

func TestCoreAliasCommand(t *testing.T) {
	config := cli.NewDefaultConfig()
	context := context.TODO()

	expected := []string{"c"}

	result := commands.NewCoreCommand(context, config)

	if len(result.Aliases) == 0 {
		t.Errorf("NewCoreCommand failed, expected %v, got nothing", expected)
	}

	for i := range result.Aliases {
		if result.Aliases[i] != expected[i] {
			t.Errorf("NewCoreCommand failed, expected %v, got %v", expected[i], result.Aliases[i])
		}
	}

}
