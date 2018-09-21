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

type NodePublicKey struct {
	NodeID    string `json:"node_id"`
	PublicKey string `json:"public_key"`
}

// type GetNodePublicKeyResult struct {
// 	PublicKey string `json:"public_key"`
// }

// type GetNodeMasterPublicKeyResult struct {
// 	MasterPublicKey string `json:"master_public_key"`
// }

type Node struct {
	Ial          float64 `json:"ial"`
	NodeID       string  `json:"node_id"`
	Active       bool    `json:"active"`
	First        bool    `json:"first"`
	TimeoutBlock int64   `json:"timeout_block"`
}

type MsqDestination struct {
	Nodes []Node `json:"nodes"`
}

// type MsqDestinationNode struct {
// 	ID     string  `json:"node_id"`
// 	Name   string  `json:"node_name"`
// 	MaxIal float64 `json:"max_ial"`
// 	MaxAal float64 `json:"max_aal"`
// }

// type GetIdpNodesResult struct {
// 	Node []MsqDestinationNode `json:"node"`
// }

type GetAccessorMethodParam struct {
	AccessorID string `json:"accessor_id"`
}

type GetAccessorMethodResult struct {
	AccessorType string `json:"accessor_type"`
	AccessorKey  string `json:"accessor_key"`
	Commitment   string `json:"commitment"`
}

type DataRequest struct {
	ServiceID            string   `json:"service_id"`
	As                   []string `json:"as_id_list"`
	Count                int      `json:"min_as"`
	RequestParamsHash    string   `json:"request_params_hash"`
	AnsweredAsIdList     []string `json:"answered_as_id_list"`
	ReceivedDataFromList []string `json:"received_data_from_list"`
}

type Request struct {
	RequestID       string        `json:"request_id"`
	MinIdp          int           `json:"min_idp"`
	MinAal          float64       `json:"min_aal"`
	MinIal          float64       `json:"min_ial"`
	Timeout         int           `json:"request_timeout"`
	DataRequestList []DataRequest `json:"data_request_list"`
	MessageHash     string        `json:"request_message_hash"`
	Responses       []Response    `json:"response_list"`
	IsClosed        bool          `json:"closed"`
	IsTimedOut      bool          `json:"timed_out"`
	CanAddAccessor  bool          `json:"can_add_accessor"`
	Owner           string        `json:"owner"`
	Mode            int           `json:"mode"`
}

type Response struct {
	Ial              float64 `json:"ial"`
	Aal              float64 `json:"aal"`
	Status           string  `json:"status"`
	Signature        string  `json:"signature"`
	IdentityProof    string  `json:"identity_proof"`
	PrivateProofHash string  `json:"private_proof_hash"`
	IdpID            string  `json:"idp_id"`
	ValidProof       *bool   `json:"valid_proof"`
	ValidIal         *bool   `json:"valid_ial"`
	ValidSignature   *bool   `json:"valid_signature"`
}

type GetRequestResult struct {
	IsClosed    bool   `json:"closed"`
	IsTimedOut  bool   `json:"timed_out"`
	MessageHash string `json:"request_message_hash"`
	Mode        int    `json:"mode"`
}

type GetRequestDetailResult struct {
	RequestID       string        `json:"request_id"`
	MinIdp          int           `json:"min_idp"`
	MinAal          float64       `json:"min_aal"`
	MinIal          float64       `json:"min_ial"`
	Timeout         int           `json:"request_timeout"`
	DataRequestList []DataRequest `json:"data_request_list"`
	MessageHash     string        `json:"request_message_hash"`
	Responses       []Response    `json:"response_list"`
	IsClosed        bool          `json:"closed"`
	IsTimedOut      bool          `json:"timed_out"`
	Special         bool          `json:"special"`
	Mode            int           `json:"mode"`
	RequesterNodeID string        `json:"requester_node_id"`
}

type ASNode struct {
	ID        string  `json:"node_id"`
	Name      string  `json:"node_name"`
	MinIal    float64 `json:"min_ial"`
	MinAal    float64 `json:"min_aal"`
	ServiceID string  `json:"service_id"`
	Active    bool    `json:"active"`
}

type GetAsNodesByServiceIdResult struct {
	Node []ASNode `json:"node"`
}

type ASNodeResult struct {
	ID     string  `json:"node_id"`
	Name   string  `json:"node_name"`
	MinIal float64 `json:"min_ial"`
	MinAal float64 `json:"min_aal"`
}

