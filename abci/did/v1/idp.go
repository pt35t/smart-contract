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

	"github.com/golang/protobuf/proto"
	"github.com/ndidplatform/smart-contract/abci/utils"
	"github.com/ndidplatform/smart-contract/abci/code"
	"github.com/ndidplatform/smart-contract/protos/data"
	"github.com/tendermint/tendermint/abci/types"
)

func (app *DIDApplication) registerAccessor(param string, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("RegisterAccessor, Parameter: %s", param)
	var funcParam RegisterAccessorParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	accessorKey := "Accessor" + "|" + funcParam.AccessorID
	var accessor data.Accessor
	accessor.AccessorType = funcParam.AccessorType
	accessor.AccessorPublicKey = funcParam.AccessorPublicKey
	accessor.AccessorGroupId = funcParam.AccessorGroupID
	accessor.Active = true
	accessor.Owner = nodeID
	accessorJSON, err := utils.ProtoDeterministicMarshal(&accessor)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	accessorGroupKey := "AccessorGroup" + "|" + funcParam.AccessorGroupID
	accessorGroup := funcParam.AccessorGroupID
	// Check duplicate accessor_id
	_, chkAccessorKeyExists := app.state.db.Get(prefixKey([]byte(accessorKey)))
	if chkAccessorKeyExists != nil {
		return app.ReturnDeliverTxLog(code.DuplicateAccessorID, "Duplicate Accessor ID", "")
	}
	// Check duplicate accessor_group_id
	_, chkAccessorGroupKeyExists := app.state.db.Get(prefixKey([]byte(accessorGroupKey)))
	if chkAccessorGroupKeyExists != nil {
		return app.ReturnDeliverTxLog(code.DuplicateAccessorGroupID, "Duplicate Accessor Group ID", "")
	}
	// Add relation AccessorGroupID -> AccessorID
	accessorInGroupKey := "AccessorInGroup" + "|" + funcParam.AccessorGroupID
	var accessorInGroup data.AccessorInGroup
	accessorInGroup.Accessors = append(accessorInGroup.Accessors, funcParam.AccessorID)
	accessorInGroupProtobuf, err := utils.ProtoDeterministicMarshal(&accessorInGroup)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(accessorKey), []byte(accessorJSON))
	app.SetStateDB([]byte(accessorGroupKey), []byte(accessorGroup))
	app.SetStateDB([]byte(accessorInGroupKey), []byte(accessorInGroupProtobuf))
	return app.ReturnDeliverTxLog(code.OK, "success", "")
}

