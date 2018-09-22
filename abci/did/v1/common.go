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
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/ndidplatform/smart-contract/abci/code"
	"github.com/ndidplatform/smart-contract/protos/data"
	pbParam "github.com/ndidplatform/smart-contract/protos/params"
	pbResult "github.com/ndidplatform/smart-contract/protos/result"
	"github.com/tendermint/tendermint/abci/types"
)

func registerMsqAddress(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("RegisterMsqAddress, Parameter: %s", param)
	var funcParam pbParam.RegisterMsqAddressParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	nodeDetailKey := "NodeID" + "|" + funcParam.NodeId
	_, value := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal(value, &nodeDetail)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	var msqAddress []*data.MQ
	for _, address := range funcParam.Addresses {
		var msq data.MQ
		msq.Ip = address.Ip
		msq.Port = address.Port
		msqAddress = append(msqAddress, &msq)
	}
	nodeDetail.Mq = msqAddress
	nodeDetailByte, err := proto.Marshal(&nodeDetail)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(nodeDetailKey), []byte(nodeDetailByte))
	return ReturnDeliverTxLog(code.OK, "success", "")
}

func getNodeMasterPublicKey(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetNodeMasterPublicKey, Parameter: %s", param)
	var funcParam pbParam.GetNodeMasterPublicKeyParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	key := "NodeID" + "|" + funcParam.NodeId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)
	var res pbResult.GetNodeMasterPublicKeyResult
	if value == nil {
		valueJSON, err := proto.Marshal(&res)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(valueJSON, "not found", app.state.db.Version64(), app)
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal(value, &nodeDetail)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	res.MasterPublicKey = nodeDetail.MasterPublicKey
	valueJSON, err := proto.Marshal(&res)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(valueJSON, "success", app.state.db.Version64(), app)

}

func getNodePublicKey(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetNodePublicKey, Parameter: %s", param)
	var funcParam pbParam.GetNodePublicKeyParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	key := "NodeID" + "|" + funcParam.NodeId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)
	var res pbResult.GetNodePublicKeyResult
	if value == nil {
		valueJSON, err := proto.Marshal(&res)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(valueJSON, "not found", app.state.db.Version64(), app)
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal(value, &nodeDetail)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	res.PublicKey = nodeDetail.PublicKey
	valueJSON, err := proto.Marshal(&res)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(valueJSON, "success", app.state.db.Version64(), app)
}

func getNodeNameByNodeID(nodeID string, app *DIDApplication) string {
	key := "NodeID" + "|" + nodeID
	_, value := app.state.db.Get(prefixKey([]byte(key)))
	if value != nil {
		var nodeDetail data.NodeDetail
		err := proto.Unmarshal([]byte(value), &nodeDetail)
		if err != nil {
			return ""
		}
		return nodeDetail.NodeName
	}
	return ""
}

func getIdpNodes(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetIdpNodes, Parameter: %s", param)
	var funcParam pbParam.GetIdpNodesParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var returnNodes pbResult.GetIdpNodesResult
	returnNodes.Node = make([]*pbResult.MsqDestinationNode, 0)

	if funcParam.HashId == "" {
		idpsKey := "IdPList"
		_, idpsValue := app.state.db.GetVersioned(prefixKey([]byte(idpsKey)), height)
		var idpsList data.IdPList
		if idpsValue != nil {
			err := proto.Unmarshal(idpsValue, &idpsList)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			for _, idp := range idpsList.NodeId {
				nodeDetailKey := "NodeID" + "|" + idp
				_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
				if nodeDetailValue == nil {
					continue
				}
				var nodeDetail data.NodeDetail
				err := proto.Unmarshal(nodeDetailValue, &nodeDetail)
				if err != nil {
					continue
				}
				// check node is active
				if !nodeDetail.Active {
					continue
				}
				// check Max IAL && AAL
				if !(nodeDetail.MaxIal >= funcParam.MinIal &&
					nodeDetail.MaxAal >= funcParam.MinAal) {
					continue
				}
				var msqDesNode pbResult.MsqDestinationNode
				msqDesNode.NodeId = idp
				msqDesNode.NodeName = nodeDetail.NodeName
				msqDesNode.MaxIal = nodeDetail.MaxIal
				msqDesNode.MaxAal = nodeDetail.MaxAal
				returnNodes.Node = append(returnNodes.Node, &msqDesNode)
			}
		}
	} else {
		key := "MsqDestination" + "|" + funcParam.HashId
		_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)

		if value != nil {
			var nodes data.MsqDesList
			err = proto.Unmarshal([]byte(value), &nodes)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}

			for _, node := range nodes.Nodes {
				// check msq destination is not active
				if !node.Active {
					continue
				}
				// check Ial > min ial
				if node.Ial < funcParam.MinIal {
					continue
				}
				// check msq destination is not timed out
				if node.TimeoutBlock != 0 && app.CurrentBlock > node.TimeoutBlock {
					continue
				}
				nodeDetailKey := "NodeID" + "|" + node.NodeId
				_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
				if nodeDetailValue == nil {
					continue
				}
				var nodeDetail data.NodeDetail
				err := proto.Unmarshal(nodeDetailValue, &nodeDetail)
				if err != nil {
					continue
				}
				// check node is active
				if !nodeDetail.Active {
					continue
				}
				// check Max IAL && AAL
				if !(nodeDetail.MaxIal >= funcParam.MinIal &&
					nodeDetail.MaxAal >= funcParam.MinAal) {
					continue
				}
				var msqDesNode pbResult.MsqDestinationNode
				msqDesNode.NodeId = node.NodeId
				msqDesNode.NodeName = nodeDetail.NodeName
				msqDesNode.MaxIal = nodeDetail.MaxIal
				msqDesNode.MaxAal = nodeDetail.MaxAal
				returnNodes.Node = append(returnNodes.Node, &msqDesNode)
			}
		}
	}

	value, err := proto.Marshal(&returnNodes)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	// return ReturnQuery(value, "success", app.state.db.Version64(), app)
	if len(returnNodes.Node) > 0 {
		return ReturnQuery(value, "success", app.state.db.Version64(), app)
	}
	return ReturnQuery(value, "not found", app.state.db.Version64(), app)
}

