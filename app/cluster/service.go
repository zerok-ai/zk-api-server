package cluster

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"main/app/cluster/errors"
	"main/app/cluster/models"
	"main/app/cluster/models/tablemux"
	clusterUtils "main/app/cluster/utils"
	"main/app/utils"
	"net/http"
	"px.dev/pxapi"
	"strings"
)

var authToken string
var apiKey string
var domain = "zkcloud01.getanton.com:443"

func init() {
	if authToken == "" {
		authToken = getTokenWith2ReTryInit()
		if authToken == "" {
			return
		}
	}
	resp := getMetaDataWithoutCtx()
	if resp.StatusCode == 401 {
		resp = getMetaDataWithoutCtx()
		if resp.StatusCode != 200 {
			return
		}
	} else if resp.StatusCode != 200 {
		return
	}

	finalResp := models.ClusterDetailsMetaDataResponse{}
	responseData, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(responseData, &finalResp)

	apiKey = finalResp.Data.ApiKey.Key
}

func listCluster(ctx iris.Context) {

	if authToken == "" {
		authToken = getTokenWith2ReTry(ctx)
		if authToken == "" {
			return
		}
	}

	resp := getMetaData(ctx)
	if resp.StatusCode == 401 {
		resp = getMetaData(ctx)
		if resp.StatusCode != 200 {
			return
		}
	} else if resp.StatusCode != 200 {
		return
	}

	finalResp := models.ClusterDetailsMetaDataResponse{}
	responseData, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(responseData, &finalResp)

	apiKey = finalResp.Data.ApiKey.Key
	clusters := make([]models.Cluster, 0)

	for _, r := range finalResp.Data.Clusters {
		clusters = append(clusters, r.FromResponseToDomainClusterDetails())
	}

	err := ctx.JSON(clusters)
	if err != nil {
		return
	}
}

func getMetaData(ctx iris.Context) http.Response {
	tr := &http.Transport{}
	client := http.Client{
		Transport: tr,
	}

	resp := clusterUtils.MakeRawApiCall("GET", nil, client, clusterUtils.CLUSTER_METADATA_URL, nil, nil, authToken)
	if resp.StatusCode == 401 {
		authToken = getTokenWith2ReTry(ctx)
		return resp
	} else if resp.StatusCode == 200 {
		return resp
	} else {
		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal Server Error"))
		return resp
	}
}

func getMetaDataWithoutCtx() http.Response {
	tr := &http.Transport{}
	client := http.Client{
		Transport: tr,
	}

	resp := clusterUtils.MakeRawApiCall("GET", nil, client, clusterUtils.CLUSTER_METADATA_URL, nil, nil, authToken)
	if resp.StatusCode == 401 {
		authToken = getTokenWith2ReTryInit()
		return resp
	} else if resp.StatusCode == 200 {
		return resp
	} else {
		return resp
	}
}

// TODO: make it better soon
func getTokenWith2ReTry(ctx iris.Context) string {
	// 1st call
	token := getAuthToken()
	if token == "" {
		// 2nd call
		token = getAuthToken()
		if token == "" {
			// 3rd call
			token = getAuthToken()
			if token == "" {
				_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
					Title("Internal Server Error, token generation error"))
				return ""
			}
		}
	}
	return token
}

// TODO: make it better soon
func getTokenWith2ReTryInit() string {
	// 1st call
	token := getAuthToken()
	if token == "" {
		// 2nd call
		token = getAuthToken()
		if token == "" {
			// 3rd call
			token = getAuthToken()
			if token == "" {

			}
		}
	}
	return token
}

// TODO: replace this soon
func getAuthToken() string {

	bodyMap := map[string]string{}
	bodyMap["email"] = clusterUtils.EMAIL
	bodyMap["password"] = clusterUtils.PASSWORD

	tr := &http.Transport{}
	client := http.Client{
		Transport: tr,
	}
	s, err := json.Marshal(bodyMap)
	if err != nil {
		return ""
	}
	bodyReader := strings.NewReader(string(s[:]))

	resp := clusterUtils.MakeRawApiCall("POST", utils.StringToPtr("application/json"), client, clusterUtils.LOGIN_URL, nil, bodyReader, "")
	var token string
	if resp.StatusCode == 200 {
		token = resp.Header.Get("Token")
	} else {
		token = ""
	}

	return token
}

func updateCluster(ctx iris.Context, cluster models.Cluster) {

	// Create a cluster if cluster.ID is missing. else, update if valid.
	if cluster.Id != nil {
		// Validation: Check if provided cluster ID is valid?
		//if !ValidateCluster(*cluster.Id, ctx) {
		//	return
		//}
		// Action: Update(replace) entire cluster info
		models.ClusterMap[*cluster.Id] = cluster
		ctx.StatusCode(iris.StatusOK)
	} else {
		// Action: Generate a UUID clusterID and add cluster info.
		clusterId := uuid.New()
		models.ClusterMap[clusterId.String()] = cluster
		ctx.StatusCode(iris.StatusCreated)
	}
}

func deleteCluster(ctx iris.Context, clusterId string) {
	//if !ValidateCluster(clusterId, ctx) {
	//	return
	//}

	delete(models.ClusterMap, clusterId)
	ctx.StatusCode(iris.StatusOK)
}

func GetClusterDetails(id string) models.Cluster {
	return models.Cluster{
		Domain:    domain,
		ApiKey:    apiKey,
		ClusterId: id,
	}
}

