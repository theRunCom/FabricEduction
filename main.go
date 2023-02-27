package main

import (
	"encoding/json"
	"fmt"
	"FabricEduction/sdkInit"
	"FabricEduction/service"
	"FabricEduction/web"
	"FabricEduction/web/controller"
	"os"
)

const (
	edu_name = "edu"
	edu_version = "1.0.0"
)

func main() {
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    1,
			OrgAnchorFile: os.Getenv("GOPATH") + "/src/FabricEduction/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
	}
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    os.Getenv("GOPATH") + "/src/education/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      edu_name,
		ChaincodePath:    os.Getenv("GOPATH")+"/src/education/chaincode/",
		ChaincodeVersion: edu_version,
	}
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}
	fmt.Println(">> 通过链码外部服务设置链码状态......")
	edu := service.Education{
		Name: "张三",
		Gender: "男",
		Nation: "汉",
		EntityID: "101",
		Place: "北京",
		BirthDay: "1991年01月01日",
		EnrollDate: "2009年9月",
		GraduationDate: "2013年7月",
		SchoolName: "中国政法大学",
		Major: "社会学",
		QuaType: "普通",
		Length: "四年",
		Mode: "普通全日制",
		Level: "本科",
		Graduation: "毕业",
		CertNo: "111",
		Photo: "/static/photo/11.png",
	}
	serviceSetup, err := service.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk)
	if err!=nil{
		fmt.Println()
		os.Exit(-1)
	}
	msg, err := serviceSetup.SaveEdu(edu)
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("信息发布成功, 交易编号为: " + msg)
	}
	result, err := serviceSetup.FindEduInfoByEntityID("101")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		var edu service.Education
		json.Unmarshal(result, &edu)
		fmt.Println("根据身份证号码查询信息成功：")
		fmt.Println(edu)
	}
	app := controller.Application{
		Setup: serviceSetup,
	}
	web.WebStart(app)
}