func getAsNodesByServiceId(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetAsNodesByServiceId, Parameter: %s", param)
	var funcParam pbParam.GetAsNodesByServiceIdParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	key := "ServiceDestination" + "|" + funcParam.ServiceId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)

	if value == nil {
		var result pbResult.GetAsNodesByServiceIdResult
		result.Node = make([]*pbResult.ASNodeInGetAsNodesByServiceIdResult, 0)
		value, err := proto.Marshal(&result)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(value, "not found", app.state.db.Version64(), app)
	}

	// filter serive is active
	serviceKey := "Service" + "|" + funcParam.ServiceId
	_, serviceValue := app.state.db.Get(prefixKey([]byte(serviceKey)))
	if serviceValue != nil {
		var service data.ServiceDetail
		err = proto.Unmarshal([]byte(serviceValue), &service)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		if service.Active == false {
			var result pbResult.GetAsNodesByServiceIdResult
			result.Node = make([]*pbResult.ASNodeInGetAsNodesByServiceIdResult, 0)
			value, err := proto.Marshal(&result)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			return ReturnQuery(value, "service is not active", app.state.db.Version64(), app)
		}
	} else {
		var result pbResult.GetAsNodesByServiceIdResult
		result.Node = make([]*pbResult.ASNodeInGetAsNodesByServiceIdResult, 0)
		value, err := proto.Marshal(&result)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(value, "not found", app.state.db.Version64(), app)
	}

	var storedData data.ServiceDesList
	err = proto.Unmarshal([]byte(value), &storedData)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var result pbResult.GetAsNodesByServiceIdResult
	result.Node = make([]*pbResult.ASNodeInGetAsNodesByServiceIdResult, 0)
	for index := range storedData.Node {

		// filter service destination is Active
		if !storedData.Node[index].Active {
			continue
		}

		// Filter approve from NDID
		approveServiceKey := "ApproveKey" + "|" + funcParam.ServiceId + "|" + storedData.Node[index].NodeId
		_, approveServiceJSON := app.state.db.Get(prefixKey([]byte(approveServiceKey)))
		if approveServiceJSON == nil {
			continue
		}
		var approveService data.ApproveService
		err = proto.Unmarshal([]byte(approveServiceJSON), &approveService)
		if err != nil {
			continue
		}
		if !approveService.Active {
			continue
		}

		nodeDetailKey := "NodeID" + "|" + storedData.Node[index].NodeId
		_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
		if nodeDetailValue == nil {
			continue
		}
		var nodeDetail data.NodeDetail
		err := proto.Unmarshal(nodeDetailValue, &nodeDetail)
		if err != nil {
			continue
		}

		// filter node is active
		if !nodeDetail.Active {
			continue
		}

		var newRow pbResult.ASNodeInGetAsNodesByServiceIdResult
		newRow.NodeId = storedData.Node[index].NodeId
		newRow.NodeName = nodeDetail.NodeName
		newRow.MinIal = storedData.Node[index].MinIal
		newRow.MinAal = storedData.Node[index].MinAal
		result.Node = append(result.Node, &newRow)
	}
	resultJSON, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(resultJSON, "success", app.state.db.Version64(), app)
}

func getMsqAddress(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetMsqAddress, Parameter: %s", param)
	var funcParam pbParam.GetMsqAddressParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	nodeDetailKey := "NodeID" + "|" + funcParam.NodeId
	_, value := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal(value, &nodeDetail)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	if value == nil {
		return ReturnQuery(nil, "not found", app.state.db.Version64(), app)
	}
	var result pbResult.GetMsqAddressResult
	result.Mq = make([]*pbResult.MsqAddressInResult, 0)
	for _, msq := range nodeDetail.Mq {
		var newRow pbResult.MsqAddressInResult
		newRow.Ip = msq.Ip
		newRow.Port = msq.Port
		result.Mq = append(result.Mq, &newRow)
	}
	resultJSON, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	if len(result.Mq) == 0 {
		return ReturnQuery(resultJSON, "not found", app.state.db.Version64(), app)
	}
	return ReturnQuery(resultJSON, "success", app.state.db.Version64(), app)
}

func getCanAddAccessor(requestID string, app *DIDApplication) bool {
	result := false
	key := "Request" + "|" + requestID
	_, value := app.state.db.Get(prefixKey([]byte(key)))
	if value != nil {
		var request data.Request
		err := proto.Unmarshal([]byte(value), &request)
		if err == nil {
			if request.CanAddAccessor {
				result = true
			}
		}
	}
	return result
}

func getRequest(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetRequest, Parameter: %s", param)
	var funcParam pbParam.GetRequestParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	key := "Request" + "|" + funcParam.RequestId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)

	if value == nil {
		return ReturnQuery(nil, "not found", app.state.db.Version64(), app)
	}
	var request data.Request
	err = proto.Unmarshal([]byte(value), &request)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var res pbResult.GetRequestResult
	res.Closed = request.Closed
	res.TimedOut = request.TimedOut
	res.RequestMessageHash = request.RequestMessageHash
	res.Mode = request.Mode

	valueJSON, err := proto.Marshal(&res)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(valueJSON, "success", app.state.db.Version64(), app)
}

