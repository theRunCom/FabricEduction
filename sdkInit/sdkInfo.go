package sdkInit

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
)

type OrgInfo struct {
	OrgAdminUser          string 
	OrgName               string 
	OrgMspId              string 
	OrgUser               string 
	orgMspClient          *mspclient.Client
	OrgAdminClientContext *contextAPI.ClientProvider
	OrgResMgmt            *resmgmt.Client
	OrgPeerNum            int
	OrgAnchorFile string 
}

type SdkEnvInfo struct {
	ChannelID     string 
	ChannelConfig string 
	Orgs []*OrgInfo
	OrdererAdminUser     string 
	OrdererOrgName       string 
	OrdererEndpoint      string
	OrdererClientContext *contextAPI.ClientProvider
	ChaincodeID      string
	ChaincodeGoPath  string
	ChaincodePath    string
	ChaincodeVersion string
}




