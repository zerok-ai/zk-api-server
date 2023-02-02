package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"px.dev/pxapi"
	"px.dev/pxapi/errdefs"
	"px.dev/pxapi/types"
)

// Define PxL script with one table output.
var (

// 	pxl = `
// import px
// df = px.DataFrame('http_events')
// df = df[['upid', 'req_path', 'remote_addr', 'req_method']]
// df = df.head(10)
// px.display(df, 'http')
// `
)

func main() {
	var API_KEY = "px-api-b65a1334-5d5a-40a3-9e84-b21f46990764"
	// var API_KEY_ID = "d820bb34-e4a0-4d8d-878c-f8d43984167d"
	var CLUSTER_ID = "89a7771a-17f4-4577-9cf9-6ec605c9b16b"
	var CLOUD_ADDR = "pxtest1.getanton.com:443"

	dat, err := os.ReadFile("./getNamespaceHTTPTraffic.pxl")
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
}

// Satisfies the TableRecordHandler interface.
type tablePrinter struct{}

func (t *tablePrinter) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *tablePrinter) HandleRecord(ctx context.Context, r *types.Record) error {
	for _, d := range r.Data {
		fmt.Printf("%s ", d.String())
	}
	fmt.Printf("\n")
	return nil
}

func (t *tablePrinter) HandleDone(ctx context.Context) error {
	return nil
}

// Satisfies the TableMuxer interface.
type tableMux struct {
}

func (s *tableMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &tablePrinter{}, nil
}
