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
	"net/http"

	"template/adapter"
	"template/infra/dao"
	"template/service"

	"github.com/gin-gonic/gin"
	"github.com/phcp-tech/common-library-golang/bootstrap"
	dbComp "github.com/phcp-tech/common-library-golang/dbsqlx/postgres/component"
	"github.com/phcp-tech/common-library-golang/env"
	envComp "github.com/phcp-tech/common-library-golang/env/component"
	ginComp "github.com/phcp-tech/common-library-golang/gin/component"
	httpComp "github.com/phcp-tech/common-library-golang/httpserver/component"
	"github.com/phcp-tech/common-library-golang/log"
	logComp "github.com/phcp-tech/common-library-golang/log/component"
	"github.com/phcp-tech/common-library-golang/network"
)

func main() {
	var router *gin.Engine
	bootstrap.New().
		Add(envComp.Component("config/app.toml")). // 1st - env
		Add(logComp.Component()).                  // 2nd - log
		AddParallel(dbComp.Component()).
		PreReady(initServices).
		Add(ginComp.Component(func(r *gin.Engine) {
			router = r
			adapter.Mount(r)
		})).
		Add(httpComp.Component(func() http.Handler { return router })).
		PostReady(func() {
			log.Infof("%s start successfully, version is %s, environment is %s, local ip address are %s",
				env.Env().String("app.name"),
				env.Env().String("app.version"),
				env.Env().String("app.env.value"),
				network.GetLocalIpAddress())
		}).
		Run()
}

// initial services
func initServices() error {
	// inject services to adapter layer for RESTful API
	userService := service.NewUserService(dao.NewUserDao())
	adapter.Svcs = &adapter.Services{UserService: userService}

	log.Info("All services initialized successfully")
	return nil
}
