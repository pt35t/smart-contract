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
	pbParam "github.com/ndidplatform/smart-contract/protos/params"
	"github.com/tendermint/tendermint/abci/types"
)

func signData(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("SignData, Parameter: %s", param)
	var signData pbParam.SignDataParams
	err := proto.Unmarshal([]byte(param), &signData)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	requestKey := "Request" + "|" + signData.RequestId
	_, requestJSON := app.state.db.Get(prefixKey([]byte(requestKey)))
	if requestJSON == nil {
		return ReturnDeliverTxLog(code.RequestIDNotFound, "Request ID not found", "")
	}
	var request data.Request
	err = proto.Unmarshal([]byte(requestJSON), &request)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check IsClosed
	if request.Closed {
		return ReturnDeliverTxLog(code.RequestIsClosed, "Request is closed", "")
	}

	// Check IsTimedOut
	if request.TimedOut {
		return ReturnDeliverTxLog(code.RequestIsTimedOut, "Request is timed out", "")
	}

	// Check Service ID
	serviceKey := "Service" + "|" + signData.ServiceId
	_, serviceJSON := app.state.db.Get(prefixKey([]byte(serviceKey)))
	if serviceJSON == nil {
		return ReturnDeliverTxLog(code.ServiceIDNotFound, "Service ID not found", "")
	}
	var service data.ServiceDetail
	err = proto.Unmarshal([]byte(serviceJSON), &service)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check service is active
	if !service.Active {
		return ReturnDeliverTxLog(code.ServiceIsNotActive, "Service is not active", "")
	}

	// Check service destination is approved by NDID
	approveServiceKey := "ApproveKey" + "|" + signData.ServiceId + "|" + nodeID
	_, approveServiceJSON := app.state.db.Get(prefixKey([]byte(approveServiceKey)))
	if approveServiceJSON == nil {
		return ReturnDeliverTxLog(code.ServiceIDNotFound, "Service ID not found", "")
	}
	var approveService data.ApproveService
	err = proto.Unmarshal([]byte(approveServiceJSON), &approveService)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	if !approveService.Active {
		return ReturnDeliverTxLog(code.ServiceDestinationIsNotActive, "Service destination is not approved by NDID", "")
	}

	// Check service destination is active
	serviceDestinationKey := "ServiceDestination" + "|" + signData.ServiceId
	_, serviceDestinationValue := app.state.db.Get(prefixKey([]byte(serviceDestinationKey)))

	if serviceDestinationValue == nil {
		return ReturnDeliverTxLog(code.ServiceDestinationNotFound, "Service destination not found", "")
	}

	var nodes data.ServiceDesList
	err = proto.Unmarshal([]byte(serviceDestinationValue), &nodes)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	for index := range nodes.Node {
		if nodes.Node[index].NodeId == nodeID {
			if !nodes.Node[index].Active {
				return ReturnDeliverTxLog(code.ServiceDestinationIsNotActive, "Service destination is not active", "")
			}
			break
		}
	}

	// if AS != [], Check nodeID is exist in as_id_list
	exist := false
	for _, dataRequest := range request.DataRequestList {
		if dataRequest.ServiceId == signData.ServiceId {
			if len(dataRequest.AsIdList) == 0 {
				exist = true
				break
			} else {
				for _, as := range dataRequest.AsIdList {
					if as == nodeID {
						exist = true
						break
					}
				}
			}
		}
	}
	if exist == false {
		return ReturnDeliverTxLog(code.NodeIDIsNotExistInASList, "Node ID is not exist in AS list", "")
	}

	// Check Duplicate AS ID
	duplicate := false
	for _, dataRequest := range request.DataRequestList {
		if dataRequest.ServiceId == signData.ServiceId {
			for _, as := range dataRequest.AnsweredAsIdList {
				if as == nodeID {
					duplicate = true
					break
				}
			}
		}
	}
	if duplicate == true {
		return ReturnDeliverTxLog(code.DuplicateAnsweredAsIDList, "Duplicate AS ID in answered AS list", "")
	}

	// Check min_as
	for _, dataRequest := range request.DataRequestList {
		if dataRequest.ServiceId == signData.ServiceId {
			if int64(len(dataRequest.AnsweredAsIdList)) >= dataRequest.MinAs {
				return ReturnDeliverTxLog(code.DataRequestIsCompleted, "Can't sign data to data request that's enough data", "")
			}
		}
	}

	signDataKey := "SignData" + "|" + nodeID + "|" + signData.ServiceId + "|" + signData.RequestId
	signDataValue := signData.Signature

	// Update answered_as_id_list in request
	for index, dataRequest := range request.DataRequestList {
		if dataRequest.ServiceId == signData.ServiceId {
			request.DataRequestList[index].AnsweredAsIdList = append(dataRequest.AnsweredAsIdList, nodeID)
		}
	}

	requestJSON, err = proto.Marshal(&request)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}

	app.SetStateDB([]byte(requestKey), []byte(requestJSON))
	app.SetStateDB([]byte(signDataKey), []byte(signDataValue))
	return ReturnDeliverTxLog(code.OK, "success", signData.RequestId)
}

