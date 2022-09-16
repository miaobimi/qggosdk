package qgproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// 接口地址
type ApiUrl string

const (
	AllocateUrl ApiUrl = "/allocate"
)

func Allocate(params map[string]interface{}) ([]byte, error) {
	params = map[string]interface{}{"format": "json", "area": "北京,上海"}
	paramsJson, _ := json.Marshal(params)

	client := new(http.Client)
	resp, err := DoReq(client, "proxy.qg.net", string(AllocateUrl), "GET", url.Values{}, paramsJson)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetReq(host, path, method string, query url.Values, body []byte) (*http.Request, error) {
	targetUrl := url.URL{
		Scheme:  "http",
		Host:    host,
		Path:    path,
		RawPath: query.Encode(),
	}
	if body == nil {
		return http.NewRequest(method, targetUrl.String(), nil)
	}
	return http.NewRequest(method, targetUrl.String(), bytes.NewReader(body))
}

func DoReq(httpCli *http.Client, host string, method, path string, query url.Values, body []byte) ([]byte, error) {
	req, err := GetReq(host, method, path, query, body)
	if err != nil {
		return nil, err
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response, %s", err)
	}

	if resp.StatusCode > http.StatusIMUsed || resp.StatusCode < http.StatusOK {
		return nil, fmt.Errorf("http code is %d, body is %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}
