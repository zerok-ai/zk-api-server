package tablemux

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"px.dev/pxapi"
	"px.dev/pxapi/errdefs"
	"text/template"
	"zk-api-server/app/utils/errors"

	"github.com/kataras/iris/v12"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
)

type PixieRepository interface {
	GetPixieData(ctx iris.Context, t pxapi.TableMuxer, tx MethodTemplate, clusterId string, apiKey string, domain string) (*pxapi.ScriptResults, *zkerrors.ZkError)
}

type pixie struct {
}

func NewPixieRepository() PixieRepository {
	return &pixie{}
}

var path string

func init() {
}

func CreateVizierClient(tx MethodTemplate, clusterId string, apiKey string, domain string) (*pxapi.VizierClient, string, context.Context, error) {
	path, err := os.Getwd()
	if err != nil {
		log.Printf("failed to get working dir, %s\n", err.Error())
		return nil, "", nil, err
	}

	pxFilePath := "/app/px/my.pxl"
	dat, err := os.ReadFile(path + pxFilePath)
	if err != nil {
		path = "/Users/vaibhavpaharia/Go/src/zk-api-server"
		dat, err = os.ReadFile(path + pxFilePath)
		if err != nil {
			log.Printf("failed to open pixel file, path: %s, err: %s\n", pxFilePath, err.Error())
			return nil, "", nil, err
		}
	}
	t := template.New("Template")
	t, _ = t.Parse(string(dat))

	var doc bytes.Buffer
	err = t.Execute(&doc, tx)
	if err != nil {
		log.Printf("failed to get working dir, %s\n", err.Error())
		return nil, "", nil, err
	}
	pxl := doc.String()

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
			return nil, errors.ErrAuthenticationFailed
		}
		return nil, err
	}

	return resultSet, nil
}

func (p *pixie) GetPixieData(ctx iris.Context, t pxapi.TableMuxer, tx MethodTemplate, clusterId string, apiKey string, domain string) (*pxapi.ScriptResults, *zkerrors.ZkError) {
	vz, pxl, ctxNew, err := CreateVizierClient(tx, clusterId, apiKey, domain)
	var zkErr zkerrors.ZkError
	if err != nil {
		log.Printf("failed to create vizier api client, error: %s\n", err.Error())
		zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return nil, &zkErr
	}

	resultSet, err := vz.ExecuteScript(ctxNew, pxl, t)
	if err != nil && err != io.EOF {
		log.Printf("failed to execute pixie script, error: %s\n", err.Error())
		zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return nil, &zkErr
	}

	resultSet, err = GetResult(resultSet)
	if err != nil {
		log.Printf("failed to get pixie data result, error: %s\n", err.Error())
		zkErr = zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return nil, &zkErr
	}

	return resultSet, nil
}