type GetAsNodesByServiceIdWithNameResult struct {
	Node []ASNodeResult `json:"node"`
}

type InitNDIDParam struct {
	NodeID          string `json:"node_id"`
	PublicKey       string `json:"public_key"`
	MasterPublicKey string `json:"master_public_key"`
}

type TransferNDIDParam struct {
	PublicKey string `json:"public_key"`
}

type RegisterNode struct {
	NodeID          string  `json:"node_id"`
	PublicKey       string  `json:"public_key"`
	MasterPublicKey string  `json:"master_public_key"`
	NodeName        string  `json:"node_name"`
	Role            string  `json:"role"`
	MaxIal          float64 `json:"max_ial"`
	MaxAal          float64 `json:"max_aal"`
}

type NodeDetail struct {
	PublicKey       string `json:"public_key"`
	MasterPublicKey string `json:"master_public_key"`
	NodeName        string `json:"node_name"`
	Role            string `json:"role"`
	Active          bool   `json:"active"`
}

type MaxIalAal struct {
	MaxIal float64 `json:"max_ial"`
	MaxAal float64 `json:"max_aal"`
}

type MsqAddress struct {
	IP   string `json:"ip"`
	Port int64  `json:"port"`
}

// type GetNodeTokenResult struct {
// 	Amount float64 `json:"amount"`
// }

// type GetPriceFuncResult struct {
// 	Price float64 `json:"price"`
// }

type RequestIDParam struct {
	RequestID string `json:"request_id"`
}