func registerServiceDestination(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("RegisterServiceDestination, Parameter: %s", param)
	var funcParam pbParam.RegisterServiceDestinationParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check Service ID
	serviceKey := "Service" + "|" + funcParam.ServiceId
	_, serviceJSON := app.state.db.Get(prefixKey([]byte(serviceKey)))
	if serviceJSON == nil {
		return ReturnDeliverTxLog(code.ServiceIDNotFound, "Service ID not found", "")
	}
	var service data.ServiceDetail
	err = proto.Unmarshal([]byte(serviceJSON), &service)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check service is active
	if !service.Active {
		return ReturnDeliverTxLog(code.ServiceIsNotActive, "Service is not active", "")
	}

	provideServiceKey := "ProvideService" + "|" + nodeID
	_, provideServiceValue := app.state.db.Get(prefixKey([]byte(provideServiceKey)))
	var services data.ServiceList
	if provideServiceValue != nil {
		err := proto.Unmarshal([]byte(provideServiceValue), &services)
		if err != nil {
			return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
	}
	// Check duplicate service ID
	for _, service := range services.Services {
		if service.ServiceId == funcParam.ServiceId {
			return ReturnDeliverTxLog(code.DuplicateServiceID, "Duplicate service ID in provide service list", "")
		}
	}

	// Check approve register service destination from NDID
	approveServiceKey := "ApproveKey" + "|" + funcParam.ServiceId + "|" + nodeID
	_, approveServiceJSON := app.state.db.Get(prefixKey([]byte(approveServiceKey)))
	if approveServiceJSON == nil {
		return ReturnDeliverTxLog(code.NoPermissionForRegisterServiceDestination, "This node does not have permission to register service destination", "")
	}
	var approveService data.ApproveService
	err = proto.Unmarshal([]byte(approveServiceJSON), &approveService)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}
	if approveService.Active == false {
		return ReturnDeliverTxLog(code.NoPermissionForRegisterServiceDestination, "This node does not have permission to register service destination", "")
	}

	// Append to ProvideService list
	var newService data.Service
	newService.ServiceId = funcParam.ServiceId
	newService.MinAal = funcParam.MinAal
	newService.MinIal = funcParam.MinIal
	newService.Active = true
	services.Services = append(services.Services, &newService)

	provideServiceJSON, err := proto.Marshal(&services)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}

	// Add ServiceDestination
	serviceDestinationKey := "ServiceDestination" + "|" + funcParam.ServiceId
	_, chkExists := app.state.db.Get(prefixKey([]byte(serviceDestinationKey)))

	if chkExists != nil {
		var nodes data.ServiceDesList
		err := proto.Unmarshal([]byte(chkExists), &nodes)
		if err != nil {
			return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}

		// Check duplicate node ID before add
		for _, node := range nodes.Node {
			if node.NodeId == nodeID {
				return ReturnDeliverTxLog(code.DuplicateNodeID, "Duplicate node ID", "")
			}
		}

		var newNode data.ASNode
		newNode.NodeId = nodeID
		newNode.MinIal = funcParam.MinIal
		newNode.MinAal = funcParam.MinAal
		newNode.ServiceId = funcParam.ServiceId
		newNode.Active = true
		nodes.Node = append(nodes.Node, &newNode)
		value, err := proto.Marshal(&nodes)
		if err != nil {
			return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
		}
		app.SetStateDB([]byte(serviceDestinationKey), []byte(value))
	} else {
		var nodes data.ServiceDesList
		var newNode data.ASNode
		newNode.NodeId = nodeID
		newNode.MinIal = funcParam.MinIal
		newNode.MinAal = funcParam.MinAal
		newNode.ServiceId = funcParam.ServiceId
		newNode.Active = true
		nodes.Node = append(nodes.Node, &newNode)
		value, err := proto.Marshal(&nodes)
		if err != nil {
			return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
		}
		app.SetStateDB([]byte(serviceDestinationKey), []byte(value))
	}
	app.SetStateDB([]byte(provideServiceKey), []byte(provideServiceJSON))
	return ReturnDeliverTxLog(code.OK, "success", "")
}