func getRequestDetail(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetRequestDetail, Parameter: %s", param)
	var funcParam pbParam.GetRequestParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	key := "Request" + "|" + funcParam.RequestId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)

	if value == nil {
		valueJSON := []byte("{}")
		return ReturnQuery(valueJSON, "not found", app.state.db.Version64(), app)
	}

	var result pbResult.GetRequestDetailResult
	var request data.Request
	err = proto.Unmarshal([]byte(value), &request)
	if err != nil {
		value = []byte("")
		return ReturnQuery(value, err.Error(), app.state.db.Version64(), app)
	}

	result.RequestId = request.RequestId
	result.MinIdp = request.MinIdp
	result.MinAal = float64(request.MinAal)
	result.MinIal = float64(request.MinIal)
	result.RequestTimeout = request.RequestTimeout
	for _, dataRequest := range request.DataRequestList {
		var newRow pbResult.DataRequestInResult
		newRow.ServiceId = dataRequest.ServiceId
		newRow.AsIdList = dataRequest.AsIdList
		newRow.MinAs = dataRequest.MinAs
		newRow.AnsweredAsIdList = dataRequest.AnsweredAsIdList
		newRow.ReceivedDataFromList = dataRequest.ReceivedDataFromList
		newRow.RequestParamsHash = dataRequest.RequestParamsHash
		if newRow.AsIdList == nil {
			newRow.AsIdList = make([]string, 0)
		}
		if newRow.AnsweredAsIdList == nil {
			newRow.AnsweredAsIdList = make([]string, 0)
		}
		if newRow.ReceivedDataFromList == nil {
			newRow.ReceivedDataFromList = make([]string, 0)
		}
		result.DataRequestList = append(result.DataRequestList, &newRow)
	}
	result.RequestMessageHash = request.RequestMessageHash
	for _, response := range request.ResponseList {
		var newRow pbResult.ResponseInResult
		newRow.Ial = response.Ial
		newRow.Aal = response.Aal
		newRow.Status = response.Status
		newRow.Signature = response.Signature
		newRow.IdentityProof = response.IdentityProof
		newRow.PrivateProofHash = response.PrivateProofHash
		newRow.IdpId = response.IdpId
		if response.ValidProof != "" {
			if response.ValidProof == "true" {
				var boolProof pbResult.ResponseInResult_ValidProofBool
				boolProof.ValidProofBool = true
				newRow.ValidProof = &boolProof
			} else {
				var boolProof pbResult.ResponseInResult_ValidProofBool
				boolProof.ValidProofBool = false
				newRow.ValidProof = &boolProof
			}

		}
		if response.ValidIal != "" {
			if response.ValidIal == "true" {
				var boolIal pbResult.ResponseInResult_ValidIalBool
				boolIal.ValidIalBool = true
				newRow.ValidIal = &boolIal
			} else {
				var boolIal pbResult.ResponseInResult_ValidIalBool
				boolIal.ValidIalBool = false
				newRow.ValidIal = &boolIal
			}

		}
		if response.ValidSignature != "" {
			if response.ValidSignature == "true" {
				var boolSignature pbResult.ResponseInResult_ValidSignatureBool
				boolSignature.ValidSignatureBool = true
				newRow.ValidSignature = &boolSignature
			} else {
				var boolSignature pbResult.ResponseInResult_ValidSignatureBool
				boolSignature.ValidSignatureBool = false
				newRow.ValidSignature = &boolSignature
			}

		}
		result.ResponseList = append(result.ResponseList, &newRow)
	}
	result.Closed = request.Closed
	result.TimedOut = request.TimedOut
	result.Mode = request.Mode

	// Check Role, If it's IdP then Set set special = true
	ownerRole := getRoleFromNodeID(request.Owner, app)
	if string(ownerRole) == "IdP" {
		result.Special = true
	}

	// Set requester_node_id
	result.RequesterNodeId = request.Owner
	resultJSON, err := proto.Marshal(&result)
	if err != nil {
		value = []byte("")
		return ReturnQuery(value, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(resultJSON, "success", app.state.db.Version64(), app)
}

func getNamespaceList(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetNamespaceList, Parameter: %s", param)
	key := "AllNamespace"
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)
	if value == nil {
		return ReturnQuery(nil, "not found", app.state.db.Version64(), app)
	}

	// result := make([]*data.Namespace, 0)
	// filter flag==true
	var result pbResult.GetNamespaceListResult
	var namespaces data.NamespaceList
	err := proto.Unmarshal([]byte(value), &namespaces)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	for _, namespace := range namespaces.Namespaces {
		if namespace.Active {
			var newRow pbResult.NamespaceInResult
			newRow.Namespace = namespace.Namespace
			newRow.Description = namespace.Description
			newRow.Active = namespace.Active
			result.Namespaces = append(result.Namespaces, &newRow)
		}
	}
	returnValue, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
}

func getServiceDetail(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetServiceDetail, Parameter: %s", param)
	var funcParam pbParam.GetServiceDetailParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	key := "Service" + "|" + funcParam.ServiceId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)
	if value == nil {
		return ReturnQuery(nil, "not found", app.state.db.Version64(), app)
	}
	var service data.ServiceDetail
	err = proto.Unmarshal(value, &service)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	var res pbResult.GetServiceDetailResult
	res.ServiceId = service.ServiceId
	res.ServiceName = service.ServiceName
	res.Active = service.Active
	returnValue, err := proto.Marshal(&res)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
}

func updateNode(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("UpdateNode, Parameter: %s", param)
	var funcParam pbParam.UpdateNodeParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	key := "NodeID" + "|" + nodeID
	_, value := app.state.db.Get(prefixKey([]byte(key)))

	if value != nil {
		var nodeDetail data.NodeDetail
		err := proto.Unmarshal([]byte(value), &nodeDetail)
		if err != nil {
			return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}

		// update MasterPublicKey
		if funcParam.MasterPublicKey != "" {
			nodeDetail.MasterPublicKey = funcParam.MasterPublicKey
		}

		// update PublicKey
		if funcParam.PublicKey != "" {
			nodeDetail.PublicKey = funcParam.PublicKey
		}

		nodeDetailValue, err := proto.Marshal(&nodeDetail)
		if err != nil {
			return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
		}
		app.SetStateDB([]byte(key), []byte(nodeDetailValue))
		return ReturnDeliverTxLog(code.OK, "success", "")
	}
	return ReturnDeliverTxLog(code.NodeIDNotFound, "Node ID not found", "")
}

func checkExistingIdentity(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("CheckExistingIdentity, Parameter: %s", param)
	var funcParam pbParam.CheckExistingIdentityParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var result pbResult.CheckExistingIdentityResult
	result.Exist = false

	key := "MsqDestination" + "|" + funcParam.HashId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)

	if value != nil {
		var nodes data.MsqDesList
		err = proto.Unmarshal([]byte(value), &nodes)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}

		msqCount := 0
		for _, node := range nodes.Nodes {
			if node.TimeoutBlock == 0 || node.TimeoutBlock > app.CurrentBlock {
				msqCount++
			}
		}

		if msqCount > 0 {
			result.Exist = true
		}
	}

	returnValue, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
}

func getAccessorGroupID(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetAccessorGroupID, Parameter: %s", param)
	var funcParam pbParam.GetAccessorGroupIDParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var result pbResult.GetAccessorGroupIDResult
	result.AccessorGroupId = ""

	key := "Accessor" + "|" + funcParam.AccessorId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)

	if value != nil {
		var accessor data.Accessor
		err = proto.Unmarshal([]byte(value), &accessor)
		if err == nil {
			result.AccessorGroupId = accessor.AccessorGroupId
		}
	}

	returnValue, err := proto.Marshal(&result)

	// If value == nil set log = "not found"
	if value == nil {
		return ReturnQuery([]byte("{}"), "not found", app.state.db.Version64(), app)
	}

	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
}

func getAccessorKey(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetAccessorKey, Parameter: %s", param)
	var funcParam pbParam.GetAccessorKeyParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var result pbResult.GetAccessorKeyResult
	result.AccessorPublicKey = ""

	key := "Accessor" + "|" + funcParam.AccessorId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)

	if value != nil {
		var accessor data.Accessor
		err = proto.Unmarshal([]byte(value), &accessor)
		if err == nil {
			result.AccessorPublicKey = accessor.AccessorPublicKey
			result.Active = accessor.Active
		}
	}

	returnValue, err := proto.Marshal(&result)

	// If value == nil set log = "not found"
	if value == nil {
		return ReturnQuery([]byte("{}"), "not found", app.state.db.Version64(), app)
	}

	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
}

