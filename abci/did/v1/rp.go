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

package did

import (
	"github.com/gogo/protobuf/proto"
	"github.com/ndidplatform/smart-contract/abci/code"
	"github.com/ndidplatform/smart-contract/protos/data"
	pbParam "github.com/ndidplatform/smart-contract/protos/param"
	"github.com/tendermint/tendermint/abci/types"
)

func createRequest(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("CreateRequest, Parameter: %s", param)
	var funcParam pbParam.CreateRequestParam
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	var request data.Request

	// set request data
	request.RequestId = funcParam.RequestId
	request.MinIdp = int64(funcParam.MinIdp)
	request.MinAal = funcParam.MinAal
	request.MinIal = funcParam.MinIal
	request.RequestTimeout = int64(funcParam.RequestTimeout)
	// request.DataRequestList = funcParam.DataRequestList
	request.RequestMessageHash = funcParam.RequestMessageHash
	request.Mode = int64(funcParam.Mode)

	// set data request
	request.DataRequestList = make([]*data.DataRequest, 0)
	for index := range funcParam.DataRequestList {
		var newRow data.DataRequest
		newRow.ServiceId = funcParam.DataRequestList[index].ServiceId
		newRow.RequestParamsHash = funcParam.DataRequestList[index].RequestParamsHash
		newRow.MinAs = int64(funcParam.DataRequestList[index].MinAs)
		newRow.AsIdList = funcParam.DataRequestList[index].AsIdList
		if funcParam.DataRequestList[index].AsIdList == nil {
			newRow.AsIdList = make([]string, 0)
		}
		newRow.AnsweredAsIdList = make([]string, 0)
		newRow.ReceivedDataFromList = make([]string, 0)
		request.DataRequestList = append(request.DataRequestList, &newRow)
	}

	// set default value
	request.Closed = false
	request.TimedOut = false
	request.CanAddAccessor = false

	// set Owner
	request.Owner = nodeID

	// set Can add accossor
	ownerRole := getRoleFromNodeID(nodeID, app)
	if string(ownerRole) == "IdP" || string(ownerRole) == "MasterIdP" {
		request.CanAddAccessor = true
	}

	// set default value
	request.ResponseList = make([]*data.Response, 0)

	// check duplicate service ID in Data Request
	serviceIDCount := make(map[string]int)
	for _, dataRequest := range request.DataRequestList {
		serviceIDCount[dataRequest.ServiceId]++
	}
	for _, count := range serviceIDCount {
		if count > 1 {
			return ReturnDeliverTxLog(code.DuplicateServiceIDInDataRequest, "Duplicate Service ID In Data Request", "")
		}
	}

	key := "Request" + "|" + request.RequestId

	value, err := proto.Marshal(&request)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}

	_, existValue := app.state.db.Get(prefixKey([]byte(key)))
	if existValue != nil {
		return ReturnDeliverTxLog(code.DuplicateRequestID, "Duplicate Request ID", "")
	}
	app.SetStateDB([]byte(key), []byte(value))
	return ReturnDeliverTxLog(code.OK, "success", request.RequestId)
}

func closeRequest(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("CloseRequest, Parameter: %s", param)
	var funcParam pbParam.CloseRequestParam
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	key := "Request" + "|" + funcParam.RequestId
	_, value := app.state.db.Get(prefixKey([]byte(key)))

	if value == nil {
		return ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}

	var request data.Request
	err = proto.Unmarshal([]byte(value), &request)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	if request.Closed {
		return ReturnDeliverTxLog(code.RequestIsClosed, "Can not set time out a closed request", "")
	}

	if request.TimedOut {
		return ReturnDeliverTxLog(code.RequestIsTimedOut, "Can not close a timed out request", "")
	}

	// // Check valid list
	// if len(funcParam.ResponseValidList) != len(request.Responses) {
	// 	return ReturnDeliverTxLog(code.IncompleteValidList, "Incomplete valid list", "")
	// }

	for _, valid := range funcParam.ResponseValidList {
		for index := range request.ResponseList {
			if valid.IdpId == request.ResponseList[index].IdpId {
				if valid.ValidProof != nil {
					if valid.GetValidProofBool() {
						request.ResponseList[index].ValidProof = "true"
					} else {
						request.ResponseList[index].ValidProof = "false"
					}
				}
				if valid.ValidIal != nil {
					if valid.GetValidIalBool() {
						request.ResponseList[index].ValidIal = "true"
					} else {
						request.ResponseList[index].ValidIal = "false"
					}
				}
				if valid.ValidSignature != nil {
					if valid.GetValidSignatureBool() {
						request.ResponseList[index].ValidSignature = "true"
					} else {
						request.ResponseList[index].ValidSignature = "false"
					}
				}
			}
		}
	}

	request.Closed = true
	value, err = proto.Marshal(&request)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(key), []byte(value))
	return ReturnDeliverTxLog(code.OK, "success", funcParam.RequestId)
}