func getResourceDetails(ctx iris.Context, clusterIdx, action, st string) {
	if !ValidateAction(action, ctx) {
		return
	}
	//if !ValidateCluster(clusterIdx, ctx) {
	//	return
	//}

	//clusterDetails := models.ClusterMap[clusterIdx]
	//clusterDetails.Domain = clusterDetails.Domain + ":443"

	clusterDetails := GetClusterDetails(clusterIdx)

	//if namespace != "all" {
	//	var v = make([]string, 0)
	//	stringListMux := tablemux.StringListMux{Table: tablemux.TablePrinterStringList{Values: v}}
	//
	//	tx := models.MethodTemplate{MethodSignature: utils.GetNamespaceMethodSignature(st), DataFrameName: "my_first_ns"}
	//	_ = GetResource(ctx, clusterDetails, &stringListMux, tx)
	//	namespaceList := stringListMux.Table.Values
	//	if !utils.Contains(namespaceList, namespace) {
	//		_ = ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
	//			Title("Invalid namespace"))
	//		return
	//	}
	//}

	var resultSet *pxapi.ScriptResults
	var result interface{}
	if strings.EqualFold(action, "list") {
		resultSet, result = getServiceDetailsList(ctx, clusterDetails, st)
	} else if strings.EqualFold(action, "map") {
		resultSet, result = getServiceDetailsMap(ctx, clusterDetails, st)
	}

	if result == nil {
		return
	}

	_ = ctx.JSON(map[string]interface{}{
		"results": result,
		"stats":   resultSet.Stats(),
		"status":  200,
	})
}

func getServiceDetailsMap(ctx iris.Context, cluster models.Cluster, st string) (*pxapi.ScriptResults, []models.ServiceMap) {
	var s = make([]models.ServiceMap, 0)
	serviceListMux := tablemux.ServiceMapMux{Table: tablemux.TablePrinterServiceMap{Values: s}}
	tx := models.MethodTemplate{MethodSignature: utils.GetServiceMapMethodSignature(st), DataFrameName: "my_first_map"}
	resultSet := GetResource(ctx, cluster, &serviceListMux, tx, 3)
	return resultSet, serviceListMux.Table.Values
}

func getServiceDetailsList(ctx iris.Context, cluster models.Cluster, st string) (*pxapi.ScriptResults, []models.Service) {
	var s = make([]models.Service, 0)
	serviceListMux := tablemux.ServiceListMux{Table: tablemux.TablePrinterServiceList{Values: s}}
	tx := models.MethodTemplate{MethodSignature: utils.GetServiceListMethodSignature(st), DataFrameName: "my_first_list"}
	resultSet := GetResource(ctx, cluster, &serviceListMux, tx, 3)
	return resultSet, serviceListMux.Table.Values
}

func GetResource(ctx iris.Context, cluster models.Cluster, t tablemux.TableRecordHandler, tx models.MethodTemplate, retryCount int) *pxapi.ScriptResults {

	if retryCount == 0 {
		return nil
	}

	vz, pxl, ctxNew, err := tablemux.CreateVizierClient(cluster, tx)
	if err != nil {
		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return nil
	}
	resultSet, err := t.ExecutePxlScript(ctxNew, vz, pxl)
	if err != nil {
		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return nil
	}
	resultSet, err = tablemux.GetResult(resultSet)
	if err == errors.ErrAuthenticationFailed {
		return retry(ctx, cluster, t, tx, retryCount)
	} else if err != nil {
		_ = ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title(err.Error()))
		return nil
	}
	return resultSet
}

func retry(ctx iris.Context, cluster models.Cluster, t tablemux.TableRecordHandler, tx models.MethodTemplate, retryCount int) *pxapi.ScriptResults {
	if retryCount == 0 {
		return nil
	}
	retryCount -= 1

	if authToken == "" {
		authToken = getTokenWith2ReTryInit()
		if authToken == "" {
			return retry(ctx, cluster, t, tx, retryCount)
		}
	}
	resp := getMetaDataWithoutCtx()
	if resp.StatusCode == 401 {
		resp = getMetaDataWithoutCtx()
		if resp.StatusCode != 200 {
			return retry(ctx, cluster, t, tx, retryCount)
		}
	} else if resp.StatusCode != 200 {
		return retry(ctx, cluster, t, tx, retryCount)
	}

	finalResp := models.ClusterDetailsMetaDataResponse{}
	responseData, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(responseData, &finalResp)
	if err != nil {
		return retry(ctx, cluster, t, tx, retryCount)
	}

	cluster.ApiKey = finalResp.Data.ApiKey.Key
	apiKey = cluster.ApiKey

	resultSet := GetResource(ctx, cluster, t, tx, retryCount)

	return resultSet
}

func getServiceDetails(ctx iris.Context, clusterIdx, name, ns, st string) {
	//if !ValidateCluster(clusterIdx, ctx) {
	//	return
	//}

	//TODO: write validation for st

	clusterDetails := GetClusterDetails(clusterIdx)

	var resultSet *pxapi.ScriptResults
	var result interface{}

	var s = make([]models.ServiceStat, 0)
	serviceStatMux := tablemux.ServiceStatMux{Table: tablemux.TablePrinterServiceStat{Values: s}}
	resultSet = GetResource(ctx, clusterDetails, &serviceStatMux, models.MethodTemplate{MethodSignature: utils.GetServiceDetailsMethodSignature(st, ns+"/"+name), DataFrameName: "my_first_graph"}, 3)
	result = serviceStatMux.Table.Values

	if result == nil {
		return
	}

	_ = ctx.JSON(map[string]interface{}{
		"results": result,
		"stats":   resultSet.Stats(),
		"status":  200,
	})
}
