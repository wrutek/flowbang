package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RawRequest(method string, uri string, headers map[string]string, data interface{}) (respBody map[string]interface{}, err error) {
	req, err := http.NewRequest(method, "https://api.github.com/"+uri, nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	for title, value := range headers {
		req.Header.Add(title, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer client.CloseIdleConnections()
	_ = json.NewDecoder(resp.Body).Decode(&respBody)
	fmt.Println(respBody)
	return
}

func DoNothing() {
	fmt.Println("I'll do nothing")
}
