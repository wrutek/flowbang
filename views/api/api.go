package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/wrutek/flowbang/config/configgetter"
)

// Get make a get request without body
func Get(uri string, out interface{}, headers map[string]string) (err error) {
	return Request("GET", uri, nil, out, headers)
}

// Post make a get request without body
func Post(uri string, data interface{}, out interface{}, headers map[string]string) (err error) {
	return Request("GET", uri, data, out, headers)
}

// Request make a request with already predefined settings like headers
func Request(method string, uri string, data interface{}, out interface{}, headers map[string]string) (err error) {
	cfg, err := configgetter.GetConfiguration()

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Authorization"] = "token " + cfg.OauthToken

	err = RawRequest(method, uri, headers, data, out)
	return
}

// RawRequest simple helper for making requests to github api.
// It will be mostly used durring configuration process where
// we do not have oauth token yet
func RawRequest(method string, uri string, headers map[string]string, data interface{}, out interface{}) (err error) {
	host := "https://api.github.com/"

	// sometimes dev will pass here a full url with api host.
	// Then we shouldn't add `host` string at the begining
	// of the request.
	if strings.Index(uri, host) != -1 {
		host = ""
	}
	req, err := http.NewRequest(method, host+uri, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	if headers != nil {
		for title, value := range headers {
			req.Header.Add(title, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer client.CloseIdleConnections()

	if resp.StatusCode > 299 {
		// TODO: http error handling
		errBodyByte, _ := ioutil.ReadAll(resp.Body)
		errBody := fmt.Sprintf("%s\n code: %d\n%s", req.URL.Path, resp.StatusCode, string(errBodyByte))
		panic(fmt.Errorf(errBody))
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(r, out)
	if err != nil {
		return
	}
	return
}
