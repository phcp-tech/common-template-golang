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
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"template/adapter"
	"template/infra/dao"
	"template/service"

	db "github.com/phcp-tech/common-library-golang/dbsqlc/postgres"
	dbLoader "github.com/phcp-tech/common-library-golang/dbsqlc/postgres/loader"
	"github.com/phcp-tech/common-library-golang/env"
	libGin "github.com/phcp-tech/common-library-golang/gin"
	"github.com/phcp-tech/common-library-golang/httpserver"
	httpserverLoader "github.com/phcp-tech/common-library-golang/httpserver/loader"
	"github.com/phcp-tech/common-library-golang/log"
	"github.com/phcp-tech/common-library-golang/shutdown"
	"golang.org/x/sync/errgroup"
)

func main() {
	// step 1: initial config file
	if err := env.InitEnv("config/app.toml"); err != nil {
		fmt.Printf("Initial environment config file failed: %s", err.Error())
		os.Exit(1)
	}

	// step 2: initial log
	cfg := log.Config{
		Level: env.Env().String("log.level"),
	}
	if env.Env().String("log.file.path") != "" {
		cfg.FilePath = env.Env().String("log.file.path")
		cfg.MaxSizeMB = env.Env().Int("log.file.max.size")
		cfg.MaxBackups = env.Env().Int("log.file.max.backups")
		cfg.MaxAgeDays = env.Env().Int("log.file.max.age")
		cfg.Compress = env.Env().Bool("log.file.compress")
	}
	log.InitLog(&cfg)
	log.Info("Initial environment config and log successfully.")
	defer func() {
		log.Info("Log file has been closed, application exit.")
		log.Close() // ensure flush logs before exit.
	}()

	// step 3: top-level recover to capture panics in main goroutine and record stack trace
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panic in main: %v\n%s", r, string(debug.Stack()))
		}
	}()

	// step 4: initial infrastructures
	initInfrastructures()
	defer func() {
		if conn := db.Default(); conn != nil {
			conn.Close()
		}
		log.Info("Database has been closed.")
	}()

	// step 5: initial services
	initServices()
	//defer service.Close() // ensure service resources are released before exit.

	// step 6: initial gin router
	var origins []string
	if strings.EqualFold(env.Env().String("app.env.value"), "prod") {
		origins = env.Env().Strings("cors.allow.origins.prod")
	} else {
		// add localhost:port origins to non-prod environment, enable local development
		origins = env.Env().Strings("cors.allow.origins.dev")
	}
	router := libGin.InitGin(origins)
	adapter.Mount(router)

	// step 7: load http server sequentially, start after all infrastructures and services are ready
	runner, _ := httpserverLoader.LoadDefault(router)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), httpserver.DefaultShutdownTimeout)
		defer cancel()
		if err := runner.Shutdown(ctx); err != nil {
			log.Errorf("http server shutdown failed: %s", err.Error())
		}
		log.Info("Http server has been shutdown.")
	}()

	// step 8: log application start info
	log.Infof("%s start successfully, version is %s, environment is %s.",
		env.Env().String("app.name"),
		env.Env().String("app.version"),
		env.Env().String("app.env.value"))

	// step 9: wait for shutdown signal
	shutdown.Wait()

	// step 10: other cleanup operations
}

// initial infrastructures concurrently
func initInfrastructures() {
	eg, _ := errgroup.WithContext(context.Background())
	// load default database
	eg.Go(func() error {
		return dbLoader.LoadDefault()
	})

	// wait for all infrastructures to be initialized
	if err := eg.Wait(); err != nil {
		log.Errorf("Init infrastructures failed, %s", err.Error())
		os.Exit(1)
	}
	log.Info("All infrastructures initialized successfully.")
}

// initial services
func initServices() {
	// initialize services
	userService := service.NewUserService(dao.NewUserDao())

	// inject services to adapter layer for RESTful API
	adapter.Svcs = &adapter.Services{UserService: userService}
	log.Info("All services initialized successfully.")
}
