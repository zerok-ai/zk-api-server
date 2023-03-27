package tablemux

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"io"
	"log"
	"main/app/utils"
	"os"
	"px.dev/pxapi"
	"px.dev/pxapi/errdefs"
	"text/template"
)

func CreateVizierClient(tx MethodTemplate, clusterId string, apiKey string, domain string) (*pxapi.VizierClient, string, context.Context, error) {
	path, err := os.Getwd()
	if err != nil {
		log.Printf("failed to get working dir, %s\n", err.Error())
		return nil, "", nil, err
	}

	pxFilePath := "/app/px/my.pxl"
	dat, err := os.ReadFile(path + pxFilePath)
	if err != nil {
		log.Printf("failed to open pixel file, path: %s, err: %s\n", pxFilePath, err.Error())
		return nil, "", nil, err
	}
	t2 := template.New("Template")
	t2, _ = t2.Parse(string(dat))

	var doc bytes.Buffer
	err = t2.Execute(&doc, tx)
	if err != nil {
		log.Printf("failed to get working dir, %s\n", err.Error())
		return nil, "", nil, err
	}
	pxl := doc.String()
	fmt.Print(pxl)

	ctx := context.Background()
	client, err := pxapi.NewClient(ctx, pxapi.WithAPIKey(apiKey), pxapi.WithCloudAddr(domain))
	if err != nil {
		log.Printf("failed to create pixie api client, error: %s\n", err.Error())
		return nil, "", nil, err
	}

	vz, err := client.NewVizierClient(ctx, clusterId)
	if err != nil {
		log.Printf("failed to create vizier api client, error: %s\n", err.Error())
		return nil, "", nil, err
	}

	return vz, pxl, ctx, nil

}

func GetResult(resultSet *pxapi.ScriptResults) (*pxapi.ScriptResults, error) {
	// Receive the PxL script results.
	defer func(resultSet *pxapi.ScriptResults) {
		err := resultSet.Close()
		if err != nil {

		}
	}(resultSet)
	if err := resultSet.Stream(); err != nil {
		if errdefs.IsCompilationError(err) {
			fmt.Printf("Got compiler error: \n %s\n", err.Error())
		} else {
			println("Error")
			fmt.Printf("Got error : %+v, while streaming\n", err)
		}
		if err.Error() == "rpc error: code = Internal desc = Auth middleware failed: failed to fetch token - unauthenticated" {
			return nil, utils.ErrAuthenticationFailed
		}
		return nil, err
	}

	// Get the execution stats for the script execution.
	stats := resultSet.Stats()
	fmt.Printf("Execution Time: %v\n", stats.ExecutionTime)
	fmt.Printf("Bytes received: %v\n", stats.TotalBytes)

	return resultSet, nil
}

func GetResource[C pxapi.TableMuxer](ctx iris.Context, t C, tx MethodTemplate, clusterId string, apiKey string, domain string) *pxapi.ScriptResults {
	vz, pxl, ctxNew, err := CreateVizierClient(tx, clusterId, apiKey, domain)
	if err != nil {
		log.Printf("failed to create vizier api client, error: %s\n", err.Error())
		ctx.StatusCode(500)
		ctx.SetErr(utils.ErrInternalServerError)
	}
	resultSet, err := vz.ExecuteScript(ctxNew, pxl, t)
	if err != nil && err != io.EOF {
		log.Printf("failed to execute pixie script, error: %s\n", err.Error())
		ctx.StatusCode(500)
		ctx.SetErr(utils.ErrInternalServerError)
	}

	resultSet, err = GetResult(resultSet)
	if err != nil {
		log.Printf("failed to get pixie data result, error: %s\n", err.Error())
		ctx.StatusCode(500)
		ctx.SetErr(utils.ErrInternalServerError)

	}
	return resultSet
}
