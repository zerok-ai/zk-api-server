package main

import (
	"main/app/cluster"
	"main/internal/model"

	"github.com/kataras/iris/v12"
	zkConfig "github.com/zerok-ai/zk-utils-go/config"
	httpConfig "github.com/zerok-ai/zk-utils-go/http/config"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkPostgres "github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres"
)

var LogTag = "main"

func main() {
	var cfg model.ZkApiServerConfig
	if err := zkConfig.ProcessArgs[model.ZkApiServerConfig](&cfg); err != nil {
		panic(err)
	}

	zkLogger.Info(LogTag, "")
	zkLogger.Info(LogTag, "********* Initializing Application *********")
	httpConfig.Init(cfg.Http.Debug)
	zkPostgres.Init(cfg.Postgres)
	zkLogger.Init(cfg.LogsConfig)
	zkLogger.Debug(LogTag, "Parsed Configuration", cfg)

	app := newApp()

	config := iris.WithConfiguration(iris.Configuration{
		DisablePathCorrection: true,
		LogLevel:              "debug",
	})
	app.Listen(":"+cfg.Server.Port, config)
}

func newApp() *iris.Application {
	app := iris.Default()

	crs := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Method() == iris.MethodOptions {
			ctx.Header("Access-Control-Methods",
				"POST, PUT, PATCH, DELETE")

			ctx.Header("Access-Control-Allow-Headers",
				"Access-Control-Allow-Origin,Content-Type")

			ctx.Header("Access-Control-Max-Age",
				"86400")

			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		ctx.Next()
	}

	app.UseRouter(crs)

	app.AllowMethods(iris.MethodOptions)

	app.Get("/healthz", func(ctx iris.Context) {
		ctx.WriteString("pong")
	}).Describe("healthcheck")

	v1 := app.Party("/v1")
	cluster.Initialize(v1)

	return app
}