func updateServiceDestination(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("UpdateServiceDestination, Parameter: %s", param)
	var funcParam pbParam.UpdateServiceDestinationParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check Service ID
	serviceKey := "Service" + "|" + funcParam.ServiceId
	_, serviceJSON := app.state.db.Get(prefixKey([]byte(serviceKey)))
	if serviceJSON == nil {
		return ReturnDeliverTxLog(code.ServiceIDNotFound, "Service ID not found", "")
	}
	var service data.ServiceDetail
	err = proto.Unmarshal([]byte(serviceJSON), &service)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Update ServiceDestination
	serviceDestinationKey := "ServiceDestination" + "|" + funcParam.ServiceId
	_, serviceDestinationValue := app.state.db.Get(prefixKey([]byte(serviceDestinationKey)))

	if serviceDestinationValue == nil {
		return ReturnDeliverTxLog(code.ServiceDestinationNotFound, "Service destination not found", "")
	}

	var nodes data.ServiceDesList
	err = proto.Unmarshal([]byte(serviceDestinationValue), &nodes)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	for index := range nodes.Node {
		if nodes.Node[index].NodeId == nodeID {
			// selective update
			if funcParam.MinAal > 0 {
				nodes.Node[index].MinAal = funcParam.MinAal
			}
			if funcParam.MinIal > 0 {
				nodes.Node[index].MinIal = funcParam.MinIal
			}
			break
		}
	}

	// Update PrivideService
	provideServiceKey := "ProvideService" + "|" + nodeID
	_, provideServiceValue := app.state.db.Get(prefixKey([]byte(provideServiceKey)))
	var services data.ServiceList
	if provideServiceValue != nil {
		err := proto.Unmarshal([]byte(provideServiceValue), &services)
		if err != nil {
			return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
	}
	for index, service := range services.Services {
		if service.ServiceId == funcParam.ServiceId {
			if funcParam.MinAal > 0 {
				services.Services[index].MinAal = funcParam.MinAal
			}
			if funcParam.MinIal > 0 {
				services.Services[index].MinIal = funcParam.MinIal
			}
			break
		}
	}
	provideServiceJSON, err := proto.Marshal(&services)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	serviceDestinationJSON, err := proto.Marshal(&nodes)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(provideServiceKey), []byte(provideServiceJSON))
	app.SetStateDB([]byte(serviceDestinationKey), []byte(serviceDestinationJSON))
	return ReturnDeliverTxLog(code.OK, "success", "")
}

