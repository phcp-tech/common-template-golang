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

	"template/adapter"
	"template/domain/model"
	"template/infra/dao"
	"template/service"

	"github.com/phcp-tech/common-library-golang/application"
	"github.com/phcp-tech/common-library-golang/db"
	"github.com/phcp-tech/common-library-golang/env"
	libGin "github.com/phcp-tech/common-library-golang/gin"
	"github.com/phcp-tech/common-library-golang/httpserver"
	"github.com/phcp-tech/common-library-golang/log"

	"github.com/gin-gonic/gin"
)

// compile-time interfae check
var _ application.IApplication = (*Application)(nil)

type Application struct {
	userService *service.UserService
}

func NewApplication() *Application {
	app := &Application{}
	app.initInfrastructures()
	app.initServices()
	return app
}

func (app *Application) initInfrastructures() {
	application.InitInfrastructures()
}

func (app *Application) initServices() {
	// initialize services
	userDao := dao.NewUserDao()
	app.userService = service.NewUserService(userDao)

	// inject services to adapter layer for RESTful API
	adapter.Svcs = &adapter.Services{UserService: app.userService}
	log.Info("All services initialized successfully.")
}

func (app *Application) Start() {
	// step 1: auto migrate
	db.AutoMigrate("", &model.User{})

	// step 2: initial gin and mount controller/pprof
	router := libGin.InitGin()
	adapter.Mount(router)

	// step 3: start http server
	port := env.Env().String("http.server.port")
	go func(r *gin.Engine) {
		if err := httpserver.Startup(r, port, "", ""); err != nil {
			log.Errorf("Startup http server failed: %s", err.Error())
			os.Exit(1)
		}
	}(router)

	log.Infof("%s start successful, version is %s, environment is %s.",
		env.Env().String("app.name"),
		env.Env().String("app.version"),
		env.Env().String("app.env.value"))
}

func (app *Application) Shutdown(sig os.Signal) {
	application.Shutdown(sig)
}
