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
	"testing"

	"github.com/gogo/protobuf/proto"
	pbParam "github.com/ndidplatform/smart-contract/protos/params"
	"github.com/tendermint/tendermint/libs/common"
)

func RegisterMsqDestination(t *testing.T, param pbParam.RegisterMsqDestinationParams, privKeyFile string, nodeID string, expected string) {
	idpKey := getPrivateKeyFromString(privKeyFile)
	idpNodeID := []byte(nodeID)
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	fnName := "RegisterMsqDestination"
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	tempPSSmessage := append([]byte(fnName), paramsByte...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, idpKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramsByte, []byte(nonce), signature, idpNodeID)
	resultObj, _ := result.(ResponseTx)
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func DeclareIdentityProof(t *testing.T, param pbParam.DeclareIdentityProofParams, privKeyFile string, nodeID string) {
	idpKey := getPrivateKeyFromString(privKeyFile)
	idpNodeID := []byte(nodeID)
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	fnName := "DeclareIdentityProof"
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	tempPSSmessage := append([]byte(fnName), paramsByte...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, idpKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramsByte, []byte(nonce), signature, idpNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func CreateIdpResponse(t *testing.T, param pbParam.CreateIdpResponseParams, privKeyFile string, nodeID string) {
	idpKey := getPrivateKeyFromString(privKeyFile)
	idpNodeID := []byte(nodeID)
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	fnName := "CreateIdpResponse"
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	tempPSSmessage := append([]byte(fnName), paramsByte...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, idpKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramsByte, []byte(nonce), signature, idpNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func CreateIdentity(t *testing.T, param pbParam.CreateIdentityParams, nodeID string) {
	idpKey := getPrivateKeyFromString(idpPrivK)
	idpNodeID := []byte(nodeID)
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	fnName := "CreateIdentity"
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	tempPSSmessage := append([]byte(fnName), paramsByte...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, idpKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramsByte, []byte(nonce), signature, idpNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func AddAccessorMethod(t *testing.T, param pbParam.AccessorMethodParams, nodeID string) {
	idpKey := getPrivateKeyFromString(idpPrivK)
	idpNodeID := []byte(nodeID)
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	fnName := "AddAccessorMethod"
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	tempPSSmessage := append([]byte(fnName), paramsByte...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, idpKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramsByte, []byte(nonce), signature, idpNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func ClearRegisterMsqDestinationTimeout(t *testing.T, param pbParam.ClearRegisterMsqDestinationTimeoutParams, privKeyFile string, nodeID string) {
	idpKey := getPrivateKeyFromString(privKeyFile)
	idpNodeID := []byte(nodeID)
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	fnName := "ClearRegisterMsqDestinationTimeout"
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	tempPSSmessage := append([]byte(fnName), paramsByte...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, idpKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramsByte, []byte(nonce), signature, idpNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}

func UpdateIdentity(t *testing.T, param pbParam.UpdateIdentityParams, nodeID string) {
	idpKey := getPrivateKeyFromString(idpPrivK2)
	idpNodeID := []byte(nodeID)
	paramsByte, err := proto.Marshal(&param)
	if err != nil {
		fmt.Println("error:", err)
	}
	fnName := "UpdateIdentity"
	nonce := base64.StdEncoding.EncodeToString([]byte(common.RandStr(12)))
	tempPSSmessage := append([]byte(fnName), paramsByte...)
	tempPSSmessage = append(tempPSSmessage, []byte(nonce)...)
	PSSmessage := []byte(base64.StdEncoding.EncodeToString(tempPSSmessage))
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, idpKey, newhash, hashed)
	result, _ := callTendermint([]byte(fnName), paramsByte, []byte(nonce), signature, idpNodeID)
	resultObj, _ := result.(ResponseTx)
	expected := "success"
	if actual := resultObj.Result.DeliverTx.Log; actual != expected {
		t.Errorf("\n"+`CheckTx log: "%s"`, resultObj.Result.CheckTx.Log)
		t.Fatalf("FAIL: %s\nExpected: %#v\nActual: %#v", fnName, expected, actual)
	}
	t.Logf("PASS: %s", fnName)
}