type Namespace struct {
	Namespace   string `json:"namespace"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type UpdateNodeParam struct {
	PublicKey       string `json:"public_key"`
	MasterPublicKey string `json:"master_public_key"`
}

type Accessor struct {
	AccessorType      string `json:"accessor_type"`
	AccessorPublicKey string `json:"accessor_public_key"`
	AccessorGroupID   string `json:"accessor_group_id"`
	Active            bool   `json:"active"`
	Owner             string `json:"owner"`
}

type CheckExistingIdentityResult struct {
	Exist bool `json:"exist"`
}

type GetAccessorGroupIDResult struct {
	AccessorGroupID string `json:"accessor_group_id"`
}

type GetAccessorKeyResult struct {
	AccessorPublicKey string `json:"accessor_public_key"`
	Active            bool   `json:"active"`
}

type ServiceDetail struct {
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	Active      bool   `json:"active"`
}

type CheckExistingResult struct {
	Exist bool `json:"exist"`
}

type GetNodeInfoResult struct {
	PublicKey       string       `json:"public_key"`
	MasterPublicKey string       `json:"master_public_key"`
	NodeName        string       `json:"node_name"`
	Role            string       `json:"role"`
	Mq              []MsqAddress `json:"mq"`
}

type GetNodeInfoIdPResult struct {
	PublicKey       string       `json:"public_key"`
	MasterPublicKey string       `json:"master_public_key"`
	NodeName        string       `json:"node_name"`
	Role            string       `json:"role"`
	MaxIal          float64      `json:"max_ial"`
	MaxAal          float64      `json:"max_aal"`
	Mq              []MsqAddress `json:"mq"`
}

type GetIdentityInfoResult struct {
	Ial float64 `json:"ial"`
}

type CloseRequestParam struct {
	RequestID         string          `json:"request_id"`
	ResponseValidList []ResponseValid `json:"response_valid_list"`
}

type TimeOutRequestParam struct {
	RequestID         string          `json:"request_id"`
	ResponseValidList []ResponseValid `json:"response_valid_list"`
}

type ResponseValid struct {
	IdpID          string `json:"idp_id"`
	ValidProof     *bool  `json:"valid_proof"`
	ValidIal       *bool  `json:"valid_ial"`
	ValidSignature *bool  `json:"valid_signature"`
}

type GetDataSignatureResult struct {
	Signature string `json:"signature"`
}

type GetIdentityProofResult struct {
	IdentityProof string `json:"identity_proof"`
}

type DisableMsqDestinationParam struct {
	HashID string `json:"hash_id"`
}

type DisableAccessorMethodParam struct {
	AccessorID string `json:"accessor_id"`
}

type Service struct {
	ServiceID string  `json:"service_id"`
	MinIal    float64 `json:"min_ial"`
	MinAal    float64 `json:"min_aal"`
	Active    bool    `json:"active"`
	Suspended bool    `json:"suspended"`
}

type GetServicesByAsIDResult struct {
	Services []Service `json:"services"`
}

type ApproveService struct {
	Active bool `json:"active"`
}

type GetIdpNodesInfoResult struct {
	Node []interface{} `json:"node"`
}

type IdpNode struct {
	NodeID    string       `json:"node_id"`
	Name      string       `json:"name"`
	MaxIal    float64      `json:"max_ial"`
	MaxAal    float64      `json:"max_aal"`
	PublicKey string       `json:"public_key"`
	Mq        []MsqAddress `json:"mq"`
}

type ASWithMqNode struct {
	ID        string       `json:"node_id"`
	Name      string       `json:"name"`
	MinIal    float64      `json:"min_ial"`
	MinAal    float64      `json:"min_aal"`
	PublicKey string       `json:"public_key"`
	Mq        []MsqAddress `json:"mq"`
}

type GetAsNodesInfoByServiceIdResult struct {
	Node []interface{} `json:"node"`
}

type GetNodeInfoResultRPandASBehindProxy struct {
	PublicKey       string `json:"public_key"`
	MasterPublicKey string `json:"master_public_key"`
	NodeName        string `json:"node_name"`
	Role            string `json:"role"`
	Proxy           struct {
		NodeID          string       `json:"node_id"`
		NodeName        string       `json:"node_name"`
		PublicKey       string       `json:"public_key"`
		MasterPublicKey string       `json:"master_public_key"`
		Mq              []MsqAddress `json:"mq"`
		Config          string       `json:"config"`
	} `json:"proxy"`
}

type GetNodeInfoResultIdPandASBehindProxy struct {
	PublicKey       string  `json:"public_key"`
	MasterPublicKey string  `json:"master_public_key"`
	NodeName        string  `json:"node_name"`
	Role            string  `json:"role"`
	MaxIal          float64 `json:"max_ial"`
	MaxAal          float64 `json:"max_aal"`
	Proxy           struct {
		NodeID          string       `json:"node_id"`
		NodeName        string       `json:"node_name"`
		PublicKey       string       `json:"public_key"`
		MasterPublicKey string       `json:"master_public_key"`
		Mq              []MsqAddress `json:"mq"`
		Config          string       `json:"config"`
	} `json:"proxy"`
}

type IdpNodeBehindProxy struct {
	NodeID    string  `json:"node_id"`
	Name      string  `json:"name"`
	MaxIal    float64 `json:"max_ial"`
	MaxAal    float64 `json:"max_aal"`
	PublicKey string  `json:"public_key"`
	Proxy     struct {
		NodeID    string       `json:"node_id"`
		PublicKey string       `json:"public_key"`
		Mq        []MsqAddress `json:"mq"`
		Config    string       `json:"config"`
	} `json:"proxy"`
}

type ASWithMqNodeBehindProxy struct {
	NodeID    string  `json:"node_id"`
	Name      string  `json:"name"`
	MinIal    float64 `json:"min_ial"`
	MinAal    float64 `json:"min_aal"`
	PublicKey string  `json:"public_key"`
	Proxy     struct {
		NodeID    string       `json:"node_id"`
		PublicKey string       `json:"public_key"`
		Mq        []MsqAddress `json:"mq"`
		Config    string       `json:"config"`
	} `json:"proxy"`
}

type GetNodesBehindProxyNodeResult struct {
	Nodes []interface{} `json:"nodes"`
}

type IdPBehindProxy struct {
	NodeID          string  `json:"node_id"`
	NodeName        string  `json:"node_name"`
	Role            string  `json:"role"`
	PublicKey       string  `json:"public_key"`
	MasterPublicKey string  `json:"master_public_key"`
	MaxIal          float64 `json:"max_ial"`
	MaxAal          float64 `json:"max_aal"`
	Config          string  `json:"config"`
}

type ASorRPBehindProxy struct {
	NodeID          string `json:"node_id"`
	NodeName        string `json:"node_name"`
	Role            string `json:"role"`
	PublicKey       string `json:"public_key"`
	MasterPublicKey string `json:"master_public_key"`
	Config          string `json:"config"`
}

type Proxy struct {
	ProxyNodeID string `json:"proxy_node_id"`
	Config      string `json:"config"`
}

// type Report struct {
// 	Method string  `json:"method"`
// 	Price  float64 `json:"price"`
// 	Data   string  `json:"data"`
// }

// type GetUsedTokenReportResult []Report

type GetNodeIDListResult struct {
	NodeIDList []string `json:"node_id_list"`
}

type GetMsqAddressResult []MsqAddress
