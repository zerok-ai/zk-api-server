package px

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/kataras/iris/v12"
	"px.dev/pxapi"
	"px.dev/pxapi/errdefs"
	"px.dev/pxapi/types"

	"main/app/cluster"
)

func GetPXData(ctx iris.Context) {
	clusterId := ctx.Params().Get("clusterId")
	clusterDetails := cluster.ClusterMap[clusterId]
	apiKey := clusterDetails.ApiKey
	cloudAddress := clusterDetails.Domain + ":443"
	resultSet := setupApiServer(apiKey, clusterId, cloudAddress)

	ctx.JSON(map[string]interface{}{
		"stats":   resultSet.Stats(),
		"results": Accumulator,
	})
}

func setupApiServer(apiKey string, clusterId string, cloudAddress string) *pxapi.ScriptResults {
	var API_KEY = apiKey
	var CLUSTER_ID = clusterId
	var CLOUD_ADDR = cloudAddress

	fmt.Println("API_KEY: %s, CLUSTER_ID: %s, CLOUD_ADDR: %s", API_KEY, CLUSTER_ID, CLOUD_ADDR)

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dat, err := os.ReadFile(path + "/app/px/getNamespaceHTTPTraffic.pxl")
	// dat, err := os.ReadFile(path + "/app/px/getMySQLData.pxl")
	if err != nil {
		panic(err)
	}
	pxl := string(dat)
	// fmt.Print(pxl)

	// Create a Pixie client.
	ctx := context.Background()
	client, err := pxapi.NewClient(ctx, pxapi.WithAPIKey(API_KEY), pxapi.WithCloudAddr(CLOUD_ADDR))
	if err != nil {
		panic(err)
	}

	// Create a connection to the cluster.
	vz, err := client.NewVizierClient(ctx, CLUSTER_ID)
	if err != nil {
		panic(err)
	}

	// Create TableMuxer to accept results table.
	tm := &tableMux{}

	// Execute the PxL script.
	resultSet, err := vz.ExecuteScript(ctx, pxl, tm)
	if err != nil && err != io.EOF {
		panic(err)
	}

	// Receive the PxL script results.
	defer resultSet.Close()
	if err := resultSet.Stream(); err != nil {
		if errdefs.IsCompilationError(err) {
			fmt.Printf("Got compiler error: \n %s\n", err.Error())
		} else {
			println("Error")
			fmt.Printf("Got error : %+v, while streaming\n", err)
		}
	}

	// Get the execution stats for the script execution.
	stats := resultSet.Stats()
	fmt.Printf("Execution Time: %v\n", stats.ExecutionTime)
	fmt.Printf("Bytes received: %v\n", stats.TotalBytes)

	return resultSet
}

var Accumulator []map[string]any

// Satisfies the TableRecordHandler interface.
type tablePrinter struct {
	TableRows []string
}

func (t *tablePrinter) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *tablePrinter) HandleRecord(ctx context.Context, r *types.Record) error {
	tempObj := make(map[string]any)
	fmt.Println(r.Data)
	for k, d := range r.Data {
		t.TableRows = append(t.TableRows, d.String())
		var colName string = r.TableMetadata.ColInfo[k].Name
		var value = d.String()
		// fmt.Println(colName, ":", value)

		var bufferSingleMap map[string]interface{}
		json.Unmarshal([]byte(value), &bufferSingleMap)
		if bufferSingleMap != nil {
			tempObj[colName] = bufferSingleMap
		} else {
			if num, err := strconv.Atoi(value); err == nil {
				tempObj[colName] = num
			} else {
				tempObj[colName] = value
			}
		}
	}
	Accumulator = append(Accumulator, tempObj)
	return nil
}

func (t *tablePrinter) HandleDone(ctx context.Context) error {
	return nil
}

// Satisfies the TableMuxer interface.
type tableMux struct {
	Table tablePrinter
}

func (s *tableMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	var Table = &tablePrinter{}
	s.Table = *Table
	return Table, nil
}
