package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAddHeadersWithValidData(t *testing.T) {
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer some-token",
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	result := addHeaders(header, req)

	for k, v := range header {
		assert.Equal(t, result.Header.Get(k), v)
	}
}

func TestAddHeadersWithEmptyHeaderMap(t *testing.T) {
	header := map[string]string{}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	result := addHeaders(header, req)

	assert.Equal(t, len(result.Header), 0)
}

func TestAddQueryParamsWithValidData(t *testing.T) {
	queryParams := map[string]string{
		"name": "John",
		"age":  "30",
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	result := addQueryParams(queryParams, req)

	// expectedQuery := "name=John&age=30"
	// since we are passing map, so we cannot be certain whether QueryParam will name=John&age=30 or age=30&name=John
	query := result.URL.RawQuery
	assert.True(t, strings.Contains(query, "name=John"))
	assert.True(t, strings.Contains(query, "&"))
	assert.True(t, strings.Contains(query, "age=30"))
}

func TestAddQueryParamsWithEmptyQueryParams(t *testing.T) {
	queryParams := map[string]string{}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	result := addQueryParams(queryParams, req)

	assert.Equal(t, result.URL.RawQuery, "")
}

func TestAddCookiesWithValidData(t *testing.T) {
	cookie1 := http.Cookie{
		Name:  "username",
		Value: "elonmusk",
	}
	cookie2 := http.Cookie{
		Name:  "sessionid",
		Value: "123456789",
	}
	cookies := []http.Cookie{cookie1, cookie2}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	result := addCookies(cookies, req)

	cookieHeader := result.Header.Get("Cookie")
	expectedCookieHeader := "username=elonmusk; sessionid=123456789"
	assert.Equal(t, cookieHeader, expectedCookieHeader)
	if cookieHeader != expectedCookieHeader {
		t.Errorf("addCookies didn't add cookies correctly, got %v but expected %v", cookieHeader, expectedCookieHeader)
	}
}

func TestAddCookiesWithEmptyCookieSlice(t *testing.T) {
	cookies := []http.Cookie{}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	result := addCookies(cookies, req)

	cookieHeader := result.Header.Get("Cookie")
	assert.Equal(t, cookieHeader, "")
}

func TestMakeGetRequestWithValidData(t *testing.T) {
	url := "http://example.com"
	authToken := "my-auth-token"

	result, err := makeGetRequest(url, authToken)

	assert.Nil(t, err)
	assert.Equal(t, result.Method, "GET")
	assert.Equal(t, result.URL.String(), url)
	assert.Equal(t, result.Header.Get("Token"), authToken)
}

func TestMakeGetRequestWithEmptyAuthToken(t *testing.T) {
	url := "http://example.com"
	authToken := ""

	result, err := makeGetRequest(url, authToken)

	assert.Nil(t, err)
	assert.Equal(t, result.Method, "GET")
	assert.Equal(t, result.URL.String(), url)
	assert.Equal(t, result.Header.Get("Token"), "")
}

func TestMakePostRequestWithValidData(t *testing.T) {
	url := "http://example.com"
	bodyStr := "request body"
	body := strings.NewReader(bodyStr)

	result, err := makePostRequest(url, body)
	x, err := io.ReadAll(result.Body)
	actualReqBody := string(x[:])

	assert.Nil(t, err)
	assert.Equal(t, result.Method, "POST")
	assert.Equal(t, result.URL.String(), url)
	assert.Equal(t, bodyStr, actualReqBody)
}

func TestMakePostRequestWithEmptyBody(t *testing.T) {
	url := "http://example.com"
	bodyStr := ""

	result, err := makePostRequest(url, strings.NewReader(bodyStr))
	x, err := io.ReadAll(result.Body)
	actualReqBody := string(x[:])

	assert.Nil(t, err)
	assert.Equal(t, result.Method, "POST")
	assert.Equal(t, result.URL.String(), url)
	assert.Equal(t, bodyStr, actualReqBody)
}

func TestMakeRawApiCall_GET_Success(t *testing.T) {
	mockResponse := http.Response{StatusCode: http.StatusOK}
	mockClient := &http.Client{Transport: &mockTransport{mockResponse, nil}}

	response, err := MakeRawApiCall("GET", "http://example.com", nil, nil, nil, "", nil, *mockClient)

	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Nil(t, err)
}

func TestMakeRawApiCall_POST_Success(t *testing.T) {
	mockResponse := http.Response{StatusCode: http.StatusOK}
	mockClient := &http.Client{Transport: &mockTransport{mockResponse, nil}}

	requestBody := strings.NewReader("request body")

	response, err := MakeRawApiCall("POST", "http://example.com", nil, requestBody, nil, "", nil, *mockClient)

	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)

}

func TestMakeRawApiCall_InvalidMethod(t *testing.T) {
	response, err := MakeRawApiCall("INVALID", "http://example.com", nil, nil, nil, "", nil, http.Client{})

	assert.Equal(t, 0, response.StatusCode)
	assert.Equal(t, "invalid method type", err.Error())

}

func TestMakeRawApiCall_FailedRequest(t *testing.T) {
	mockClient := &http.Client{Transport: &mockTransport{http.Response{}, errors.New("new error")}}

	response, err := MakeRawApiCall("GET", "http://example.com", nil, nil, nil, "", nil, *mockClient)

	assert.Equal(t, 0, response.StatusCode)
	assert.Equal(t, "Get \"http://example.com\": new error", err.Error())
}

type mockTransport struct {
	response http.Response
	err      error
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// return the mock response or error
	if t.err != nil {
		return nil, t.err
	}
	return &t.response, nil
}
