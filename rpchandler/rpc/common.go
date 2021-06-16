package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RPCServer struct {
	url string
}

func NewRPCServer(url string) *RPCServer {
	return &RPCServer{url: url}
}

func (server RPCServer) GetURL() string {
	return server.url
}

func (server *RPCServer) SendPostRequestWithQuery(query string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return nil, fmt.Errorf("rpc server not set")
	}

	var jsonStr = []byte(query)
	req, _ := http.NewRequest("POST", server.url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, err
		}
		return body, nil
	}
}