func getServiceList(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetServiceList, Parameter: %s", param)
	key := "AllService"
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)
	if value == nil {
		var result pbResult.GetServiceListResult
		value, err := proto.Marshal(&result)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(value, "not found", app.state.db.Version64(), app)
	}

	var result pbResult.GetServiceListResult
	// filter flag==true
	var services data.ServiceDetailList
	err := proto.Unmarshal([]byte(value), &services)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	for _, service := range services.Services {
		if service.Active {
			var newRow pbResult.ServiceDetailInResult
			newRow.ServiceId = service.ServiceId
			newRow.ServiceName = service.ServiceName
			newRow.Active = service.Active
			result.Services = append(result.Services, &newRow)
		}
	}
	returnValue, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
}

func getServiceNameByServiceID(serviceID string, app *DIDApplication) string {
	key := "Service" + "|" + serviceID
	_, value := app.state.db.Get(prefixKey([]byte(key)))
	var result data.ServiceDetail
	if value != nil {
		err := proto.Unmarshal([]byte(value), &result)
		if err != nil {
			return ""
		}
		return result.ServiceName
	}
	return ""
}

func checkExistingAccessorID(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("CheckExistingAccessorID, Parameter: %s", param)
	var funcParam pbParam.CheckExistingAccessorIDParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var result pbResult.CheckExistingAccessorIDResult
	result.Exist = false

	accessorKey := "Accessor" + "|" + funcParam.AccessorId
	_, accessorValue := app.state.db.GetVersioned(prefixKey([]byte(accessorKey)), height)
	if accessorValue != nil {
		var accessor data.Accessor
		err = proto.Unmarshal([]byte(accessorValue), &accessor)
		if err == nil {
			result.Exist = true
		}
	}

	returnValue, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
}

func checkExistingAccessorGroupID(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("CheckExistingAccessorGroupID, Parameter: %s", param)
	var funcParam pbParam.CheckExistingAccessorGroupIDParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var result pbResult.CheckExistingAccessorIDResult
	result.Exist = false

	accessorGroupKey := "AccessorGroup" + "|" + funcParam.AccessorGroupId
	_, accessorGroupValue := app.state.db.GetVersioned(prefixKey([]byte(accessorGroupKey)), height)
	if accessorGroupValue != nil {
		result.Exist = true
	}

	returnValue, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
}

func getNodeInfo(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetNodeInfo, Parameter: %s", param)
	var funcParam pbParam.GetNodeInfoParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	nodeDetailKey := "NodeID" + "|" + funcParam.NodeId
	_, nodeDetailValue := app.state.db.GetVersioned(prefixKey([]byte(nodeDetailKey)), height)
	if nodeDetailValue == nil {
		return ReturnQuery([]byte("{}"), "not found", app.state.db.Version64(), app)
	}
	var nodeDetail data.NodeDetail
	err = proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	// If node behind proxy
	proxyKey := "Proxy" + "|" + funcParam.NodeId
	_, proxyValue := app.state.db.Get(prefixKey([]byte(proxyKey)))
	if proxyValue != nil {

		// Get proxy node ID
		var proxy data.Proxy
		err = proto.Unmarshal([]byte(proxyValue), &proxy)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		proxyNodeID := proxy.ProxyNodeId

		// Get proxy node detail
		proxyNodeDetailKey := "NodeID" + "|" + string(proxyNodeID)
		_, proxyNodeDetailValue := app.state.db.GetVersioned(prefixKey([]byte(proxyNodeDetailKey)), height)
		if proxyNodeDetailValue == nil {
			return ReturnQuery([]byte("{}"), "not found", app.state.db.Version64(), app)
		}
		var proxyNode data.NodeDetail
		err = proto.Unmarshal([]byte(proxyNodeDetailValue), &proxyNode)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		if nodeDetail.Role == "IdP" {
			var result pbResult.GetNodeInfoResult
			result.PublicKey = nodeDetail.PublicKey
			result.MasterPublicKey = nodeDetail.MasterPublicKey
			result.NodeName = nodeDetail.NodeName
			result.Role = nodeDetail.Role
			result.MaxIal = nodeDetail.MaxIal
			result.MaxAal = nodeDetail.MaxAal
			var proxyObj pbResult.ProxyInGetNodeInfoResult
			proxyObj.NodeId = string(proxyNodeID)
			proxyObj.NodeName = proxyNode.NodeName
			proxyObj.PublicKey = proxyNode.PublicKey
			proxyObj.MasterPublicKey = proxyNode.MasterPublicKey
			proxyObj.Config = proxy.Config
			if proxyNode.Mq != nil {
				for _, mq := range proxyNode.Mq {
					var msq pbResult.MsqAddressInResult
					msq.Ip = mq.Ip
					msq.Port = mq.Port
					proxyObj.Mq = append(proxyObj.Mq, &msq)
				}
			}
			result.Proxy = &proxyObj
			value, err := proto.Marshal(&result)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			return ReturnQuery(value, "success", app.state.db.Version64(), app)
		}
		var result pbResult.GetNodeInfoResult
		result.PublicKey = nodeDetail.PublicKey
		result.MasterPublicKey = nodeDetail.MasterPublicKey
		result.NodeName = nodeDetail.NodeName
		result.Role = nodeDetail.Role
		var proxyObj pbResult.ProxyInGetNodeInfoResult
		proxyObj.NodeId = string(proxyNodeID)
		proxyObj.NodeName = proxyNode.NodeName
		proxyObj.PublicKey = proxyNode.PublicKey
		proxyObj.MasterPublicKey = proxyNode.MasterPublicKey
		proxyObj.Config = proxy.Config
		if proxyNode.Mq != nil {
			for _, mq := range proxyNode.Mq {
				var msq pbResult.MsqAddressInResult
				msq.Ip = mq.Ip
				msq.Port = mq.Port
				proxyObj.Mq = append(proxyObj.Mq, &msq)
			}
		}
		result.Proxy = &proxyObj
		value, err := proto.Marshal(&result)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(value, "success", app.state.db.Version64(), app)
	} else {
		if nodeDetail.Role == "IdP" {
			var result pbResult.GetNodeInfoResult
			result.PublicKey = nodeDetail.PublicKey
			result.MasterPublicKey = nodeDetail.MasterPublicKey
			result.NodeName = nodeDetail.NodeName
			result.Role = nodeDetail.Role
			result.MaxIal = nodeDetail.MaxIal
			result.MaxAal = nodeDetail.MaxAal
			if nodeDetail.Mq != nil {
				for _, mq := range nodeDetail.Mq {
					var msq pbResult.MsqAddressInResult
					msq.Ip = mq.Ip
					msq.Port = mq.Port
					result.Mq = append(result.Mq, &msq)
				}
			}
			value, err := proto.Marshal(&result)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			return ReturnQuery(value, "success", app.state.db.Version64(), app)
		}
		var result pbResult.GetNodeInfoResult
		result.PublicKey = nodeDetail.PublicKey
		result.MasterPublicKey = nodeDetail.MasterPublicKey
		result.NodeName = nodeDetail.NodeName
		result.Role = nodeDetail.Role
		if nodeDetail.Mq != nil {
			for _, mq := range nodeDetail.Mq {
				var msq pbResult.MsqAddressInResult
				msq.Ip = mq.Ip
				msq.Port = mq.Port
				result.Mq = append(result.Mq, &msq)
			}
		}
		value, err := proto.Marshal(&result)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(value, "success", app.state.db.Version64(), app)
	}
}

