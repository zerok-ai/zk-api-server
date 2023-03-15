package tablemux

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"main/app/tablemux/handlerimplementation"
	"main/app/utils"
	"net/http"
	"os"
	"px.dev/pxapi"
	"px.dev/pxapi/errdefs"
	"strings"
	"text/template"
)

var authToken string
var apiKey string
var domain = "zkcloud01.getanton.com:443"
var LOGIN_URL = "http://zk-auth-demo.getanton.com:80/v1/auth/login"
var EMAIL = "admin@default.com"
var PASSWORD = "admin"
var CLUSTER_METADATA_URL = "http://zk-auth-demo.getanton.com/v1/org/cluster/metadata"

type Cluster struct {
	Domain    string
	ApiKey    string
	ClusterId string
}

func getClusterDetails(id string) Cluster {
	return Cluster{
		Domain:    domain,
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
		// TODO: add logs
		return nil, "", nil, err
	}

	pxFilePath := "/app/px/my.pxl"
	dat, err := os.ReadFile(path + pxFilePath)
	if err != nil {
		// TODO: add logs
		return nil, "", nil, err
	}
	t2 := template.New("Template")
	t2, _ = t2.Parse(string(dat))

	var doc bytes.Buffer
	err = t2.Execute(&doc, tx)
	if err != nil {
		// TODO: add logs
		return nil, "", nil, err
	}
	pxl := doc.String()
	fmt.Print(pxl)

	ctx := context.Background()
	client, err := pxapi.NewClient(ctx, pxapi.WithAPIKey(cluster.ApiKey), pxapi.WithCloudAddr(cluster.Domain))
	if err != nil {
		// TODO: add logs
		return nil, "", nil, err
	}

	vz, err := client.NewVizierClient(ctx, cluster.ClusterId)
	if err != nil {
		// TODO: add logs
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

	cluster := getClusterDetails(id)

	vz, pxl, ctxNew, err := CreateVizierClient(cluster, tx)
	if err != nil {
		ctx.StatusCode(500)
		ctx.SetErr(utils.ErrInternalServerError)
		// TODO: remove line below
		// TODO: add logs

		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return nil
	}
	resultSet, err := t.ExecutePxlScript(ctxNew, vz, pxl)
	if err != nil {
		ctx.StatusCode(500)
		ctx.SetErr(utils.ErrInternalServerError)
		// TODO: remove line below
		// TODO: add logs

		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return nil
	}
	resultSet, err = GetResult(resultSet)
	if err == utils.ErrAuthenticationFailed {
		return retry(ctx, cluster, t, tx, retryCount, id)
	} else if err != nil {
		ctx.StatusCode(500)
		ctx.SetErr(utils.ErrInternalServerError)
		// TODO: remove line below
		// TODO: add logs

		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return nil
	}
	return resultSet
}

// todo: check the implementation, retry seems to be having some logical flaw
func retry(ctx iris.Context, cluster Cluster, t TableRecordHandler, tx MethodTemplate, retryCount int, id string) *pxapi.ScriptResults {
	if retryCount == 0 {
		return nil
	}
	retryCount -= 1

	if authToken == "" {
		authToken = GetAuthTokenWith2ReTry(3)
		if authToken == "" {
			return retry(ctx, cluster, t, tx, retryCount, id)
		}
	}
	resp := GetMetaDataWithRetry(3)
	if resp.StatusCode == 401 {
		// TODO: add logs
		resp = GetMetaDataWithRetry(3)
		if resp.StatusCode != 200 {
			// TODO: add logs
			return retry(ctx, cluster, t, tx, retryCount, id)
		}
	} else if resp.StatusCode != 200 {
		// TODO: add logs
		return retry(ctx, cluster, t, tx, retryCount, id)
	}

	finalResp := handlerimplementation.ClusterDetailsMetaDataResponse{}
	responseData, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(responseData, &finalResp)
	if err != nil {
		// TODO: add logs
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
	responseData, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(responseData, &finalResp)

	apiKey = finalResp.Data.ApiKey.Key
}

func GetMetaDataWithRetry(retryCount int) http.Response {

	var resp http.Response
	for i := 0; i < retryCount; i++ {
		resp = GetMetaData()
		if resp.StatusCode == 200 {
			return resp
		} else if resp.StatusCode == 401 {
			// TODO: add logs
			authToken = GetAuthTokenWith2ReTry(3)
		}
	}
	return resp
}

func GetMetaData() http.Response {
	tr := &http.Transport{}
	client := http.Client{
		Transport: tr,
	}

	return utils.MakeRawApiCall("GET", nil, client, CLUSTER_METADATA_URL, nil, nil, authToken)
}

// GetAuthTokenWith2ReTry TODO: make it better soon
func GetAuthTokenWith2ReTry(retryCount int) string {
	for i := 0; i < retryCount; i++ {
		authToken = getAuthToken()
		if authToken != "" {
			return authToken
		}
	}
	return authToken
}

// getAuthToken TODO: replace this soon
func getAuthToken() string {

	bodyMap := map[string]string{}
	bodyMap["email"] = EMAIL
	bodyMap["password"] = PASSWORD

	tr := &http.Transport{}
	client := http.Client{
		Transport: tr,
	}
	s, err := json.Marshal(bodyMap)
	if err != nil {
		// TODO: add logs
		return ""
	}
	bodyReader := strings.NewReader(string(s[:]))

	resp := utils.MakeRawApiCall("POST", utils.StringToPtr("application/json"), client, LOGIN_URL, nil, bodyReader, "")
	var token string
	if resp.StatusCode == 200 {
		token = resp.Header.Get("Token")
	} else {
		// TODO: add logs
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
