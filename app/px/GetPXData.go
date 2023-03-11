package px

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/app/cluster/models"
	"os"
	"strconv"
	"text/template"

	"github.com/kataras/iris/v12"
	"px.dev/pxapi"
	"px.dev/pxapi/errdefs"
	"px.dev/pxapi/types"
)

type Template struct {
	StartTime string
	Head      int
	Filter    string
}

func GetPXData(ctx iris.Context) {
	clusterMapId := ctx.URLParamDefault("cluster_id", "1")
	clusterDetails := models.ClusterMap[clusterMapId]
	apiKey := clusterDetails.ApiKey
	cloudAddress := clusterDetails.Domain + ":443"
	clusterId := clusterDetails.ClusterId
	resultSet, err := setupApiServer(apiKey, clusterId, cloudAddress)

	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return
	}

	ctx.JSON(map[string]interface{}{
		"stats":   resultSet.Stats(),
		"results": Accumulator,
	})
}

func setupApiServer(apiKey string, clusterId string, cloudAddress string) (*pxapi.ScriptResults, error) {
	var API_KEY = apiKey
	var CLUSTER_ID = clusterId
	var CLOUD_ADDR = cloudAddress

	fmt.Printf("API_KEY: %s, CLUSTER_ID: %s, CLOUD_ADDR: %s\n", API_KEY, CLUSTER_ID, CLOUD_ADDR)

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dat, err := os.ReadFile(path + "/app/px/getROI.pxl")
	// dat, err := os.ReadFile(path + "/app/px/getMySQLData.pxl")
	if err != nil {
		return nil, err
	}
	t2 := template.New("Template")
	t2, _ = t2.Parse(string(dat))
	tx := Template{"-20s", 100, "{}"}

	var doc bytes.Buffer
	t2.Execute(&doc, tx)
	pxl := doc.String()
	fmt.Print(pxl)

	// Create a Pixie client.
	ctx := context.Background()
	client, err := pxapi.NewClient(ctx, pxapi.WithAPIKey(API_KEY), pxapi.WithCloudAddr(CLOUD_ADDR))
	if err != nil {
		return nil, err
	}

	// Create a connection to the cluster.
	vz, err := client.NewVizierClient(ctx, CLUSTER_ID)
	if err != nil {
		return nil, err
	}

	// Create TableMuxer to accept results table.
	tm := &tableMux{}

	// Execute the PxL script.
	resultSet, err := vz.ExecuteScript(ctx, pxl, tm)
	if err != nil && err != io.EOF {
		return nil, err
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
		return nil, err
	}

	// Get the execution stats for the script execution.
	stats := resultSet.Stats()
	fmt.Printf("Execution Time: %v\n", stats.ExecutionTime)
	fmt.Printf("Bytes received: %v\n", stats.TotalBytes)

	return resultSet, nil
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
	// fmt.Println(r.Data)
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