func disableServiceDestination(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("DisableServiceDestination, Parameter: %s", param)
	var funcParam pbParam.DisableServiceDestinationParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check Service ID
	serviceKey := "Service" + "|" + funcParam.ServiceId
	_, serviceJSON := app.state.db.Get(prefixKey([]byte(serviceKey)))
	if serviceJSON == nil {
		return ReturnDeliverTxLog(code.ServiceIDNotFound, "Service ID not found", "")
	}
	var service data.ServiceDetail
	err = proto.Unmarshal([]byte(serviceJSON), &service)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Update ServiceDestination
	serviceDestinationKey := "ServiceDestination" + "|" + funcParam.ServiceId
	_, serviceDestinationValue := app.state.db.Get(prefixKey([]byte(serviceDestinationKey)))

	if serviceDestinationValue == nil {
		return ReturnDeliverTxLog(code.ServiceDestinationNotFound, "Service destination not found", "")
	}

	var nodes data.ServiceDesList
	err = proto.Unmarshal([]byte(serviceDestinationValue), &nodes)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	for index := range nodes.Node {
		if nodes.Node[index].NodeId == nodeID {
			nodes.Node[index].Active = false
			break
		}
	}

	// Update ProvideService
	provideServiceKey := "ProvideService" + "|" + nodeID
	_, provideServiceValue := app.state.db.Get(prefixKey([]byte(provideServiceKey)))
	var services data.ServiceList
	if provideServiceValue != nil {
		err := proto.Unmarshal([]byte(provideServiceValue), &services)
		if err != nil {
			return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
	}
	for index, service := range services.Services {
		if service.ServiceId == funcParam.ServiceId {
			services.Services[index].Active = false
			break
		}
	}
	provideServiceJSON, err := proto.Marshal(&services)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}

	serviceDestinationJSON, err := proto.Marshal(&nodes)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(provideServiceKey), []byte(provideServiceJSON))
	app.SetStateDB([]byte(serviceDestinationKey), []byte(serviceDestinationJSON))
	return ReturnDeliverTxLog(code.OK, "success", "")
}

func enableServiceDestination(param []byte, app *DIDApplication, nodeID string) types.ResponseDeliverTx {
	app.logger.Infof("EnableServiceDestination, Parameter: %s", param)
	var funcParam pbParam.DisableServiceDestinationParams
	err := proto.Unmarshal([]byte(param), &funcParam)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Check Service ID
	serviceKey := "Service" + "|" + funcParam.ServiceId
	_, serviceJSON := app.state.db.Get(prefixKey([]byte(serviceKey)))
	if serviceJSON == nil {
		return ReturnDeliverTxLog(code.ServiceIDNotFound, "Service ID not found", "")
	}
	var service data.ServiceDetail
	err = proto.Unmarshal([]byte(serviceJSON), &service)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	// Update ServiceDestination
	serviceDestinationKey := "ServiceDestination" + "|" + funcParam.ServiceId
	_, serviceDestinationValue := app.state.db.Get(prefixKey([]byte(serviceDestinationKey)))

	if serviceDestinationValue == nil {
		return ReturnDeliverTxLog(code.ServiceDestinationNotFound, "Service destination not found", "")
	}

	var nodes data.ServiceDesList
	err = proto.Unmarshal([]byte(serviceDestinationValue), &nodes)
	if err != nil {
		return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
	}

	for index := range nodes.Node {
		if nodes.Node[index].NodeId == nodeID {
			nodes.Node[index].Active = true
			break
		}
	}

	// Update ProvideService
	provideServiceKey := "ProvideService" + "|" + nodeID
	_, provideServiceValue := app.state.db.Get(prefixKey([]byte(provideServiceKey)))
	var services data.ServiceList
	if provideServiceValue != nil {
		err := proto.Unmarshal([]byte(provideServiceValue), &services)
		if err != nil {
			return ReturnDeliverTxLog(code.UnmarshalError, err.Error(), "")
		}
	}
	for index, service := range services.Services {
		if service.ServiceId == funcParam.ServiceId {
			services.Services[index].Active = true
			break
		}
	}
	provideServiceJSON, err := proto.Marshal(&services)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}

	serviceDestinationJSON, err := proto.Marshal(&nodes)
	if err != nil {
		return ReturnDeliverTxLog(code.MarshalError, err.Error(), "")
	}
	app.SetStateDB([]byte(provideServiceKey), []byte(provideServiceJSON))
	app.SetStateDB([]byte(serviceDestinationKey), []byte(serviceDestinationJSON))
	return ReturnDeliverTxLog(code.OK, "success", "")
}
