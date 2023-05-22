package main

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/kataras/iris/v12"
	"main/app/cluster"
	httpConfig "main/app/utils/http/config"
	zkLogger "main/app/utils/logs"
	zkPostgres "main/app/utils/postgres"
	"main/internal/model"
	"os"
)

var LOG_TAG = "main"

type Args struct {
	ConfigPath string
}

func main() {
	var cfg model.ZkApiServerConfig
	// path := "/opt/zk-auth-configmap.yaml"
	args := ProcessArgs(&cfg)

	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(args.ConfigPath, &cfg); err != nil {
		zkLogger.Error(LOG_TAG, err)
		os.Exit(2)
	}

	zkLogger.Info(LOG_TAG, "")
	zkLogger.Info(LOG_TAG, "********* Initializing Application *********")
	httpConfig.Init(cfg.Http.Debug)
	zkPostgres.Init(cfg.Postgres)
	zkLogger.Init(cfg.LogsConfig)
	zkLogger.Debug(LOG_TAG, "Parsed Configuration", cfg)

	app := newApp()

	config := iris.WithConfiguration(iris.Configuration{
		DisablePathCorrection: true,
		LogLevel:              "debug",
	})
	app.Listen(":"+cfg.Server.Port, config)
}

func ProcessArgs(cfg interface{}) Args {
	var a Args

	f := flag.NewFlagSet("Example server", 1)
	f.StringVar(&a.ConfigPath, "c", "config.yaml", "Path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
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

	v1 := app.Party("/v1")
	cluster.Initialize(v1)

	return app
}
