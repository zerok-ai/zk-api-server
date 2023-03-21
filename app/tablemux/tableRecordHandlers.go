package tablemux

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"io"
	"log"
	"main/app/tablemux/handlerimplementation"
	"main/app/utils"
	"net/http"
	"os"
	"px.dev/pxapi"
	"px.dev/pxapi/errdefs"
	"strings"
	"text/template"
	"time"
)

type clusterDetails struct {
	Domain string `json:"Domain"`
	Url    string `json:"Url"`
}

var authToken string
var apiKey string
var details clusterDetails

var LoginEndpoint = "/v1/auth/login"
var Email = "admin@default.com"
var Password = "admin"
var ClusterMetadataEndpoint = "/v1/org/cluster/metadata"
var LoginUrl string
var ClusterMetadataUrl string

func init() {
	//configFilePath := "/Users/vaibhavpaharia/Go/src/zk-api-server/k8s/a.txt"
	configFilePath := "/opt/cluster.conf"

	jsonFile, err := os.Open(configFilePath)

	if err != nil {
		log.Println(err)
		os.Exit(2)
		return
	} else {
		defer jsonFile.Close()

		err = json.NewDecoder(jsonFile).Decode(&details)
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}

	}

	LoginUrl = details.Url + LoginEndpoint
	ClusterMetadataUrl = details.Url + ClusterMetadataEndpoint
}

type Cluster struct {
	Domain    string
	ApiKey    string
	ClusterId string
}

func getClusterDetails(id string) Cluster {
	return Cluster{
		Domain:    details.Domain,
		ApiKey:    apiKey,
		ClusterId: id,
	}
}

type TableRecordHandler interface {
	ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error)
}

func CreateVizierClient(cluster Cluster, tx MethodTemplate) (*pxapi.VizierClient, string, context.Context, error) {
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
	client, err := pxapi.NewClient(ctx, pxapi.WithAPIKey(cluster.ApiKey), pxapi.WithCloudAddr(cluster.Domain))
	if err != nil {
		log.Printf("failed to create pixie api client, error: %s\n", err.Error())
		return nil, "", nil, err
	}

	vz, err := client.NewVizierClient(ctx, cluster.ClusterId)
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

func GetResource(ctx iris.Context, id string, t TableRecordHandler, tx MethodTemplate, retryCount int) *pxapi.ScriptResults {
	if retryCount == 0 {
		return nil
	}

	retryCount -= 1

	cluster := getClusterDetails(id)

	vz, pxl, ctxNew, err := CreateVizierClient(cluster, tx)
	if err != nil {
		log.Printf("failed to create vizier api client, error: %s\n", err.Error())
		ctx.StatusCode(500)
		ctx.SetErr(utils.ErrInternalServerError)
	}
	resultSet, err := t.ExecutePxlScript(ctxNew, vz, pxl)
	if err != nil {
		log.Printf("failed to execute pixie script, error: %s\n", err.Error())
		ctx.StatusCode(500)
		ctx.SetErr(utils.ErrInternalServerError)
	}
	resultSet, err = GetResult(resultSet)
	if err == utils.ErrAuthenticationFailed {
		return retry(ctx, cluster, t, tx, retryCount, id)
	} else if err != nil {
		log.Printf("failed to get pixie data result, error: %s\n", err.Error())
		ctx.StatusCode(500)
		ctx.SetErr(utils.ErrInternalServerError)

	}
	return resultSet
}

func retry(ctx iris.Context, cluster Cluster, t TableRecordHandler, tx MethodTemplate, retryCount int, id string) *pxapi.ScriptResults {

	if authToken == "" {
		authToken = GetAuthTokenWith2ReTry(retryCount)
		if authToken == "" {
			return nil
		}
	}
	resp := GetMetaDataWithRetry(retryCount)
	if resp.StatusCode == 401 {
		log.Printf("metadata not retrieved, error: %d\n", resp.StatusCode)
		return nil
	} else if resp.StatusCode != 200 {
		log.Printf("metadata not retrieved, error: %d\n", resp.StatusCode)
		return nil
	}

	finalResp := handlerimplementation.ClusterDetailsMetaDataResponse{}
	responseData, _ := io.ReadAll(resp.Body)
	err := json.Unmarshal(responseData, &finalResp)
	if err != nil {
		log.Printf("json unmarshall error, error: %s\n", err.Error())
		return retry(ctx, cluster, t, tx, retryCount, id)
	}

	cluster.ApiKey = finalResp.Data.ApiKey.Key
	apiKey = cluster.ApiKey

	resultSet := GetResource(ctx, id, t, tx, retryCount)

	return resultSet
}

func PopulateApiKey() {
	r := GetMetaDataWithRetry(3)
	finalResp := handlerimplementation.ClusterDetailsMetaDataResponse{}
	responseData, _ := io.ReadAll(r.Body)
	_ = json.Unmarshal(responseData, &finalResp)

	apiKey = finalResp.Data.ApiKey.Key
}

func GetMetaDataWithRetry(retryCount int) http.Response {

	var resp http.Response
	if utils.IsEmpty(authToken) {
		return resp
	}

	for i := 0; i < retryCount; i++ {
		resp = GetMetaData()
		if resp.StatusCode == 200 {
			return resp
		} else if resp.StatusCode == 401 {
			log.Printf("metadata not retrieved, error: %d\n", resp.StatusCode)
			authToken = getAuthToken()
		}
	}
	return resp
}

func GetMetaData() http.Response {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	return utils.MakeRawApiCall("GET", ClusterMetadataUrl, nil, nil, nil, authToken, nil, client)
}

func GetAuthTokenWith2ReTry(retryCount int) string {
	for i := 0; i < retryCount; i++ {
		authToken = getAuthToken()
		if authToken != "" {
			return authToken
		}
	}
	return authToken
}

func getAuthToken() string {

	bodyMap := map[string]string{}
	bodyMap["email"] = Email
	bodyMap["password"] = Password

	tr := &http.Transport{}
	client := http.Client{
		Transport: tr,
	}
	s, err := json.Marshal(bodyMap)
	if err != nil {
		log.Printf("json unmarshall error, error: %s\n", err.Error())
		return ""
	}
	bodyReader := strings.NewReader(string(s[:]))

	resp := utils.MakeRawApiCall("POST", LoginUrl, nil, bodyReader, map[string]string{"content-type": "application/json"}, "", nil, client)
	var token string
	if resp.StatusCode == 200 {
		token = resp.Header.Get("Token")
	} else {
		log.Printf("token could not be fetched, error: %d\n", resp.StatusCode)
		token = ""
	}

	return token
}

func UpdateApiKey(k string) {
	if utils.IsEmpty(k) {
		return
	}
	apiKey = k
}