func getIdentityInfo(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetIdentityInfo, Parameter: %s", param)
	var funcParam pbParam.GetIdentityInfoParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var result pbResult.GetIdentityInfoResult

	key := "MsqDestination" + "|" + funcParam.HashId
	_, chkExists := app.state.db.GetVersioned(prefixKey([]byte(key)), height)

	if chkExists != nil {
		var nodes data.MsqDesList
		err = proto.Unmarshal([]byte(chkExists), &nodes)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}

		for _, node := range nodes.Nodes {
			if node.NodeId == funcParam.NodeId {
				result.Ial = float64(node.Ial)
				break
			}
		}
	}

	returnValue, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	if result.Ial > 0.0 {
		return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
	}
	return ReturnQuery([]byte("{}"), "not found", app.state.db.Version64(), app)
}

func getDataSignature(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetDataSignature, Parameter: %s", param)
	var funcParam pbParam.GetDataSignatureParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	signDataKey := "SignData" + "|" + funcParam.NodeId + "|" + funcParam.ServiceId + "|" + funcParam.RequestId
	_, signDataValue := app.state.db.GetVersioned(prefixKey([]byte(signDataKey)), height)

	var result pbResult.GetDataSignatureResult

	if signDataValue != nil {
		result.Signature = string(signDataValue)
	}

	returnValue, err := proto.Marshal(&result)
	if signDataValue != nil {
		return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
	}
	return ReturnQuery([]byte("{}"), "not found", app.state.db.Version64(), app)
}

func getIdentityProof(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetIdentityProof, Parameter: %s", param)
	var funcParam pbParam.GetIdentityProofParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	identityProofKey := "IdentityProof" + "|" + funcParam.RequestId + "|" + funcParam.IdpId
	_, identityProofValue := app.state.db.GetVersioned(prefixKey([]byte(identityProofKey)), height)
	var result pbResult.GetIdentityProofResult
	if identityProofValue != nil {
		result.IdentityProof = string(identityProofValue)
	}
	returnValue, err := proto.Marshal(&result)
	if identityProofValue != nil {
		return ReturnQuery(returnValue, "success", app.state.db.Version64(), app)
	}
	return ReturnQuery(nil, "not found", app.state.db.Version64(), app)
}

func getServicesByAsID(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetServicesByAsID, Parameter: %s", param)
	var funcParam pbParam.GetServicesByAsIDParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var result pbResult.GetServicesByAsIDResult
	result.Services = make([]*pbResult.ServiceInResult, 0)

	provideServiceKey := "ProvideService" + "|" + funcParam.AsId
	_, provideServiceValue := app.state.db.Get(prefixKey([]byte(provideServiceKey)))
	var services data.ServiceList
	if provideServiceValue != nil {
		err := proto.Unmarshal([]byte(provideServiceValue), &services)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
	}

	nodeDetailKey := "NodeID" + "|" + funcParam.AsId
	_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
	var nodeDetail data.NodeDetail
	if nodeDetailValue != nil {
		err := proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
	}

	for index, provideService := range services.Services {
		serviceKey := "Service" + "|" + provideService.ServiceId
		_, serviceValue := app.state.db.Get(prefixKey([]byte(serviceKey)))
		var service data.ServiceDetail
		if serviceValue != nil {
			err = proto.Unmarshal([]byte(serviceValue), &service)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
		}
		if nodeDetail.Active && service.Active {
			// Set suspended from NDID
			approveServiceKey := "ApproveKey" + "|" + provideService.ServiceId + "|" + funcParam.AsId
			_, approveServiceJSON := app.state.db.Get(prefixKey([]byte(approveServiceKey)))
			if approveServiceJSON != nil {
				var approveService data.ApproveService
				err = proto.Unmarshal([]byte(approveServiceJSON), &approveService)
				if err == nil {
					services.Services[index].Suspended = !approveService.Active
				}
			}
			var newRow pbResult.ServiceInResult
			newRow.Active = services.Services[index].Active
			newRow.MinAal = services.Services[index].MinAal
			newRow.MinIal = services.Services[index].MinIal
			newRow.ServiceId = services.Services[index].ServiceId
			newRow.Suspended = &wrappers.BoolValue{Value: services.Services[index].Suspended}
			result.Services = append(result.Services, &newRow)
		}
	}

	resultJSON, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	if len(result.Services) > 0 {
		return ReturnQuery(resultJSON, "success", app.state.db.Version64(), app)
	} else {
		return ReturnQuery(resultJSON, "not found", app.state.db.Version64(), app)
	}
}

