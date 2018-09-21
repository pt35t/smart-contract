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
	"encoding/json"

	"github.com/gogo/protobuf/proto"
	"github.com/ndidplatform/smart-contract/abci/code"
	"github.com/ndidplatform/smart-contract/protos/data"
	pbParam "github.com/ndidplatform/smart-contract/protos/params"
	"github.com/tendermint/tendermint/abci/types"
)

func createIdentity(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("CreateIdentity, Parameter: %s", param)
	var funcParam pbParam.CreateIdentityParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	accessorKey := "Accessor" + "|" + funcParam.AccessorId
	var accessor data.Accessor
	accessor.AccessorType = funcParam.AccessorType
	accessor.AccessorPublicKey = funcParam.AccessorPublicKey
	accessor.AccessorGroupId = funcParam.AccessorGroupId
	accessor.Active = true
	accessor.Owner = nodeID

	accessorJSON, err := proto.Marshal(&accessor)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}

	accessorGroupKey := "AccessorGroup" + "|" + funcParam.AccessorGroupId
	accessorGroup := funcParam.AccessorGroupId

	// Check duplicate accessor_id
	_, chkAccessorKeyExists := app.state.db.Get(prefixKey([]byte(accessorKey)))
	if chkAccessorKeyExists != nil {
		return ReturnDeliverTxLog(code.DuplicateAccessorID, "Duplicate Accessor ID", "")
	}

	// Check duplicate accessor_group_id
	_, chkAccessorGroupKeyExists := app.state.db.Get(prefixKey([]byte(accessorGroupKey)))
	if chkAccessorGroupKeyExists != nil {
		return ReturnDeliverTxLog(code.DuplicateAccessorGroupID, "Duplicate Accessor Group ID", "")
	}

	app.SetStateDB([]byte(accessorKey), []byte(accessorJSON))
	app.SetStateDB([]byte(accessorGroupKey), []byte(accessorGroup))

	return ReturnDeliverTxLog(code.OK, "success", "")
}

func setCanAddAccessorToFalse(requestID string, app *DIDApplication) {
	key := "Request" + "|" + requestID
	_, value := app.state.db.Get(prefixKey([]byte(key)))
	if value != nil {
		var request data.Request
		err := proto.Unmarshal([]byte(value), &request)
		if err == nil {
			request.CanAddAccessor = false
			value, err := proto.Marshal(&request)
			if err == nil {
				app.SetStateDB([]byte(key), []byte(value))
			}
		}
	}
}

func addAccessorMethod(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("AddAccessorMethod, Parameter: %s", param)
	var funcParam pbParam.AccessorMethodParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// AccessorGroupID: must already exist
	accessorGroupKey := "AccessorGroup" + "|" + funcParam.AccessorGroupId
	_, chkAccessorGroupKeyExists := app.state.db.Get(prefixKey([]byte(accessorGroupKey)))
	if chkAccessorGroupKeyExists == nil {
		return ReturnDeliverTxLog(code.AccessorGroupIDNotFound, "Accessor Group ID not found", "")
	}

	// AccessorID: must not duplicate
	accessorKey := "Accessor" + "|" + funcParam.AccessorId
	_, chkAccessorKeyExists := app.state.db.Get(prefixKey([]byte(accessorKey)))
	if chkAccessorKeyExists != nil {
		return ReturnDeliverTxLog(code.DuplicateAccessorID, "Duplicate Accessor ID", "")
	}

	// Request must be completed, can be used only once, special type
	var getRequestparam GetRequestParam
	getRequestparam.RequestID = funcParam.RequestId
	getRequestparamJSON, err := json.Marshal(getRequestparam)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	var request = getRequest(getRequestparamJSON, app, app.state.db.Version64())
	var requestDetail = getRequestDetail(getRequestparamJSON, app, app.state.db.Version64())
	var requestResult GetRequestResult
	var requestDetailResult GetRequestDetailResult
	err = json.Unmarshal([]byte(request.Value), &requestResult)
	if err != nil {
		return ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	err = json.Unmarshal([]byte(requestDetail.Value), &requestDetailResult)
	if err != nil {
		return ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}

	// Check accept result >= min_idp
	acceptCount := 0
	for _, response := range requestDetailResult.Responses {
		if response.Status == "accept" {
			acceptCount++
		}
	}
	if acceptCount < requestDetailResult.MinIdp {
		return ReturnDeliverTxLog(code.RequestIsNotCompleted, "Request is not completed", "")
	}

	if requestDetailResult.Mode != 3 {
		return ReturnDeliverTxLog(code.InvalidMode, "Onboard request must be mode 3", "")
	}

	if requestDetailResult.MinIdp < 1 {
		return ReturnDeliverTxLog(code.InvalidMinIdp, "Onboard request min_idp must be at least 1", "")
	}
	// check special type of Request && set can used only once
	canAddAccessor := getCanAddAccessor(funcParam.RequestId, app)
	if canAddAccessor != true {
		return ReturnDeliverTxLog(code.RequestIsNotSpecial, "Request is not special", "")
	}
	setCanAddAccessorToFalse(funcParam.RequestId, app)

	var accessor data.Accessor
	accessor.AccessorType = funcParam.AccessorType
	accessor.AccessorPublicKey = funcParam.AccessorPublicKey
	accessor.AccessorGroupId = funcParam.AccessorGroupId
	accessor.Active = true
	accessor.Owner = nodeID

	accessorJSON, err := proto.Marshal(&accessor)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}

	app.SetStateDB([]byte(accessorKey), []byte(accessorJSON))
	return ReturnDeliverTxLog(code.OK, "success", "")
}

