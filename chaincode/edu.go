package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"fmt"
	"encoding/json"
	"bytes"
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

type EducationChaincode struct {

}

func (t *EducationChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")
	return shim.Success(nil)
}

func (t *EducationChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "addEdu" {
		return t.addEdu(stub, args)	
	} else if fn == "queryEduByCertNoAndName" {
		return t.queryEduByCertNoAndName(stub, args)		
	} else if fn == "queryEduInfoByEntityID" {
		return t.queryEduInfoByEntityID(stub, args)
	} else if fn == "updateEdu" {
		return t.updateEdu(stub, args)		
	} else if fn == "delEdu" {
		return t.delEdu(stub, args)	
	}
	return shim.Error("Invoke fn error")
}

const DOC_TYPE = "eduObj"

func PutEdu(stub shim.ChaincodeStubInterface, edu Education) ([]byte, bool) {
	edu.ObjectType = DOC_TYPE
	b, err := json.Marshal(edu)
	if err != nil {
		return nil, false
	}

	err = stub.PutState(edu.EntityID, b)
	if err != nil {
		return nil, false
	}
	return b, true
}

func GetEduInfo(stub shim.ChaincodeStubInterface, entityID string) (Education, bool)  {
	var edu Education
	b, err := stub.GetState(entityID)
	if err != nil {
		return edu, false
	}
	if b == nil {
		return edu, false
	}
	err = json.Unmarshal(b, &edu)
	if err != nil {
		return edu, false
	}
	return edu, true
}

func getEduByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer  resultsIterator.Close()
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func (t *EducationChaincode) addEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}
	var edu Education
	err := json.Unmarshal([]byte(args[0]), &edu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}
	_, exist := GetEduInfo(stub, edu.EntityID)
	if exist {
		return shim.Error("要添加的身份证号码已存在")
	}
	_, bl := PutEdu(stub, edu)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("信息添加成功"))
}

func (t *EducationChaincode) queryEduByCertNoAndName(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	CertNo := args[0]
	name := args[1]
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"CertNo\":\"%s\", \"Name\":\"%s\"}}", DOC_TYPE, CertNo, name)
	result, err := getEduByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据证书编号及姓名查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的证书编号及姓名没有查询到相关的信息")
	}
	return shim.Success(result)
}

func (t *EducationChaincode) queryEduInfoByEntityID(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	b, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("根据身份证号码查询信息失败")
	}
	if b == nil {
		return shim.Error("根据身份证号码没有查询到相关的信息")
	}
	var edu Education
	err = json.Unmarshal(b, &edu)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}
	iterator, err := stub.GetHistoryForKey(edu.EntityID)
	if err != nil {
		return shim.Error("根据指定的身份证号码查询对应的历史变更数据失败")
	}
	defer iterator.Close()
	var historys []HistoryItem
	var hisEdu Education
	for iterator.HasNext() {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取edu的历史变更数据失败")
		}
		var historyItem HistoryItem
		historyItem.TxId = hisData.TxId
		json.Unmarshal(hisData.Value, &hisEdu)
		if hisData.Value == nil {
			var empty Education
			historyItem.Education = empty
		}else {
			historyItem.Education = hisEdu
		}
		historys = append(historys, historyItem)
	}
	edu.Historys = historys
	result, err := json.Marshal(edu)
	if err != nil {
		return shim.Error("序列化edu信息时发生错误")
	}
	return shim.Success(result)
}

func (t *EducationChaincode) updateEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}
	var info Education
	err := json.Unmarshal([]byte(args[0]), &info)
	if err != nil {
		return  shim.Error("反序列化edu信息失败")
	}
	result, bl := GetEduInfo(stub, info.EntityID)
	if !bl{
		return shim.Error("根据身份证号码查询信息时发生错误")
	}
	result.Name = info.Name
	result.BirthDay = info.BirthDay
	result.Nation = info.Nation
	result.Gender = info.Gender
	result.Place = info.Place
	result.EntityID = info.EntityID
	result.Photo = info.Photo
	result.EnrollDate = info.EnrollDate
	result.GraduationDate = info.GraduationDate
	result.SchoolName = info.SchoolName
	result.Major = info.Major
	result.QuaType = info.QuaType
	result.Length = info.Length
	result.Mode = info.Mode
	result.Level = info.Level
	result.Graduation = info.Graduation
	result.CertNo = info.CertNo;
	_, bl = PutEdu(stub, result)
	if !bl {
		return shim.Error("保存信息信息时发生错误")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("信息更新成功"))
}

func (t *EducationChaincode) delEdu(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2{
		return shim.Error("给定的参数个数不符合要求")
	}
	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("删除信息时发生错误")
	}
	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("信息删除成功"))
}

func main() {
	err := shim.Start(new(EducationChaincode))
	if err != nil{
		fmt.Println("Start Error!")
	}
}

