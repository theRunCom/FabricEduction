
package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"time"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"FabricEduction/sdkInit"
)

type Education struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`		
	Gender	string	`json:"Gender"`		
	Nation	string	`json:"Nation"`		
	EntityID	string	`json:"EntityID"`		
	Place	string	`json:"Place"`		
	BirthDay	string	`json:"BirthDay"`		
	EnrollDate	string	`json:"EnrollDate"`		
	GraduationDate	string	`json:"GraduationDate"`	
	SchoolName	string	`json:"SchoolName"`	
	Major	string	`json:"Major"`	
	QuaType	string	`json:"QuaType"`	
	Length	string	`json:"Length"`	
	Mode	string	`json:"Mode"`	
	Level	string	`json:"Level"`	
	Graduation	string	`json:"Graduation"`	
	CertNo	string	`json:"CertNo"`	
	Photo	string	`json:"Photo"`	
	Historys	[]HistoryItem	
}

type HistoryItem struct {
	TxId	string
	Education	Education
}

type ServiceSetup struct {
	ChaincodeID	string
	Client	*channel.Client
}

func regitserEvent(client *channel.Client, chaincodeID, eventID string) (fab.Registration, <-chan *fab.CCEvent) {

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventID)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}
	return reg, notifier
}

func eventResult(notifier <-chan *fab.CCEvent, eventID string) error {
	select {
	case ccEvent := <-notifier:
		fmt.Printf("接收到链码事件: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("不能根据指定的事件ID接收到相应的链码事件(%s)", eventID)
	}
	return nil
}

func InitService(chaincodeID, channelID string, org *sdkInit.OrgInfo, sdk *fabsdk.FabricSDK) (*ServiceSetup, error) {
	handler := &ServiceSetup{
		ChaincodeID:chaincodeID,
	}
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	client, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new channel client: %s", err)
	}
	handler.Client = client
	return handler, nil
}