func registerMsqDestination(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("RegisterMsqDestination, Parameter: %s", param)
	var funcParam RegisterMsqDestinationParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	nodeDetailKey := "NodeID" + "|" + nodeID
	_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
	if nodeDetailValue == nil {
		return ReturnDeliverTxLog(code.NodeIDNotFound, "Node ID not found", "")
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// Validate user's ial is <= node's max_ial
	for _, user := range funcParam.Users {
		if user.Ial > nodeDetail.MaxIal {
			return ReturnDeliverTxLog(code.IALError, "IAL must be less than or equals to registered node's MAX IAL", "")
		}
	}

	timeOutKey := "TimeOutBlockRegisterMsqDestination"
	var timeOut data.TimeOutBlockRegisterMsqDestination
	_, timeOutValue := app.state.db.Get(prefixKey([]byte(timeOutKey)))
	if timeOutValue != nil {
		err := proto.Unmarshal([]byte(timeOutValue), &timeOut)
		if err != nil {
			return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
	} else {
		timeOut.TimeOutBlock = 500
	}
	timeOutBlockInStateDB := timeOut.TimeOutBlock

	// If validate passed then add Msq Destination
	for _, user := range funcParam.Users {
		key := "MsqDestination" + "|" + user.HashID
		_, chkExists := app.state.db.Get(prefixKey([]byte(key)))

		if chkExists != nil {
			var nodes data.MsqDesList
			err = proto.Unmarshal([]byte(chkExists), &nodes)
			if err != nil {
				return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
			}

			timeoutBlock := app.CurrentBlock + timeOutBlockInStateDB
			var newNode data.Node
			newNode.Ial = user.Ial
			newNode.NodeId = nodeID
			newNode.Active = true
			newNode.First = user.First
			newNode.TimeoutBlock = timeoutBlock
			if !user.First {
				newNode.TimeoutBlock = 0
			}
			// Check duplicate before add
			chkDup := false
			for _, node := range nodes.Nodes {
				if &newNode == node {
					chkDup = true
					break
				}
			}

			// Check first
			if user.First {
				for _, node := range nodes.Nodes {
					if node.TimeoutBlock != 0 {
						if node.TimeoutBlock > app.CurrentBlock {
							return ReturnDeliverTxLog(code.NotFirstIdP, "This node is not first IdP", "")
						}
					}
				}
			}

			if chkDup == false {
				nodes.Nodes = append(nodes.Nodes, &newNode)
				value, err := proto.Marshal(&nodes)
				if err != nil {
					return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
				}
				app.SetStateDB([]byte(key), []byte(value))
			}

		} else {
			var nodes data.MsqDesList
			timeoutBlock := app.CurrentBlock + timeOutBlockInStateDB
			var newNode data.Node
			newNode.Ial = user.Ial
			newNode.NodeId = nodeID
			newNode.Active = true
			newNode.First = user.First
			newNode.TimeoutBlock = timeoutBlock
			if !user.First {
				newNode.TimeoutBlock = 0
			}
			nodes.Nodes = append(nodes.Nodes, &newNode)
			value, err := proto.Marshal(&nodes)
			if err != nil {
				return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
			}
			app.SetStateDB([]byte(key), []byte(value))
		}
	}

	return ReturnDeliverTxLog(code.OK, "success", "")
}

func createIdpResponse(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("CreateIdpResponse, Parameter: %s", param)
	var funcParam CreateIdpResponseParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	key := "Request" + "|" + funcParam.RequestID
	var response data.Response
	response.Ial = funcParam.Ial
	response.Aal = funcParam.Aal
	response.Status = funcParam.Status
	response.Signature = funcParam.Signature
	response.IdpId = nodeID
	response.IdentityProof = funcParam.IdentityProof
	response.PrivateProofHash = funcParam.PrivateProofHash
	_, value := app.state.db.Get(prefixKey([]byte(key)))

	if value == nil {
		return ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	var request data.Request
	err = proto.Unmarshal([]byte(value), &request)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check duplicate before add
	chk := false
	for _, oldResponse := range request.ResponseList {
		if &response == oldResponse {
			chk = true
			break
		}
	}

	// Check AAL
	if request.MinAal > response.Aal {
		return ReturnDeliverTxLog(code.AALError, "Response's AAL is less than min AAL", "")
	}

	// Check IAL
	if request.MinIal > response.Ial {
		return ReturnDeliverTxLog(code.IALError, "Response's IAL is less than min IAL", "")
	}

	// Check AAL, IAL with MaxIalAal
	nodeDetailKey := "NodeID" + "|" + nodeID
	_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
	if nodeDetailValue == nil {
		return ReturnDeliverTxLog(code.NodeIDNotFound, "Node ID not found", "")
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	if response.Aal > nodeDetail.MaxAal {
		return ReturnDeliverTxLog(code.AALError, "Response's AAL is greater than max AAL", "")
	}
	if response.Ial > nodeDetail.MaxIal {
		return ReturnDeliverTxLog(code.IALError, "Response's IAL is greater than max IAL", "")
	}

	// Check min_idp
	if int64(len(request.ResponseList)) >= request.MinIdp {
		return ReturnDeliverTxLog(code.RequestIsCompleted, "Can't response a request that's complete response", "")
	}

	// Check IsClosed
	if request.Closed {
		return ReturnDeliverTxLog(code.RequestIsClosed, "Can't response a request that's closed", "")
	}

	// Check IsTimedOut
	if request.TimedOut {
		return ReturnDeliverTxLog(code.RequestIsTimedOut, "Can't response a request that's timed out", "")
	}

	// Check identity proof if mode == 3
	if request.Mode == 3 {
		identityProofKey := "IdentityProof" + "|" + funcParam.RequestID + "|" + nodeID
		_, identityProofValue := app.state.db.Get(prefixKey([]byte(identityProofKey)))
		proofPassed := false
		if identityProofValue != nil {
			if funcParam.IdentityProof == string(identityProofValue) {
				proofPassed = true
			}
		}
		if proofPassed == false {
			return ReturnDeliverTxLog(code.WrongIdentityProof, "Identity proof is wrong", "")
		}
	}

	if chk == false {
		request.ResponseList = append(request.ResponseList, &response)
		value, err := proto.Marshal(&request)
		if err != nil {
			return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
		}
		app.SetStateDB([]byte(key), []byte(value))
		return ReturnDeliverTxLog(code.OK, "success", funcParam.RequestID)
	}
	return ReturnDeliverTxLog(code.DuplicateResponse, "Duplicate Response", "")
}

func updateIdentity(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("UpdateIdentity, Parameter: %s", param)
	var funcParam UpdateIdentityParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check IAL must less than Max IAL
	nodeDetailKey := "NodeID" + "|" + nodeID
	_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
	if nodeDetailValue == nil {
		return ReturnDeliverTxLog(code.NodeIDNotFound, "Node ID not found", "")
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	if funcParam.Ial > nodeDetail.MaxIal {
		return ReturnDeliverTxLog(code.IALError, "New IAL is greater than max IAL", "")
	}

	msqDesKey := "MsqDestination" + "|" + funcParam.HashID
	_, msqDesValue := app.state.db.Get(prefixKey([]byte(msqDesKey)))
	if msqDesValue != nil {
		var msqDes data.MsqDesList
		err := proto.Unmarshal([]byte(msqDesValue), &msqDes)
		if err != nil {
			return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
		// Selective update
		if funcParam.Ial > 0 {
			for index := range msqDes.Nodes {
				if msqDes.Nodes[index].NodeId == nodeID {
					msqDes.Nodes[index].Ial = funcParam.Ial
					break
				}
			}
		}
		msqDesJSON, err := proto.Marshal(&msqDes)
		if err != nil {
			return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
		}
		app.SetStateDB([]byte(msqDesKey), []byte(msqDesJSON))
		return ReturnDeliverTxLog(code.OK, "success", "")
	}
	return ReturnDeliverTxLog(code.HashIDNotFound, "Hash ID not found", "")
}

func declareIdentityProof(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("DeclareIdentityProof, Parameter: %s", param)
	var funcParam DeclareIdentityProofParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check the request
	requestKey := "Request" + "|" + funcParam.RequestID
	_, requestValue := app.state.db.Get(prefixKey([]byte(requestKey)))

	if requestValue == nil {
		return ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	var request data.Request
	err = proto.Unmarshal([]byte(requestValue), &request)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// check number of responses
	if int64(len(request.ResponseList)) >= request.MinIdp {
		return ReturnDeliverTxLog(code.RequestIsCompleted, "Can't declare identity proof for the request that's completed response", "")
	}

	// Check IsClosed
	if request.Closed {
		return ReturnDeliverTxLog(code.RequestIsClosed, "Can't declare identity proof for the request that's closed", "")
	}

	// Check IsTimedOut
	if request.TimedOut {
		return ReturnDeliverTxLog(code.RequestIsTimedOut, "Can't declare identity proof for the request that's timed out", "")
	}

	identityProofKey := "IdentityProof" + "|" + funcParam.RequestID + "|" + nodeID
	_, identityProofValue := app.state.db.Get(prefixKey([]byte(identityProofKey)))
	if identityProofValue == nil {
		identityProofValue := funcParam.IdentityProof
		app.SetStateDB([]byte(identityProofKey), []byte(identityProofValue))
		return ReturnDeliverTxLog(code.OK, "success", "")
	}
	return ReturnDeliverTxLog(code.DuplicateIdentityProof, "Duplicate Identity Proof", "")
}

func clearRegisterMsqDestinationTimeout(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("ClearRegisterMsqDestinationTimeout, Parameter: %s", param)
	var funcParam ClearRegisterMsqDestinationTimeoutParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	msqDesKey := "MsqDestination" + "|" + funcParam.HashID
	_, msqDesValue := app.state.db.Get(prefixKey([]byte(msqDesKey)))

	if msqDesValue != nil {
		var nodes data.MsqDesList
		err = proto.Unmarshal([]byte(msqDesValue), &nodes)
		if err != nil {
			return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}

		// Check is not timeout
		for index := range nodes.Nodes {
			if nodes.Nodes[index].NodeId == nodeID {
				if nodes.Nodes[index].TimeoutBlock <= app.CurrentBlock {
					return ReturnDeliverTxLog(code.MsqDestinationIsTimedOut, "Can not clear msq destination that is timed out", "")
				}
				break
			}
		}

		for index := range nodes.Nodes {
			if nodes.Nodes[index].NodeId == nodeID {
				nodes.Nodes[index].TimeoutBlock = 0
				break
			}
		}

		msqDesJSON, err := proto.Marshal(&nodes)
		if err != nil {
			return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
		}
		app.SetStateDB([]byte(msqDesKey), []byte(msqDesJSON))
		return ReturnDeliverTxLog(code.OK, "success", "")
	}
	return ReturnDeliverTxLog(code.HashIDNotFound, "Hash ID not found", "")
}