func getIdpNodesInfo(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetIdpNodesInfo, Parameter: %s", param)
	var funcParam pbParam.GetIdpNodesParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	var result pbResult.GetIdpNodesInfoResult
	result.Node = make([]*pbResult.IdpNode, 0)

	// Make mapping
	mapNodeIDList := map[string]bool{}
	for _, nodeID := range funcParam.NodeIdList {
		mapNodeIDList[nodeID] = true
	}

	if funcParam.HashId == "" {
		idpsKey := "IdPList"
		_, idpsValue := app.state.db.GetVersioned(prefixKey([]byte(idpsKey)), height)
		var idpsList data.IdPList
		if idpsValue != nil {
			err := proto.Unmarshal(idpsValue, &idpsList)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			for _, idp := range idpsList.NodeId {

				// filter from node_id_list
				if len(mapNodeIDList) > 0 {
					if mapNodeIDList[idp] == false {
						continue
					}
				}

				nodeDetailKey := "NodeID" + "|" + idp
				_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
				if nodeDetailValue == nil {
					continue
				}
				var nodeDetail data.NodeDetail
				err := proto.Unmarshal(nodeDetailValue, &nodeDetail)
				if err != nil {
					continue
				}
				// check node is active
				if !nodeDetail.Active {
					continue
				}
				// check Max IAL && AAL
				if !(nodeDetail.MaxIal >= funcParam.MinIal &&
					nodeDetail.MaxAal >= funcParam.MinAal) {
					continue
				}

				// If node is behind proxy
				proxyKey := "Proxy" + "|" + idp
				_, proxyValue := app.state.db.Get(prefixKey([]byte(proxyKey)))
				if proxyValue != nil {

					// Get proxy node ID
					var proxy data.Proxy
					err = proto.Unmarshal([]byte(proxyValue), &proxy)
					if err != nil {
						return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
					}
					proxyNodeID := proxy.ProxyNodeId

					// Get proxy node detail
					proxyNodeDetailKey := "NodeID" + "|" + string(proxyNodeID)
					_, proxyNodeDetailValue := app.state.db.GetVersioned(prefixKey([]byte(proxyNodeDetailKey)), height)
					if proxyNodeDetailValue == nil {
						return ReturnQuery([]byte("{}"), "not found", app.state.db.Version64(), app)
					}
					var proxyNode data.NodeDetail
					err = proto.Unmarshal([]byte(proxyNodeDetailValue), &proxyNode)
					if err != nil {
						return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
					}
					var msqDesNode pbResult.IdpNode
					msqDesNode.NodeId = idp
					msqDesNode.Name = nodeDetail.NodeName
					msqDesNode.MaxIal = nodeDetail.MaxIal
					msqDesNode.MaxAal = nodeDetail.MaxAal
					msqDesNode.PublicKey = nodeDetail.PublicKey
					var idpProxy pbResult.ProxyInResult
					idpProxy.NodeId = string(proxyNodeID)
					idpProxy.PublicKey = proxyNode.PublicKey
					idpProxy.Config = proxy.Config
					if proxyNode.Mq != nil {
						for _, mq := range proxyNode.Mq {
							var msq pbResult.MsqAddressInResult
							msq.Ip = mq.Ip
							msq.Port = mq.Port
							idpProxy.Mq = append(idpProxy.Mq, &msq)
						}
					}
					msqDesNode.Proxy = &idpProxy
					result.Node = append(result.Node, &msqDesNode)
				} else {
					var msq []*pbResult.MsqAddressInResult
					for _, mq := range nodeDetail.Mq {
						var msqAddress pbResult.MsqAddressInResult
						msqAddress.Ip = mq.Ip
						msqAddress.Port = mq.Port
						msq = append(msq, &msqAddress)
					}
					var msqDesNode pbResult.IdpNode
					msqDesNode.NodeId = idp
					msqDesNode.Name = nodeDetail.NodeName
					msqDesNode.MaxIal = nodeDetail.MaxIal
					msqDesNode.MaxAal = nodeDetail.MaxAal
					msqDesNode.PublicKey = nodeDetail.PublicKey
					msqDesNode.Mq = msq
					result.Node = append(result.Node, &msqDesNode)
				}

			}
		}
	} else {
		key := "MsqDestination" + "|" + funcParam.HashId
		_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)
		if value != nil {
			var nodes data.MsqDesList
			err = proto.Unmarshal([]byte(value), &nodes)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			for _, node := range nodes.Nodes {
				// filter from node_id_list
				if len(mapNodeIDList) > 0 {
					if mapNodeIDList[node.NodeId] == false {
						continue
					}
				}
				// check msq destination is not active
				if !node.Active {
					continue
				}
				// check Ial > min ial
				if node.Ial < funcParam.MinIal {
					continue
				}
				// check msq destination is not timed out
				if node.TimeoutBlock != 0 && app.CurrentBlock > node.TimeoutBlock {
					continue
				}
				nodeDetailKey := "NodeID" + "|" + node.NodeId
				_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
				if nodeDetailValue == nil {
					continue
				}
				var nodeDetail data.NodeDetail
				err := proto.Unmarshal(nodeDetailValue, &nodeDetail)
				if err != nil {
					continue
				}
				// check node is active
				if !nodeDetail.Active {
					continue
				}
				// check Max IAL && AAL
				if !(nodeDetail.MaxIal >= funcParam.MinIal &&
					nodeDetail.MaxAal >= funcParam.MinAal) {
					continue
				}

				// If node is behind proxy
				proxyKey := "Proxy" + "|" + node.NodeId
				_, proxyValue := app.state.db.Get(prefixKey([]byte(proxyKey)))
				if proxyValue != nil {

					// Get proxy node ID
					var proxy data.Proxy
					err = proto.Unmarshal([]byte(proxyValue), &proxy)
					if err != nil {
						return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
					}
					proxyNodeID := proxy.ProxyNodeId

					// Get proxy node detail
					proxyNodeDetailKey := "NodeID" + "|" + string(proxyNodeID)
					_, proxyNodeDetailValue := app.state.db.GetVersioned(prefixKey([]byte(proxyNodeDetailKey)), height)
					if proxyNodeDetailValue == nil {
						return ReturnQuery([]byte("{}"), "not found", app.state.db.Version64(), app)
					}
					var proxyNode data.NodeDetail
					err = proto.Unmarshal([]byte(proxyNodeDetailValue), &proxyNode)
					if err != nil {
						return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
					}
					var msqDesNode pbResult.IdpNode
					msqDesNode.NodeId = node.NodeId
					msqDesNode.Name = nodeDetail.NodeName
					msqDesNode.MaxIal = nodeDetail.MaxIal
					msqDesNode.MaxAal = nodeDetail.MaxAal
					msqDesNode.PublicKey = nodeDetail.PublicKey
					var idpProxy pbResult.ProxyInResult
					idpProxy.NodeId = string(proxyNodeID)
					idpProxy.PublicKey = proxyNode.PublicKey
					idpProxy.Config = proxy.Config
					if proxyNode.Mq != nil {
						for _, mq := range proxyNode.Mq {
							var msq pbResult.MsqAddressInResult
							msq.Ip = mq.Ip
							msq.Port = mq.Port
							idpProxy.Mq = append(idpProxy.Mq, &msq)
						}
					}
					msqDesNode.Proxy = &idpProxy
					result.Node = append(result.Node, &msqDesNode)
				} else {
					var msq []*pbResult.MsqAddressInResult
					for _, mq := range nodeDetail.Mq {
						var msqAddress pbResult.MsqAddressInResult
						msqAddress.Ip = mq.Ip
						msqAddress.Port = mq.Port
						msq = append(msq, &msqAddress)
					}
					var msqDesNode pbResult.IdpNode
					msqDesNode.NodeId = node.NodeId
					msqDesNode.Name = nodeDetail.NodeName
					msqDesNode.MaxIal = nodeDetail.MaxIal
					msqDesNode.MaxAal = nodeDetail.MaxAal
					msqDesNode.PublicKey = nodeDetail.PublicKey
					msqDesNode.Mq = msq
					result.Node = append(result.Node, &msqDesNode)
				}
			}
		}
	}

	value, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	if len(result.Node) > 0 {
		return ReturnQuery(value, "success", app.state.db.Version64(), app)
	}
	return ReturnQuery(value, "not found", app.state.db.Version64(), app)
}

