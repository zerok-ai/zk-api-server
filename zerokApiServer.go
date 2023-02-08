package zerokApiServer

import (
	"zerokApiServer/cluster"
	"zerokApiServer/px"

	"github.com/kataras/iris/v12"
)

func newApp() *iris.Application {
	app := iris.Default()

	v1 := app.Party("/v1")
	clusterAPI := v1.Party("/cluster")
	{
		clusterAPI.Use(iris.Compression)

		clusterAPI.Get("/", cluster.List)
	}

	pxAPI := v1.Party("/px")
	{
		pxAPI.Get("/", px.GetPXData)
	}

	return app
}

func Start() {
	app := newApp()
	app.Listen(":8080")
}
