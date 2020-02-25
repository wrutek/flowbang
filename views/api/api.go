package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func RawRequest(method string, uri string, headers *map[string]string, data interface{}, out interface{}) (err error) {
	req, err := http.NewRequest(method, "https://api.github.com/"+uri, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	if headers != nil {
		for title, value := range *headers {
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

func DoNothing() {
	fmt.Println("I'll do nothing")
}
