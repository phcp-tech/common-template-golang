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

	"template/adapter"
	"template/infra/dao"
	"template/service"

	db "github.com/phcp-tech/common-library-golang/dbsqlc/postgres"
	dbLoader "github.com/phcp-tech/common-library-golang/dbsqlc/postgres/loader"
	"github.com/phcp-tech/common-library-golang/env"
	libGin "github.com/phcp-tech/common-library-golang/gin"
	"github.com/phcp-tech/common-library-golang/httpserver"
	"github.com/phcp-tech/common-library-golang/log"
	"github.com/phcp-tech/common-library-golang/network"
	"github.com/phcp-tech/common-library-golang/shutdown"
	"golang.org/x/sync/errgroup"
)

func main() {
	// step 1: top-level recover to capture panics in main goroutine and record stack trace
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panic in main: %v\nstack: %s", r, string(debug.Stack()))
		}
	}()

	// step 2: initial config file
	if err := env.InitEnv("config/app.toml"); err != nil {
		fmt.Printf("Initial environment config file failed: %s", err.Error())
		os.Exit(1)
	}

	// step 3: initial log
	log.InitLog(&log.Config{
		Level: env.Env().String("log.level"),
	})
	log.Info("Initial environment config and log successfully")
	defer func() {
		log.Info("Log file has been closed, application exit")
		log.Close() // ensure flush logs before exit.
	}()

	// step 4: initial infrastructures
	if err := initInfrastructures(); err != nil {
		log.Errorf("Initial infrastructures failed: %s", err.Error())
		os.Exit(1)
	}
	defer func() {
		if conn := db.Default(); conn != nil {
			conn.Close()
		}
		log.Info("Database has been closed")
	}()

	// step 5: initial services
	if err := initServices(); err != nil {
		log.Errorf("Initial services failed: %s", err.Error())
		os.Exit(1)
	}
	//defer service.Close() // ensure service resources are released before exit.

	// step 6: initial gin router
	router := libGin.InitGin(env.Env().Strings("cors.allow.origins"))
	adapter.Mount(router)

	// step 7: create http server runner, then start it in a goroutine
	port := env.Env().String("http.server.port")
	httpServer := httpserver.NewHttpServer(httpserver.Config{
		Port: port,
	})
	log.Infof("Http server is running under Virtual Machine, listen on port %s", port)
	defer func() {
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Errorf("http server shutdown failed: %s", err.Error())
		}
		log.Info("Http server has been shutdown")
	}()
	// start http server in a separate goroutine, capture panics and log stack trace if any panic happens
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Panic in http server goroutine: %v\nstack: %s", r, string(debug.Stack()))
				shutdown.Trigger()
			}
		}()
		if err := httpServer.Start(router); err != nil {
			log.Errorf("Http server start with error: %s", err.Error())
			shutdown.Trigger()
		}
	}()

	// step 8: log application start info
	log.Infof("%s start successfully, version is %s, environment is %s, local ip address are %s",
		env.Env().String("app.name"),
		env.Env().String("app.version"),
		env.Env().String("app.env.value"),
		network.GetLocalIpAddress())

	// step 9: wait for shutdown signal
	shutdown.Wait()

	// step 10: other cleanup operations
}

// initial infrastructures concurrently
func initInfrastructures() error {
	eg, _ := errgroup.WithContext(context.Background())
	// load default database
	eg.Go(func() error {
		return dbLoader.LoadFromEnv()
	})

	// wait for all infrastructures to be initialized
	if err := eg.Wait(); err != nil {
		return err
	}

	log.Info("All infrastructures initialized successfully")
	return nil
}

// initial services
func initServices() error {
	// inject services to adapter layer for RESTful API
	userService := service.NewUserService(dao.NewUserDao())
	adapter.Svcs = &adapter.Services{UserService: userService}

	log.Info("All services initialized successfully")
	return nil
}
