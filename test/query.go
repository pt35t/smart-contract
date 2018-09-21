/**
 * Copyright (c) 2018, 2019 National Digital ID COMPANY LIMITED
 *
 * This file is part of NDID software.
 *
 * NDID is the free software: you can redistribute it and/or modify it under
 * the terms of the Affero GNU General Public License as published by the
 * Free Software Foundation, either version 3 of the License, or any later
 * version.
 *
 * NDID is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
 * See the Affero GNU General Public License for more details.
 *
 * You should have received a copy of the Affero GNU General Public License
 * along with the NDID source code. If not, see https://www.gnu.org/licenses/agpl.txt.
 *
 * Please contact info@ndid.co.th for any further questions
 *
 */

package test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/gogo/protobuf/proto"
	did "github.com/ndidplatform/smart-contract/abci/did/v1"
	pbParam "github.com/ndidplatform/smart-contract/protos/params"
	pbResult "github.com/ndidplatform/smart-contract/protos/result"
)

func GetNodePublicKey(t *testing.T, param pbParam.GetNodePublicKeyParams, expected string) {
	fnName := "GetNodePublicKey"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res pbResult.GetNodePublicKeyResult
	err = proto.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := res.PublicKey; actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetNodeMasterPublicKey(t *testing.T, param pbParam.GetNodeMasterPublicKeyParams, expected string) {
	fnName := "GetNodeMasterPublicKey"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res pbResult.GetNodeMasterPublicKeyResult
	err = proto.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := res.MasterPublicKey; actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetNodeToken(t *testing.T, param pbParam.GetNodeTokenParams, expected pbResult.GetNodeTokenResult) {
	fnName := "GetNodeToken"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res pbResult.GetNodeTokenResult
	err = proto.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetNodeTokenExpectString(t *testing.T, param pbParam.GetNodeTokenParams, expected string) {
	fnName := "GetNodeToken"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res pbResult.GetNodeTokenResult
	err = proto.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetIdpNodes(t *testing.T, param pbParam.GetIdpNodesParams, expected pbResult.GetIdpNodesResult) {
	fnName := "GetIdpNodes"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res pbResult.GetIdpNodesResult
	err = proto.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetIdpNodesForDisable(t *testing.T, param pbParam.GetIdpNodesParams) []*pbResult.MsqDestinationNode {
	fnName := "GetIdpNodes"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res pbResult.GetIdpNodesResult
	err = proto.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	return res.Node
}

// func GetIdpNodesExpectString(t *testing.T, param pbParam.GetIdpNodesParams, expected string) {
// 	fnName := "GetIdpNodes"
// 	paramsByte, err := proto.Marshal(&param)
// 	if err != nil {
// 		fmt.Println("error:", err)
// 	}
// 	result, _ := queryTendermint([]byte(fnName), paramsByte)
// 	resultObj, _ := result.(ResponseQuery)
// 	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
// 	if resultObj.Result.Response.Log == expected {
// 		t.Logf("PASS: %s", fnName)
// 		return
// 	}
// 	if actual := string(resultString); actual != expected {
// 		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
// 	}
// 	t.Logf("PASS: %s", fnName)
// }

func GetMsqAddress(t *testing.T, param pbParam.GetMsqAddressParams, expected []did.MsqAddress) {
	fnName := "GetMsqAddress"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res []did.MsqAddress
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetMsqAddressExpectString(t *testing.T, param pbParam.GetMsqAddressParams, expected string) {
	fnName := "GetMsqAddress"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res []did.MsqAddress
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetServiceDetail(t *testing.T, param pbParam.GetServiceDetailParams, expected did.ServiceDetail) {
	fnName := "GetServiceDetail"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		log.Fatal(err.Error())
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res did.ServiceDetail
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetServiceDetailExpectString(t *testing.T, param pbParam.GetServiceDetailParams, expected string) {
	fnName := "GetServiceDetail"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		log.Fatal(err.Error())
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res did.ServiceDetail
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetAsNodesByServiceId(t *testing.T, param pbParam.GetAsNodesByServiceIdParams, expected string) {
	fnName := "GetAsNodesByServiceId"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		log.Fatal(err.Error())
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res did.GetAsNodesByServiceIdResult
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetIdentityProof(t *testing.T, param pbParam.GetIdentityProofParams, expected did.GetIdentityProofResult) {
	fnName := "GetIdentityProof"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res did.GetIdentityProofResult
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetIdentityProofExpectString(t *testing.T, param pbParam.GetIdentityProofParams, expected string) {
	fnName := "GetIdentityProof"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res did.GetIdentityProofResult
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetPriceFunc(t *testing.T, param pbParam.GetPriceFuncParams, expected pbResult.GetPriceFuncResult) {
	fnName := "GetPriceFunc"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)

	var res pbResult.GetPriceFuncResult
	err = proto.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetUsedTokenReport(t *testing.T, param pbParam.GetUsedTokenReportParams, expected pbResult.GetUsedTokenReportResult) {
	fnName := "GetUsedTokenReport"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res pbResult.GetUsedTokenReportResult
	err = proto.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetRequestDetail(t *testing.T, param pbParam.GetRequestParams, expected string) {
	fnName := "GetRequestDetail"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetRequest(t *testing.T, param pbParam.GetRequestParams, expected did.GetRequestResult) {
	fnName := "GetRequest"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res did.GetRequestResult
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetRequestExpectString(t *testing.T, param pbParam.GetRequestParams, expected string) {
	fnName := "GetRequest"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res did.GetRequestResult
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetNamespaceList(t *testing.T, expected []did.Namespace) {
	fnName := "GetNamespaceList"
	paramsByte := []byte("")
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res []did.Namespace
	err := json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetNamespaceListForDisable(t *testing.T) (expected []did.Namespace) {
	fnName := "GetNamespaceList"
	paramsByte := []byte("")
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res []did.Namespace
	err := json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	return res
}

func GetNamespaceListExpectString(t *testing.T, expected string) {
	fnName := "GetNamespaceList"
	paramsByte := []byte("")
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if actual := string(resultString); actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func CheckExistingIdentity(t *testing.T, param pbParam.CheckExistingIdentityParams, expected string) {
	fnName := "CheckExistingIdentity"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if actual := string(resultString); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetAccessorGroupID(t *testing.T, param pbParam.GetAccessorGroupIDParams, expected string) {
	fnName := "GetAccessorGroupID"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetAccessorKey(t *testing.T, param pbParam.GetAccessorKeyParams, expected string) {
	fnName := "GetAccessorKey"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetServiceList(t *testing.T, expected []did.ServiceDetail) {
	fnName := "GetServiceList"
	paramsByte := []byte("")
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res []did.ServiceDetail
	err := json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if actual := res; !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetServiceListForDisable(t *testing.T) (expected []did.ServiceDetail) {
	fnName := "GetServiceList"
	paramsByte := []byte("")
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res []did.ServiceDetail
	err := json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	return res
}

func GetNodeInfo(t *testing.T, param pbParam.GetNodeInfoParams, expected string) {
	fnName := "GetNodeInfo"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func CheckExistingAccessorID(t *testing.T, param pbParam.CheckExistingAccessorIDParams, expected string) {
	fnName := "CheckExistingAccessorID"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if actual := string(resultString); actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func CheckExistingAccessorGroupID(t *testing.T, param pbParam.CheckExistingAccessorGroupIDParams, expected string) {
	fnName := "CheckExistingAccessorGroupID"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if actual := string(resultString); actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetIdentityInfo(t *testing.T, param pbParam.GetIdentityInfoParams, expected string) {
	fnName := "GetIdentityInfo"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetServicesByAsID(t *testing.T, param pbParam.GetServicesByAsIDParams, expected string) {
	fnName := "GetServicesByAsID"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		log.Fatal(err.Error())
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetDataSignature(t *testing.T, param pbParam.GetDataSignatureParams, expected string) {
	fnName := "GetDataSignature"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetIdpNodesInfo(t *testing.T, param pbParam.GetIdpNodesParams, expected string) {
	fnName := "GetIdpNodesInfo"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if actual := string(resultString); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetAsNodesInfoByServiceId(t *testing.T, param pbParam.GetAsNodesByServiceIdParams, expected string) {
	fnName := "GetAsNodesInfoByServiceId"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		log.Fatal(err.Error())
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res did.GetAsNodesByServiceIdResult
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetIdpNodesInfoParamJSON(t *testing.T, paramsByte string, expected string) {
	fnName := "GetIdpNodesInfo"
	result, _ := queryTendermint([]byte(fnName), []byte(paramsByte))
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if actual := string(resultString); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetNodesBehindProxyNode(t *testing.T, param pbParam.GetNodesBehindProxyNodeParams, expected string) {
	fnName := "GetNodesBehindProxyNode"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetNodeIDList(t *testing.T, param pbParam.GetNodeIDListParams, expected string) {
	fnName := "GetNodeIDList"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	if resultObj.Result.Response.Log == expected {
		t.Logf("PASS: %s", fnName)
		return
	}
	if actual := string(resultString); actual != expected {
		t.Fatalf("FAIL: %s\nExpected: %s\nActual: %s", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func GetNodeIDListForDisable(t *testing.T, param pbParam.GetNodeIDListParams) []string {
	fnName := "GetNodeIDList"
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	result, _ := queryTendermint([]byte(fnName), paramsByte)
	resultObj, _ := result.(ResponseQuery)
	resultString, _ := base64.StdEncoding.DecodeString(resultObj.Result.Response.Value)
	var res did.GetNodeIDListResult
	err = json.Unmarshal(resultString, &res)
	if err != nil {
		log.Fatal(err.Error())
	}
	return res.NodeIDList
}