func (app *DIDApplication) addAccessorMethod(param string, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("AddAccessorMethod, Parameter: %s", param)
	var funcParam AccessorMethod
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// AccessorGroupID: must already exist
	accessorGroupKey := "AccessorGroup" + "|" + funcParam.AccessorGroupID
	_, chkAccessorGroupKeyExists := app.state.db.Get(prefixKey([]byte(accessorGroupKey)))
	if chkAccessorGroupKeyExists == nil {
		return app.ReturnDeliverTxLog(code.AccessorGroupIDNotFound, "Accessor Group ID not found", "")
	}
	// AccessorID: must not duplicate
	accessorKey := "Accessor" + "|" + funcParam.AccessorID
	_, chkAccessorKeyExists := app.state.db.Get(prefixKey([]byte(accessorKey)))
	if chkAccessorKeyExists != nil {
		return app.ReturnDeliverTxLog(code.DuplicateAccessorID, "Duplicate Accessor ID", "")
	}
	// Request must be completed, can be used only once, special type
	var getRequestparam GetRequestParam
	getRequestparam.RequestID = funcParam.RequestID
	getRequestparamJSON, err := json.Marshal(getRequestparam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	var requestDetail = app.getRequestDetail(string(getRequestparamJSON), app.state.db.Version64())
	var requestDetailResult GetRequestDetailResult
	err = json.Unmarshal([]byte(requestDetail.Value), &requestDetailResult)
	if err != nil {
		return app.ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	// Check accept result >= min_idp
	acceptCount := 0
	for _, response := range requestDetailResult.Responses {
		// Check status is 'accept', valid IAL = true, valid Proof = true and valid signature = true
		if response.ValidIal == nil {
			continue
		}
		if response.ValidProof == nil {
			continue
		}
		if response.ValidSignature == nil {
			continue
		}
		if response.Status == "accept" && *response.ValidIal && *response.ValidProof && *response.ValidSignature {
			acceptCount++
		}
	}
	if acceptCount < requestDetailResult.MinIdp {
		return app.ReturnDeliverTxLog(code.RequestIsNotCompleted, "Request is not completed", "")
	}
	// check request is closed
	if !requestDetailResult.IsClosed {
		return app.ReturnDeliverTxLog(code.RequestIsNotClosed, "Request is not closed", "")
	}
	if requestDetailResult.Mode != 3 {
		return app.ReturnDeliverTxLog(code.InvalidMode, "Onboard request must be mode 3", "")
	}
	if requestDetailResult.MinIdp < 1 {
		return app.ReturnDeliverTxLog(code.InvalidMinIdp, "Onboard request min_idp must be at least 1", "")
	}
	requestKey := "Request" + "|" + funcParam.RequestID
	_, requestValue := app.state.db.Get(prefixKey([]byte(requestKey)))
	if requestValue == nil {
		return app.ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	var request data.Request
	err = proto.Unmarshal([]byte(requestValue), &request)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// check special type of Request && set can used only once
	if request.Purpose != "AddAccessor" || request.UseCount != 0 {
		return app.ReturnDeliverTxLog(code.RequestIsNotSpecial, "Request is not special", "")
	}
	// set use count
	request.UseCount = 1
	var accessor data.Accessor
	accessor.AccessorType = funcParam.AccessorType
	accessor.AccessorPublicKey = funcParam.AccessorPublicKey
	accessor.AccessorGroupId = funcParam.AccessorGroupID
	accessor.Active = true
	accessor.Owner = nodeID
	accessorJSON, err := utils.ProtoDeterministicMarshal(&accessor)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	// Add relation AccessorGroupID -> AccessorID
	accessorInGroupKey := "AccessorInGroup" + "|" + funcParam.AccessorGroupID
	_, accessorInGroupKeyValue := app.state.db.Get(prefixKey([]byte(accessorInGroupKey)))
	var accessors data.AccessorInGroup
	err = proto.Unmarshal(accessorInGroupKeyValue, &accessors)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	accessors.Accessors = append(accessors.Accessors, funcParam.AccessorID)
	accessorInGroupProtobuf, err := utils.ProtoDeterministicMarshal(&accessors)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	requestProtobuf, err := utils.ProtoDeterministicMarshal(&request)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(requestKey), []byte(requestProtobuf))
	app.SetStateDB([]byte(accessorInGroupKey), []byte(accessorInGroupProtobuf))
	app.SetStateDB([]byte(accessorKey), []byte(accessorJSON))
	return app.ReturnDeliverTxLog(code.OK, "success", "")
}

func (app *DIDApplication) registerIdentity(param string, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("RegisterIdentity, Parameter: %s", param)
	var funcParam RegisterIdentityParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	nodeDetailKey := "NodeID" + "|" + nodeID
	_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
	if nodeDetailValue == nil {
		return app.ReturnDeliverTxLog(code.NodeIDNotFound, "Node ID not found", "")
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// Validate user's ial is <= node's max_ial
	for _, user := range funcParam.Users {
		if user.Ial > nodeDetail.MaxIal {
			return app.ReturnDeliverTxLog(code.IALError, "IAL must be less than or equals to registered node's MAX IAL", "")
		}
	}
	timeOutKey := "TimeOutBlockRegisterIdentity"
	var timeOut data.TimeOutBlockRegisterIdentity
	_, timeOutValue := app.state.db.Get(prefixKey([]byte(timeOutKey)))
	if timeOutValue != nil {
		err := proto.Unmarshal([]byte(timeOutValue), &timeOut)
		if err != nil {
			return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
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
				return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
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
							return app.ReturnDeliverTxLog(code.NotFirstIdP, "This node is not first IdP", "")
						}
					}
				}
			}
			if chkDup == false {
				nodes.Nodes = append(nodes.Nodes, &newNode)
				value, err := utils.ProtoDeterministicMarshal(&nodes)
				if err != nil {
					return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
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
			value, err := utils.ProtoDeterministicMarshal(&nodes)
			if err != nil {
				return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
			}
			app.SetStateDB([]byte(key), []byte(value))
		}
	}
	return app.ReturnDeliverTxLog(code.OK, "success", "")
}

func (app *DIDApplication) createIdpResponse(param string, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("CreateIdpResponse, Parameter: %s", param)
	var funcParam CreateIdpResponseParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
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
		return app.ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	var request data.Request
	err = proto.Unmarshal([]byte(value), &request)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// Check duplicate before add
	chkDup := false
	for _, oldResponse := range request.ResponseList {
		if &response == oldResponse {
			chkDup = true
			break
		}
	}
	// Check AAL
	if request.MinAal > response.Aal {
		return app.ReturnDeliverTxLog(code.AALError, "Response's AAL is less than min AAL", "")
	}
	// Check IAL
	if request.MinIal > response.Ial {
		return app.ReturnDeliverTxLog(code.IALError, "Response's IAL is less than min IAL", "")
	}
	// Check AAL, IAL with MaxIalAal
	nodeDetailKey := "NodeID" + "|" + nodeID
	_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
	if nodeDetailValue == nil {
		return app.ReturnDeliverTxLog(code.NodeIDNotFound, "Node ID not found", "")
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	if response.Aal > nodeDetail.MaxAal {
		return app.ReturnDeliverTxLog(code.AALError, "Response's AAL is greater than max AAL", "")
	}
	if response.Ial > nodeDetail.MaxIal {
		return app.ReturnDeliverTxLog(code.IALError, "Response's IAL is greater than max IAL", "")
	}
	// Check min_idp
	if int64(len(request.ResponseList)) >= request.MinIdp {
		return app.ReturnDeliverTxLog(code.RequestIsCompleted, "Can't response a request that's complete response", "")
	}
	// Check IsClosed
	if request.Closed {
		return app.ReturnDeliverTxLog(code.RequestIsClosed, "Can't response a request that's closed", "")
	}
	// Check IsTimedOut
	if request.TimedOut {
		return app.ReturnDeliverTxLog(code.RequestIsTimedOut, "Can't response a request that's timed out", "")
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
			return app.ReturnDeliverTxLog(code.WrongIdentityProof, "Identity proof is wrong", "")
		}
	}
	// Check nodeID is exist in idp_id_list
	exist := false
	for _, idpID := range request.IdpIdList {
		if idpID == nodeID {
			exist = true
			break
		}
	}
	if exist == false {
		return app.ReturnDeliverTxLog(code.NodeIDDoesNotExistInIdPList, "Node ID does not exist in IdP list", "")
	}
	if chkDup == true {
		return app.ReturnDeliverTxLog(code.DuplicateResponse, "Duplicate Response", "")
	}
	request.ResponseList = append(request.ResponseList, &response)
	value, err = utils.ProtoDeterministicMarshal(&request)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(key), []byte(value))
	return app.ReturnDeliverTxLog(code.OK, "success", funcParam.RequestID)
}

