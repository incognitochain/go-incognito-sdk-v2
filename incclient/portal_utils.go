package incclient

import (
	"github.com/incognitochain/go-incognito-sdk-v2/metadata"
	"github.com/incognitochain/go-incognito-sdk-v2/rpchandler"
)

// GetPortalUnShieldingRequestStatus retrieves the status of a port un-shielding request.
func (client *IncClient) GetPortalUnShieldingRequestStatus(unShieldID string) (*metadata.PortalUnshieldRequestStatus, error) {
	responseInBytes, err := client.rpcServer.GetPortalUnShieldingRequestStatus(unShieldID)
	if err != nil {
		return nil, err
	}

	var res *metadata.PortalUnshieldRequestStatus
	err = rpchandler.ParseResponse(responseInBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}