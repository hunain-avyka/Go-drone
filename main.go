// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"os"

	"github.com/hunain-avyka/Go-drone/command"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(new(command.Azure), "")
	subcommands.Register(new(command.Bitbucket), "")
	subcommands.Register(new(command.Circle), "")
	subcommands.Register(new(command.Cloudbuild), "")
	subcommands.Register(new(command.Drone), "")
	subcommands.Register(new(command.Github), "")
	subcommands.Register(new(command.Gitlab), "")
	subcommands.Register(new(command.Jenkins), "")
	subcommands.Register(new(command.Travis), "")
	subcommands.Register(new(command.Downgrade), "")
	subcommands.Register(new(command.JenkinsJson), "")
	subcommands.Register(new(command.JenkinsXml), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
