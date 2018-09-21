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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	pbParam "github.com/ndidplatform/smart-contract/protos/params"
	"github.com/tendermint/tendermint/libs/common"
)

func InitNDID(t *testing.T) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidpublicKeyBytes, err := generatePublicKey(&ndidKey.PublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}
	var initNDIDparam pbParam.InitNDIDParams
	initNDIDparam.NodeId = "NDID"
	initNDIDparam.PublicKey = string(ndidpublicKeyBytes)
	initNDIDparam.MasterPublicKey = string(ndidpublicKeyBytes)

	paramJSON, err := proto.Marshal(&initNDIDparam)
	if err != nil {
		fmt.Println("error:", err)
	}

	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "InitNDID"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	startTime := time.Now()
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(initNDIDparam.NodeId))
	resultObj, _ := result.(ResponseTx)
	if resultObj.Result.CheckTx.Log == "NDID node is already existed" {
		t.SkipNow()
	}
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	stopTime := time.Now()
	writeLog(fnName, (stopTime.UnixNano()-startTime.UnixNano())/int64(time.Millisecond))
	t.Logf("PASS: %s", fnName)
}

func RegisterNode(t *testing.T, param pbParam.RegisterNodeParams) {
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}

	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := []byte("NDID")

	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "RegisterNode"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, ndidNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func SetTimeOutBlockRegisterMsqDestination(t *testing.T) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	var param pbParam.TimeOutBlockRegisterMsqDestinationParams
	param.TimeOutBlock = 100
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "SetTimeOutBlockRegisterMsqDestination"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte("NDID"))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func AddNodeToken(t *testing.T, param pbParam.AddNodeTokenParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "AddNodeToken"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func ReduceNodeToken(t *testing.T, param pbParam.ReduceNodeTokenParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "ReduceNodeToken"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func SetNodeToken(t *testing.T, param pbParam.SetNodeTokenParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "SetNodeToken"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func AddService(t *testing.T, param pbParam.AddServiceParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "AddService"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func DisableService(t *testing.T, param pbParam.DisableServiceParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "DisableService"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func EnableService(t *testing.T, param pbParam.DisableServiceParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "EnableService"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func RegisterServiceDestinationByNDID(t *testing.T, param pbParam.RegisterServiceDestinationByNDIDParams) {
	key := getPrivateKeyFromString(ndidPrivK)
	nodeID := []byte("NDID")
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		log.Fatal(err.Error())
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "RegisterServiceDestinationByNDID"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, nodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func UpdateService(t *testing.T, param pbParam.UpdateServiceParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "UpdateService"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func SetPriceFunc(t *testing.T, param pbParam.SetPriceFuncParams) {
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := []byte("NDID")
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "SetPriceFunc"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, ndidNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func AddNamespace(t *testing.T, param pbParam.AddNamespaceParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	nodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "AddNamespace"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(nodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func DisableNamespace(t *testing.T, param pbParam.DisableNamespaceParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	nodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "DisableNamespace"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(nodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func EnableNamespace(t *testing.T, param pbParam.DisableNamespaceParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	nodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "EnableNamespace"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(nodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func SetValidator(t *testing.T, param pbParam.SetValidatorParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "SetValidator"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func UpdateNodeByNDID(t *testing.T, param pbParam.UpdateNodeByNDIDParams) {
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := []byte("NDID")
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "UpdateNodeByNDID"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, ndidNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func DisableNode(t *testing.T, param pbParam.DisableNodeParams) {
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := []byte("NDID")
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "DisableNode"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, ndidNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func EnableNode(t *testing.T, param pbParam.DisableNodeParams) {
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := []byte("NDID")
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "EnableNode"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, ndidNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func DisableServiceDestinationByNDID(t *testing.T, param pbParam.DisableServiceDestinationByNDIDParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "DisableServiceDestinationByNDID"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func EnableServiceDestinationByNDID(t *testing.T, param pbParam.DisableServiceDestinationByNDIDParams) {
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := "NDID"
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "EnableServiceDestinationByNDID"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, []byte(ndidNodeID))
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func AddNodeToProxyNode(t *testing.T, param pbParam.AddNodeToProxyNodeParams, expected string) {
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}

	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := []byte("NDID")

	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "AddNodeToProxyNode"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, ndidNodeID)
	resultObj, _ := result.(ResponseTx)
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func UpdateNodeProxyNode(t *testing.T, param pbParam.UpdateNodeProxyNodeParams, expected string) {
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := []byte("NDID")
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "UpdateNodeProxyNode"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, ndidNodeID)
	resultObj, _ := result.(ResponseTx)
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func RemoveNodeFromProxyNode(t *testing.T, param pbParam.RemoveNodeFromProxyNodeParams, expected string) {
	paramJSON, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	ndidKey := getPrivateKeyFromString(ndidPrivK)
	ndidNodeID := []byte("NDID")
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	fnName := "RemoveNodeFromProxyNode"
	tempPSSmessage := append([]byte(fnName), paramJSON...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, ndidKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramJSON, []byte(nonce), signature, ndidNodeID)
	resultObj, _ := result.(ResponseTx)
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}
