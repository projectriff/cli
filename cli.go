/*
 * Copyright 2018 the original author or authors.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *  
 *        http://www.apache.org/licenses/LICENSE-2.0
 *  
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package main

import (
	"github.com/projectriff/riff/riff-cli/cmd"
	"github.com/projectriff/riff/riff-cli/global"
	"fmt"
	"os"
	"github.com/projectriff/riff/riff-cli/pkg/docker"
)

var version = "Unknown"

func main() {
	global.CLI_VERSION = version

	rootCmd := cmd.CreateAndWireRootCommand(docker.RealDocker(), docker.DryRunDocker())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}