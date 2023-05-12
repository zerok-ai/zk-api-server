package utils

import (
	"github.com/kataras/iris/v12"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateApiKeyMiddleware_Fail(t *testing.T) {
	app := iris.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := app.ContextPool.Acquire(httptest.NewRecorder(), req)

	ctx.Request().Header.Set("Zk_api_key", "")

	ValidateApiKeyMiddleware(ctx)

	assert.Equal(t, http.StatusUnauthorized, ctx.ResponseWriter().StatusCode())
}

func TestValidateApiKeyMiddleware_Success(t *testing.T) {
	app := iris.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := app.ContextPool.Acquire(httptest.NewRecorder(), req)

	ctx.Request().Header.Set("Zk_api_key", "SOME_VALUE_HERE")

	ValidateApiKeyMiddleware(ctx)

	assert.Equal(t, http.StatusOK, ctx.ResponseWriter().StatusCode())
}
