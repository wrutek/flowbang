package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RespItem struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

func (resp RespItem) GetId() int {
	return resp.Id
}

func (resp RespItem) GetName() string {
	return resp.Name
}

func (resp RespItem) GetFullName() string {
	return resp.FullName
}

func RawRequest(method string, uri string, headers *map[string]string, data interface{}) (respBody []RespItem, err error) {
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

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(r, &respBody)
	if err != nil {
		return
	}
	return
}

func DoNothing() {
	fmt.Println("I'll do nothing")
}
