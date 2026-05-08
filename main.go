// Copyright(C) 2020-2026 PHCP Technologies. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"

	"github.com/phcp-tech/common-library-golang/application"
	"github.com/phcp-tech/common-library-golang/env"
	"github.com/phcp-tech/common-library-golang/log"
)

func main() {
	// step 1: initial config file
	if err := env.InitEnv("config/app.toml"); err != nil {
		log.Errorf("Initial environment config file failed: %s", err.Error())
		os.Exit(1)
	}
	log.Info("Initial environment config file successful.")

	// step 2: wire and start application
	app := NewApplication()
	app.Start()

	// step 3: waiting for graceful exit
	application.WaitingForExitSignal(app)
}