func getAsNodesInfoByServiceId(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetAsNodesInfoByServiceId, Parameter: %s", param)
	var funcParam pbParam.GetAsNodesByServiceIdParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	key := "ServiceDestination" + "|" + funcParam.ServiceId
	_, value := app.state.db.GetVersioned(prefixKey([]byte(key)), height)

	if value == nil {
		var result pbResult.GetAsNodesInfoByServiceIdResult
		result.Node = make([]*pbResult.ASWithMqNode, 0)
		value, err := proto.Marshal(&result)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(value, "not found", app.state.db.Version64(), app)
	}

	// filter serive is active
	serviceKey := "Service" + "|" + funcParam.ServiceId
	_, serviceValue := app.state.db.Get(prefixKey([]byte(serviceKey)))
	if serviceValue != nil {
		var service data.ServiceDetail
		err = proto.Unmarshal([]byte(serviceValue), &service)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		if service.Active == false {
			var result pbResult.GetAsNodesInfoByServiceIdResult
			result.Node = make([]*pbResult.ASWithMqNode, 0)
			value, err := proto.Marshal(&result)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			return ReturnQuery(value, "service is not active", app.state.db.Version64(), app)
		}
	} else {
		var result pbResult.GetAsNodesInfoByServiceIdResult
		result.Node = make([]*pbResult.ASWithMqNode, 0)
		value, err := proto.Marshal(&result)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(value, "not found", app.state.db.Version64(), app)
	}

	var storedData data.ServiceDesList
	err = proto.Unmarshal([]byte(value), &storedData)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	// Make mapping
	mapNodeIDList := map[string]bool{}
	for _, nodeID := range funcParam.NodeIdList {
		mapNodeIDList[nodeID] = true
	}

	var result pbResult.GetAsNodesInfoByServiceIdResult
	result.Node = make([]*pbResult.ASWithMqNode, 0)
	for index := range storedData.Node {

		// filter from node_id_list
		if len(mapNodeIDList) > 0 {
			if mapNodeIDList[storedData.Node[index].NodeId] == false {
				continue
			}
		}

		// filter service destination is Active
		if !storedData.Node[index].Active {
			continue
		}

		// Filter approve from NDID
		approveServiceKey := "ApproveKey" + "|" + funcParam.ServiceId + "|" + storedData.Node[index].NodeId
		_, approveServiceJSON := app.state.db.Get(prefixKey([]byte(approveServiceKey)))
		if approveServiceJSON == nil {
			continue
		}
		var approveService data.ApproveService
		err = proto.Unmarshal([]byte(approveServiceJSON), &approveService)
		if err != nil {
			continue
		}
		if !approveService.Active {
			continue
		}

		nodeDetailKey := "NodeID" + "|" + storedData.Node[index].NodeId
		_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
		if nodeDetailValue == nil {
			continue
		}
		var nodeDetail data.NodeDetail
		err := proto.Unmarshal(nodeDetailValue, &nodeDetail)
		if err != nil {
			continue
		}
		// filter node is active
		if !nodeDetail.Active {
			continue
		}

		// If node is behind proxy
		proxyKey := "Proxy" + "|" + storedData.Node[index].NodeId
		_, proxyValue := app.state.db.Get(prefixKey([]byte(proxyKey)))
		if proxyValue != nil {

			// Get proxy node ID
			var proxy data.Proxy
			err = proto.Unmarshal([]byte(proxyValue), &proxy)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			proxyNodeID := proxy.ProxyNodeId

			// Get proxy node detail
			proxyNodeDetailKey := "NodeID" + "|" + string(proxyNodeID)
			_, proxyNodeDetailValue := app.state.db.GetVersioned(prefixKey([]byte(proxyNodeDetailKey)), height)
			if proxyNodeDetailValue == nil {
				return ReturnQuery([]byte("{}"), "not found", app.state.db.Version64(), app)
			}
			var proxyNode data.NodeDetail
			err = proto.Unmarshal([]byte(proxyNodeDetailValue), &proxyNode)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			var as pbResult.ASWithMqNode
			as.NodeId = storedData.Node[index].NodeId
			as.Name = nodeDetail.NodeName
			as.MinIal = storedData.Node[index].MinIal
			as.MinAal = storedData.Node[index].MinAal
			as.PublicKey = nodeDetail.PublicKey
			var proxyAS pbResult.ProxyInResult
			proxyAS.NodeId = string(proxyNodeID)
			proxyAS.PublicKey = proxyNode.PublicKey
			proxyAS.Config = proxy.Config
			if proxyNode.Mq != nil {
				for _, mq := range proxyNode.Mq {
					var msq pbResult.MsqAddressInResult
					msq.Ip = mq.Ip
					msq.Port = mq.Port
					proxyAS.Mq = append(proxyAS.Mq, &msq)
				}
			}
			as.Proxy = &proxyAS
			result.Node = append(result.Node, &as)
		} else {
			var msqAddress []*pbResult.MsqAddressInResult
			for _, mq := range nodeDetail.Mq {
				var msq pbResult.MsqAddressInResult
				msq.Ip = mq.Ip
				msq.Port = mq.Port
				msqAddress = append(msqAddress, &msq)
			}
			var as pbResult.ASWithMqNode
			as.NodeId = storedData.Node[index].NodeId
			as.Name = nodeDetail.NodeName
			as.MinIal = storedData.Node[index].MinIal
			as.MinAal = storedData.Node[index].MinAal
			as.PublicKey = nodeDetail.PublicKey
			as.Mq = msqAddress
			result.Node = append(result.Node, &as)
		}
	}
	resultJSON, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	return ReturnQuery(resultJSON, "success", app.state.db.Version64(), app)
}

