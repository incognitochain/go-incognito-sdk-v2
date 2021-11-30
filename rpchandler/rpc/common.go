package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// RPCServer represents a RPC host server.
type RPCServer struct {
	url string
}

// NewRPCServer creates a new RPCServer pointing to the given url.
func NewRPCServer(url string) *RPCServer {
	return &RPCServer{url: url}
}

// GetURL returns the url of a RPCServer.
func (server *RPCServer) GetURL() string {
	return server.url
}

// InitToURL points a RPCServer to a given url.
func (server *RPCServer) InitToURL(url string) *RPCServer {
	server.url = url
	return server
}

// SendQuery sends a query to the remote server given the method and parameters.
func (server *RPCServer) SendQuery(method string, params []interface{}) ([]byte, error) {
	if params == nil {
		params = make([]interface{}, 0)
	}
	request := rpchandler.CreateJsonRequest("1.0", method, params, 1)

	query, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return server.SendPostRequestWithQuery(string(query))
}

// SendPostRequestWithQuery sends a query to the remote server using the POST method.
func (server *RPCServer) SendPostRequestWithQuery(query string) ([]byte, error) {
	if server == nil || len(server.url) == 0 {
		return []byte{}, fmt.Errorf("server has not been set")
	}
	var jsonStr = []byte(query)
	req, _ := http.NewRequest("POST", server.url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client.Timeout = 10 * time.Minute
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("DoReq %v\n", err)
		return []byte{}, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%v", resp.Status)
	} else {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Printf("BodyClose %v\n", err)
			}
		}()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ReadAll %v\n", err)
			return []byte{}, err
		}
		return body, nil
	}
}
