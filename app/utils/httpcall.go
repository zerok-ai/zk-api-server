package utils

import (
	"io"
	"log"
	"net/http"
)

func MakeRawApiCall(method string, urlToBeCalled string, queryParams map[string]string, requestBody io.Reader, headers map[string]string, authToken string, cookiesTobeAdded []http.Cookie, client http.Client) http.Response {
	var req *http.Request
	var err error

	if method == "GET" {
		req, err = makeGetRequest(urlToBeCalled, authToken)
	} else if method == "POST" {
		req, err = makePostRequest(urlToBeCalled, requestBody)
	} else {
		log.Printf("Invalid method type, %s", method)
	}

	if err != nil {
		return http.Response{}
	}

	addHeaders(headers, req)
	addCookies(cookiesTobeAdded, req)
	addQueryParams(queryParams, req)

	response, err := client.Do(req)
	if err != nil {
		log.Printf("Got error %s\n", err.Error())
		return http.Response{}
	}

	return *response
}

func makePostRequest(url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Printf("Got error %s\n", err.Error())
		return nil, err
	}
	return req, nil

}

func makeGetRequest(url, authToken string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Got error %s\n", err.Error())
		return nil, err
	}

	if !IsEmpty(authToken) {
		req.Header.Add("Token", authToken)
	}
	return req, nil
}

func addCookies(c []http.Cookie, req *http.Request) *http.Request {
	if c != nil {
		for _, element := range c {
			req.AddCookie(&element)
		}
	}
	return req
}

func addHeaders(h map[string]string, req *http.Request) *http.Request {
	if h != nil {
		for k, v := range h {
			req.Header.Add(k, v)
		}
	}
	return req
}

func addQueryParams(queryParams map[string]string, req *http.Request) *http.Request {
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return req
}
