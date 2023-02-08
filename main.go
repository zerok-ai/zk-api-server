package zerokApiServer

import (
	"github.com/kataras/iris"
	"zerok.ai/zerok-api-server/cluster"
	"zerok.ai/zerok-api-server/px"
)

func main() {
	app := iris.New()

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

	app.Listen(":8080")
}
