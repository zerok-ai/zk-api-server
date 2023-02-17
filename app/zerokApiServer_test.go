package app

import (
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func TestNewApp(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	e.GET("/v1/cluster").Expect().Status(httptest.StatusOK)
	e.GET("/v1/px").Expect().Status(httptest.StatusOK)

}