func timeOutRequest(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("TimeOutRequest, Parameter: %s", param)
	var funcParam pbParam.TimeOutRequestParam
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	key := "Request" + "|" + funcParam.RequestId
	_, value := app.state.db.Get(prefixKey([]byte(key)))

	if value == nil {
		return ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}

	var request data.Request
	err = proto.Unmarshal([]byte(value), &request)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	if request.TimedOut {
		return ReturnDeliverTxLog(code.RequestIsTimedOut, "Can not close a timed out request", "")
	}

	if request.Closed {
		return ReturnDeliverTxLog(code.RequestIsClosed, "Can not set time out a closed request", "")
	}

	// // Check valid list
	// if len(funcParam.ResponseValidList) != len(request.Responses) {
	// 	return ReturnDeliverTxLog(code.IncompleteValidList, "Incomplete valid list", "")
	// }

	for _, valid := range funcParam.ResponseValidList {
		for index := range request.ResponseList {
			if valid.IdpId == request.ResponseList[index].IdpId {
				if valid.ValidProof != nil {
					if valid.GetValidProofBool() {
						request.ResponseList[index].ValidProof = "true"
					} else {
						request.ResponseList[index].ValidProof = "false"
					}
				}
				if valid.ValidIal != nil {
					if valid.GetValidIalBool() {
						request.ResponseList[index].ValidIal = "true"
					} else {
						request.ResponseList[index].ValidIal = "false"
					}
				}
				if valid.ValidSignature != nil {
					if valid.GetValidSignatureBool() {
						request.ResponseList[index].ValidSignature = "true"
					} else {
						request.ResponseList[index].ValidSignature = "false"
					}
				}
			}
		}
	}

	request.TimedOut = true
	value, err = proto.Marshal(&request)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}

	app.SetStateDB([]byte(key), []byte(value))
	return ReturnDeliverTxLog(code.OK, "success", funcParam.RequestId)
}

func setDataReceived(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("SetDataReceived, Parameter: %s", param)
	var funcParam pbParam.SetDataReceivedParam
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	key := "Request" + "|" + funcParam.RequestId
	_, value := app.state.db.Get(prefixKey([]byte(key)))

	if value == nil {
		return ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}

	var request data.Request
	err = proto.Unmarshal([]byte(value), &request)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check as_id is exist in as_id_list
	exist := false
	for _, dataRequest := range request.DataRequestList {
		if dataRequest.ServiceId == funcParam.ServiceId {
			for _, as := range dataRequest.AnsweredAsIdList {
				if as == funcParam.AsId {
					exist = true
					break
				}
			}
		}
	}
	if exist == false {
		return ReturnDeliverTxLog(code.AsIDIsNotExistInASList, "AS ID is not exist in answered AS list", "")
	}

	// Check Duplicate AS ID
	duplicate := false
	for _, dataRequest := range request.DataRequestList {
		if dataRequest.ServiceId == funcParam.ServiceId {
			for _, as := range dataRequest.ReceivedDataFromList {
				if as == funcParam.AsId {
					duplicate = true
					break
				}
			}
		}
	}
	if duplicate == true {
		return ReturnDeliverTxLog(code.DuplicateASInDataRequest, "Duplicate AS ID in data request", "")
	}

	// Update received_data_from_list in request
	for index, dataRequest := range request.DataRequestList {
		if dataRequest.ServiceId == funcParam.ServiceId {
			request.DataRequestList[index].ReceivedDataFromList = append(dataRequest.ReceivedDataFromList, funcParam.AsId)
		}
	}

	// Request has data request. If received data, signed answer > data request count on each data request
	// dataRequestCompletedCount := 0
	// for _, dataRequest := range request.DataRequestList {
	// 	if len(dataRequest.AnsweredAsIdList) >= dataRequest.Count &&
	// 		len(dataRequest.ReceivedDataFromList) >= dataRequest.Count {
	// 		dataRequestCompletedCount++
	// 	}
	// }
	// if dataRequestCompletedCount == len(request.DataRequestList) {
	// 	app.logger.Info("Auto close")
	// 	request.IsClosed = true
	// } else {
	// 	app.logger.Info("Auto close")
	// }

	value, err = proto.Marshal(&request)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(key), []byte(value))
	return ReturnDeliverTxLog(code.OK, "success", funcParam.RequestId)
}