func (app *DIDApplication) updateIdentity(param string, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("UpdateIdentity, Parameter: %s", param)
	var funcParam UpdateIdentityParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// Check IAL must less than Max IAL
	nodeDetailKey := "NodeID" + "|" + nodeID
	_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
	if nodeDetailValue == nil {
		return app.ReturnDeliverTxLog(code.NodeIDNotFound, "Node ID not found", "")
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	if funcParam.Ial > nodeDetail.MaxIal {
		return app.ReturnDeliverTxLog(code.IALError, "New IAL is greater than max IAL", "")
	}
	msqDesKey := "MsqDestination" + "|" + funcParam.HashID
	_, msqDesValue := app.state.db.Get(prefixKey([]byte(msqDesKey)))
	if msqDesValue == nil {
		return app.ReturnDeliverTxLog(code.HashIDNotFound, "Hash ID not found", "")
	}
	var msqDes data.MsqDesList
	err = proto.Unmarshal([]byte(msqDesValue), &msqDes)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
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
	msqDesJSON, err := utils.ProtoDeterministicMarshal(&msqDes)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(msqDesKey), []byte(msqDesJSON))
	return app.ReturnDeliverTxLog(code.OK, "success", "")
}

func (app *DIDApplication) declareIdentityProof(param string, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("DeclareIdentityProof, Parameter: %s", param)
	var funcParam DeclareIdentityProofParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// Check the request
	requestKey := "Request" + "|" + funcParam.RequestID
	_, requestValue := app.state.db.Get(prefixKey([]byte(requestKey)))
	if requestValue == nil {
		return app.ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	var request data.Request
	err = proto.Unmarshal([]byte(requestValue), &request)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// check number of responses
	if int64(len(request.ResponseList)) >= request.MinIdp {
		return app.ReturnDeliverTxLog(code.RequestIsCompleted, "Can't declare identity proof for the request that's completed response", "")
	}
	// Check IsClosed
	if request.Closed {
		return app.ReturnDeliverTxLog(code.RequestIsClosed, "Can't declare identity proof for the request that's closed", "")
	}
	// Check IsTimedOut
	if request.TimedOut {
		return app.ReturnDeliverTxLog(code.RequestIsTimedOut, "Can't declare identity proof for the request that's timed out", "")
	}
	identityProofKey := "IdentityProof" + "|" + funcParam.RequestID + "|" + nodeID
	_, identityProofValue := app.state.db.Get(prefixKey([]byte(identityProofKey)))
	if identityProofValue != nil {
		return app.ReturnDeliverTxLog(code.DuplicateIdentityProof, "Duplicate Identity Proof", "")
	}
	identityProofValue = []byte(funcParam.IdentityProof)
	app.SetStateDB([]byte(identityProofKey), identityProofValue)
	return app.ReturnDeliverTxLog(code.OK, "success", "")
}

func (app *DIDApplication) clearRegisterIdentityTimeout(param string, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("ClearRegisterIdentityTimeout, Parameter: %s", param)
	var funcParam ClearRegisterIdentityTimeoutParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	msqDesKey := "MsqDestination" + "|" + funcParam.HashID
	_, msqDesValue := app.state.db.Get(prefixKey([]byte(msqDesKey)))
	if msqDesValue == nil {
		return app.ReturnDeliverTxLog(code.HashIDNotFound, "Hash ID not found", "")
	}
	var nodes data.MsqDesList
	err = proto.Unmarshal([]byte(msqDesValue), &nodes)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// Check is not timeout
	for index := range nodes.Nodes {
		if nodes.Nodes[index].NodeId == nodeID {
			if nodes.Nodes[index].TimeoutBlock <= app.CurrentBlock {
				return app.ReturnDeliverTxLog(code.RegisterIdentityIsTimedOut, "Cannot clear register identity that is timed out", "")
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
	msqDesJSON, err := utils.ProtoDeterministicMarshal(&nodes)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(msqDesKey), []byte(msqDesJSON))
	return app.ReturnDeliverTxLog(code.OK, "success", "")
}

func (app *DIDApplication) revokeAccessorMethod(param string, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("RevokeAccessorMethod, Parameter: %s", param)
	var funcParam RevokeAccessorMethodParam
	err := json.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// Request must be completed, can be used only once, special type
	var getRequestparam GetRequestParam
	getRequestparam.RequestID = funcParam.RequestID
	getRequestparamJSON, err := json.Marshal(getRequestparam)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	var requestDetail = app.getRequestDetail(string(getRequestparamJSON), app.state.db.Version64())
	var requestDetailResult GetRequestDetailResult
	err = json.Unmarshal([]byte(requestDetail.Value), &requestDetailResult)
	if err != nil {
		return app.ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	if requestDetailResult.RequestID == "" {
		return app.ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	// Check accept result >= min_idp
	acceptCount := 0
	for _, response := range requestDetailResult.Responses {
		// Check status is 'accept', valid IAL = true, valid Proof = true and valid signature = true
		if response.ValidIal == nil {
			continue
		}
		if response.ValidProof == nil {
			continue
		}
		if response.ValidSignature == nil {
			continue
		}
		// All of the response(s) in the request must belong to IDP that call the function
		if response.IdpID != nodeID {
			continue
		}
		if response.Status == "accept" && *response.ValidIal && *response.ValidProof && *response.ValidSignature {
			acceptCount++
		}
	}
	if acceptCount < requestDetailResult.MinIdp {
		return app.ReturnDeliverTxLog(code.RequestIsNotCompleted, "Request is not completed", "")
	}
	// check request is closed
	if !requestDetailResult.IsClosed {
		return app.ReturnDeliverTxLog(code.RequestIsNotClosed, "Request is not closed", "")
	}
	if requestDetailResult.Mode != 3 {
		return app.ReturnDeliverTxLog(code.InvalidMode, "Onboard request must be mode 3", "")
	}
	if requestDetailResult.MinIdp < 1 {
		return app.ReturnDeliverTxLog(code.InvalidMinIdp, "Onboard request min_idp must be at least 1", "")
	}
	requestKey := "Request" + "|" + funcParam.RequestID
	_, requestValue := app.state.db.Get(prefixKey([]byte(requestKey)))
	if requestValue == nil {
		return app.ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	var request data.Request
	err = proto.Unmarshal([]byte(requestValue), &request)
	if err != nil {
		return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	// check special type of Request && set can used only once
	// Request must be complete and closed, with purpose "RevokeAccessor"
	if request.Purpose != "RevokeAccessor" || request.UseCount != 0 {
		return app.ReturnDeliverTxLog(code.RequestIsNotSpecial, "Request is not special", "")
	}
	// set use count
	request.UseCount = 1
	// All accessor key with id in "accessor_id_list" must be created by IDP that call the function
	for _, accessorID := range funcParam.AccessorIDList {
		accessorKey := "Accessor" + "|" + accessorID
		_, accessorValue := app.state.db.Get(prefixKey([]byte(accessorKey)))
		if accessorValue == nil {
			return app.ReturnDeliverTxLog(code.AccessorIDNotFound, "Accessor ID not found", "")
		}
		var accessor data.Accessor
		err = proto.Unmarshal([]byte(accessorValue), &accessor)
		if err != nil {
			return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
		// Check node ID is owner of accessor
		if nodeID != accessor.Owner {
			return app.ReturnDeliverTxLog(code.NotOwnerOfAccessor, "Node ID is not owner of accessor", "")
		}
	}
	// If all condition pass, disable all accessor key with id in "accessor_id_list"
	for _, accessorID := range funcParam.AccessorIDList {
		accessorKey := "Accessor" + "|" + accessorID
		_, accessorValue := app.state.db.Get(prefixKey([]byte(accessorKey)))
		if accessorValue == nil {
			return app.ReturnDeliverTxLog(code.AccessorIDNotFound, "Accessor ID not found", "")
		}
		var accessor data.Accessor
		err = proto.Unmarshal([]byte(accessorValue), &accessor)
		if err != nil {
			return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
		// Set disable
		accessor.Active = false
		accessorProtobuf, err := utils.ProtoDeterministicMarshal(&accessor)
		if err != nil {
			return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
		}
		// Remove AccessorID from AccessorInGroup
		accessorInGroupKey := "AccessorInGroup" + "|" + accessor.AccessorGroupId
		_, accessorInGroupKeyValue := app.state.db.Get(prefixKey([]byte(accessorInGroupKey)))
		var accessors data.AccessorInGroup
		err = proto.Unmarshal(accessorInGroupKeyValue, &accessors)
		if err != nil {
			return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
		for i, accessorIDInList := range accessors.Accessors {
			if accessorIDInList == accessorID {
				copy(accessors.Accessors[i:], accessors.Accessors[i+1:])
				accessors.Accessors[len(accessors.Accessors)-1] = ""
				accessors.Accessors = accessors.Accessors[:len(accessors.Accessors)-1]
			}
		}
		accessorInGroupProtobuf, err := utils.ProtoDeterministicMarshal(&accessors)
		if err != nil {
			return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
		}
		// Add accessor ID to revokedAccessorInGroup
		revokedAccessorInGroupKey := "RevokedAccessorInGroup" + "|" + accessor.AccessorGroupId
		_, revokedAccessorInGroupValue := app.state.db.Get(prefixKey([]byte(revokedAccessorInGroupKey)))
		var revokedAccessorInGroup data.AccessorInGroup
		err = proto.Unmarshal(revokedAccessorInGroupValue, &revokedAccessorInGroup)
		if err != nil {
			return app.ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
		revokedAccessorInGroup.Accessors = append(revokedAccessorInGroup.Accessors, accessorID)
		revokedAccessorInGroupProtobuf, err := utils.ProtoDeterministicMarshal(&revokedAccessorInGroup)
		if err != nil {
			return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
		}
		app.SetStateDB([]byte(accessorKey), []byte(accessorProtobuf))
		app.SetStateDB([]byte(accessorInGroupKey), []byte(accessorInGroupProtobuf))
		app.SetStateDB([]byte(revokedAccessorInGroupKey), []byte(revokedAccessorInGroupProtobuf))
	}
	requestProtobuf, err := utils.ProtoDeterministicMarshal(&request)
	if err != nil {
		return app.ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(requestKey), []byte(requestProtobuf))
	return app.ReturnDeliverTxLog(code.OK, "success", "")
}
