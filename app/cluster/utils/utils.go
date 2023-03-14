package utils

import (
	"io"
	"log"
	"main/app/utils"
	"net/http"
	"os"
	"px.dev/pxapi/types"
	"strconv"
)

func GetString(key string, r *types.Record) string {
	return r.GetDatum(key).String()
}

func GetFloat(key string, r *types.Record, bitSize int) (float64, error) {
	return strconv.ParseFloat(GetString(key, r), bitSize)
}

func GetInteger(key string, r *types.Record) (int, error) {
	return strconv.Atoi(GetString(key, r))
}

func GetStringPtr(key string, r *types.Record) *string {
	return utils.StringToPtr(GetString(key, r))
}

func GetFloat64Ptr(key string, r *types.Record) *float64 {
	v, err := GetFloat(key, r, 64)
	if err != nil {
		return nil
	} else {
		return &v
	}
}

func GetFloat32Ptr(key string, r *types.Record) *float64 {
	v, err := GetFloat(key, r, 32)
	if err != nil {
		return nil
	} else {
		return &v
	}
}

func GetIntegerPtr(key string, r *types.Record) *int {
	v, err := strconv.Atoi(GetString(key, r))
	if err != nil {
		return nil
	} else {
		return &v
	}
}

// TODO: the method needs major refactoring or even could be re-written
func MakeRawApiCall(method string, contentType *string, client http.Client, urlToBeCalled string, cookiesTobeAdded []http.Cookie, requestBody io.Reader, authToken string) http.Response {

	req, err := http.NewRequest(method, urlToBeCalled, requestBody)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
	}

	if cookiesTobeAdded != nil {
		for _, element := range cookiesTobeAdded {
			req.AddCookie(&element)
		}
	}

	if contentType != nil {
		// log.Println("Adding content type ", *contentType)
		req.Header.Add("Content-Type", *contentType)
	}
	if method == "GET" {
		req.Header.Add("Token", authToken)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}

	return *response
}

var LOGIN_URL = "http://zk-auth-demo.getanton.com:80/v1/auth/login"
var EMAIL = "admin@default.com"
var PASSWORD = "admin"

var CLUSTER_METADATA_URL = "http://zk-auth-demo.getanton.com/v1/org/cluster/metadata"
