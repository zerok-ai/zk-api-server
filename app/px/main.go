package px

import "github.com/kataras/iris/v12/core/router"

func Initialize(app router.Party) {
	pxAPI := app.Party("/px")
	{
		pxAPI.Get("/", GetPXData)
	}
}