func getNodesBehindProxyNode(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetNodesBehindProxyNode, Parameter: %s", param)
	var funcParam pbParam.GetNodesBehindProxyNodeParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	var result pbResult.GetNodesBehindProxyNodeResult
	result.Nodes = make([]*pbResult.NodeBehindProxy, 0)
	behindProxyNodeKey := "BehindProxyNode" + "|" + funcParam.ProxyNodeId
	_, behindProxyNodeValue := app.state.db.Get(prefixKey([]byte(behindProxyNodeKey)))
	if behindProxyNodeValue == nil {
		resultJSON, err := proto.Marshal(&result)
		if err != nil {
			return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
		}
		return ReturnQuery(resultJSON, "not found", app.state.db.Version64(), app)
	}
	var nodes data.BehindNodeList
	nodes.Nodes = make([]string, 0)
	err = proto.Unmarshal([]byte(behindProxyNodeValue), &nodes)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	for _, node := range nodes.Nodes {
		nodeDetailKey := "NodeID" + "|" + node
		_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
		if nodeDetailValue == nil {
			continue
		}
		var nodeDetail data.NodeDetail
		err := proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
		if err != nil {
			continue
		}

		// Get proxy detail
		proxyKey := "Proxy" + "|" + node
		_, proxyValue := app.state.db.Get(prefixKey([]byte(proxyKey)))
		if proxyValue == nil {
			continue
		}
		var proxy data.Proxy
		err = proto.Unmarshal([]byte(proxyValue), &proxy)
		if err != nil {
			continue
		}

		if nodeDetail.Role == "IdP" {
			var row pbResult.NodeBehindProxy
			row.NodeId = node
			row.NodeName = nodeDetail.NodeName
			row.Role = nodeDetail.Role
			row.PublicKey = nodeDetail.PublicKey
			row.MasterPublicKey = nodeDetail.MasterPublicKey
			row.MaxIal = nodeDetail.MaxIal
			row.MaxAal = nodeDetail.MaxAal
			row.Config = proxy.Config
			result.Nodes = append(result.Nodes, &row)
		} else {
			var row pbResult.NodeBehindProxy
			row.NodeId = node
			row.NodeName = nodeDetail.NodeName
			row.Role = nodeDetail.Role
			row.PublicKey = nodeDetail.PublicKey
			row.MasterPublicKey = nodeDetail.MasterPublicKey
			row.Config = proxy.Config
			result.Nodes = append(result.Nodes, &row)
		}

	}
	resultJSON, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	if len(result.Nodes) == 0 {
		return ReturnQuery(resultJSON, "not found", app.state.db.Version64(), app)
	}
	return ReturnQuery(resultJSON, "success", app.state.db.Version64(), app)
}

func getNodeIDList(param []byte, app *DIDApplication, height int64) types.ResponseQuery {
	app.logger.Infof("GetNodeIDList, Parameter: %s", param)
	var funcParam pbParam.GetNodeIDListParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}

	var result pbResult.GetNodeIDListResult
	result.NodeIdList = make([]string, 0)

	if strings.ToLower(funcParam.Role) == "rp" {
		var rpsList data.RPList
		rpsKey := "rpList"
		_, rpsValue := app.state.db.Get(prefixKey([]byte(rpsKey)))
		if rpsValue != nil {
			err := proto.Unmarshal(rpsValue, &rpsList)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			for _, nodeID := range rpsList.NodeId {
				nodeDetailKey := "NodeID" + "|" + nodeID
				_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
				if nodeDetailValue != nil {
					var nodeDetail data.NodeDetail
					err := proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
					if err != nil {
						continue
					}
					if nodeDetail.Active {
						result.NodeIdList = append(result.NodeIdList, nodeID)
					}
				}
			}
		}
	} else if strings.ToLower(funcParam.Role) == "idp" {
		var idpsList data.IdPList
		idpsKey := "IdPList"
		_, idpsValue := app.state.db.Get(prefixKey([]byte(idpsKey)))
		if idpsValue != nil {
			err := proto.Unmarshal(idpsValue, &idpsList)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			for _, nodeID := range idpsList.NodeId {
				nodeDetailKey := "NodeID" + "|" + nodeID
				_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
				if nodeDetailValue != nil {
					var nodeDetail data.NodeDetail
					err := proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
					if err != nil {
						continue
					}
					if nodeDetail.Active {
						result.NodeIdList = append(result.NodeIdList, nodeID)
					}
				}
			}
		}
	} else if strings.ToLower(funcParam.Role) == "as" {
		var asList data.ASList
		asKey := "asList"
		_, asValue := app.state.db.Get(prefixKey([]byte(asKey)))
		if asValue != nil {
			err := proto.Unmarshal(asValue, &asList)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			for _, nodeID := range asList.NodeId {
				nodeDetailKey := "NodeID" + "|" + nodeID
				_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
				if nodeDetailValue != nil {
					var nodeDetail data.NodeDetail
					err := proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
					if err != nil {
						continue
					}
					if nodeDetail.Active {
						result.NodeIdList = append(result.NodeIdList, nodeID)
					}
				}
			}
		}
	} else {
		var allList data.AllList
		allKey := "allList"
		_, allValue := app.state.db.Get(prefixKey([]byte(allKey)))
		if allValue != nil {
			err := proto.Unmarshal(allValue, &allList)
			if err != nil {
				return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
			}
			for _, nodeID := range allList.NodeId {
				nodeDetailKey := "NodeID" + "|" + nodeID
				_, nodeDetailValue := app.state.db.Get(prefixKey([]byte(nodeDetailKey)))
				if nodeDetailValue != nil {
					var nodeDetail data.NodeDetail
					err := proto.Unmarshal([]byte(nodeDetailValue), &nodeDetail)
					if err != nil {
						continue
					}
					if nodeDetail.Active {
						result.NodeIdList = append(result.NodeIdList, nodeID)
					}
				}
			}
		}
	}

	resultJSON, err := proto.Marshal(&result)
	if err != nil {
		return ReturnQuery(nil, err.Error(), app.state.db.Version64(), app)
	}
	if len(result.NodeIdList) == 0 {
		return ReturnQuery(resultJSON, "not found", app.state.db.Version64(), app)
	}
	return ReturnQuery(resultJSON, "success", app.state.db.Version64(), app)
}
