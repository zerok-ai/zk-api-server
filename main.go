package main

import (
	zkapp "main/app"
	clusterHandler "main/app/cluster/handler"
	scenarioHandler "main/app/scenario/handler"
	"main/app/scenario/repository"
	"main/app/scenario/service"
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
	zkLogger.Init(cfg.LogsConfig)
	zkLogger.Debug(LogTag, "Parsed Configuration", cfg)
	zkPostgresRepo, err := zkPostgres.NewZkPostgresRepo(cfg.Postgres)
	if err != nil {
		return
	}

	zkLogger.Debug(LogTag, "Parsed Configuration", cfg)

	rr := repository.NewZkPostgresRepo(zkPostgresRepo)
	rs := service.NewScenarioService(rr)
	rh := scenarioHandler.NewScenarioHandler(rs)
	ch := clusterHandler.NewClusterHandler()

	app := newApp(rh, ch)

	config := iris.WithConfiguration(iris.Configuration{
		DisablePathCorrection: true,
		LogLevel:              "debug",
	})
	app.Listen(":"+cfg.Server.Port, config)
}

func newApp(rh scenarioHandler.ScenarioHandler, ch clusterHandler.ClusterHandler) *iris.Application {
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
	zkapp.Initialize(v1, rh, ch)

	return app
}
