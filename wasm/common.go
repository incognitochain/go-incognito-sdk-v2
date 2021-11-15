package main

import "github.com/incognitochain/go-incognito-sdk-v2/incclient"

var ic *incclient.IncClient

func initClient(remoteHost string, network ...string) error {
	var err error
	ic, err = incclient.NewIncClient(remoteHost, incclient.MainNetETHHost, 2, network...)
	if err != nil {
		return err
	}

	return err
}
