package main

import (
	zkapp "zk-api-server/app"
	attributeHandler "zk-api-server/app/attribute/handler"
	attributeRepo "zk-api-server/app/attribute/repository"
	attributeService "zk-api-server/app/attribute/service"
	integrationsHandler "zk-api-server/app/integrations/handler"
	integrationRepo "zk-api-server/app/integrations/repository"
	integrationService "zk-api-server/app/integrations/service"
	scenarioHandler "zk-api-server/app/scenario/handler"
	"zk-api-server/app/scenario/repository"
	"zk-api-server/app/scenario/service"
	"zk-api-server/internal/model"

	"github.com/kataras/iris/v12"
	zkConfig "github.com/zerok-ai/zk-utils-go/config"
	httpConfig "github.com/zerok-ai/zk-utils-go/http/config"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkPostgres "github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres"

	obfuscationHandler "zk-api-server/app/obfuscation/handler"
	obfuscationRepo "zk-api-server/app/obfuscation/repository"
	obfuscationService "zk-api-server/app/obfuscation/service"
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
	rh := scenarioHandler.NewScenarioHandler(rs, cfg)

	ir := integrationRepo.NewZkPostgresRepo(zkPostgresRepo)
	is := integrationService.NewIntegrationsService(ir)
	ih := integrationsHandler.NewIntegrationsHandler(is, cfg)

	ar := attributeRepo.NewZkPostgresRepo(zkPostgresRepo)
	as := attributeService.NewAttributeService(ar)
	ah := attributeHandler.NewAttributeHandler(as, cfg)

	or := obfuscationRepo.NewZkPostgresObfuscationRepo(zkPostgresRepo)
	os := obfuscationService.NewObfuscationService(or)
	oh := obfuscationHandler.NewObfuscationHandler(os, cfg)

	app := newApp(rh, ih, ah, oh)

	config := iris.WithConfiguration(iris.Configuration{
		DisablePathCorrection: true,
		LogLevel:              "debug",
	})
	app.Listen(":"+cfg.Server.Port, config)
}

func newApp(rh scenarioHandler.ScenarioHandler, ih integrationsHandler.IntegrationsHandler, ah attributeHandler.AttributeHandler, oh obfuscationHandler.ObfuscationHandler) *iris.Application {
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
	zkapp.Initialize(v1, rh, ih, ah, oh)

	return app
}
