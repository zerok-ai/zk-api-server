package tablemux

import (
	"bytes"
	"context"
	"fmt"
	"main/app/cluster/models"
	"os"
	"px.dev/pxapi"
	"px.dev/pxapi/errdefs"
	"text/template"
)

type TableRecordHandler interface {
	ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error)
}

func CreateVizierClient(cluster models.Cluster, tx models.MethodTemplate) (*pxapi.VizierClient, string, context.Context, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, "", nil, err
	}

	pxFilePath := "/app/px/my.pxl"
	dat, err := os.ReadFile(path + pxFilePath)
	if err != nil {
		return nil, "", nil, err
	}
	t2 := template.New("Template")
	t2, _ = t2.Parse(string(dat))

	var doc bytes.Buffer
	err = t2.Execute(&doc, tx)
	if err != nil {
		return nil, "", nil, err
	}
	pxl := doc.String()
	fmt.Print(pxl)

	ctx := context.Background()
	client, err := pxapi.NewClient(ctx, pxapi.WithAPIKey(cluster.ApiKey), pxapi.WithCloudAddr(cluster.Domain))
	if err != nil {
		return nil, "", nil, err
	}

	vz, err := client.NewVizierClient(ctx, cluster.ClusterId)
	if err != nil {
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
		return nil, err
	}

	// Get the execution stats for the script execution.
	stats := resultSet.Stats()
	fmt.Printf("Execution Time: %v\n", stats.ExecutionTime)
	fmt.Printf("Bytes received: %v\n", stats.TotalBytes)

	return resultSet, nil
}
