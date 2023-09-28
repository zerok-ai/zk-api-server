package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type AttributeInfo struct {
	CommonId         string `json:"common_id"`
	VersionId        string `json:"version_id"`
	Field            string `json:"field"`
	Input            string `json:"input"`
	Values           string `json:"values"`
	DataType         string `json:"data_type"`
	Description      string `json:"description,omitempty"`
	Examples         string `json:"examples,omitempty"`
	RequirementLevel string `json:"requirement_level,omitempty"`
	Executor         string `json:"executor,omitempty"`
	Protocol         string `json:"protocol,omitempty"`
	SendToFrontEnd   bool   `json:"send_to_front_end,omitempty"`
}

type AttributeInfoResp struct {
	Id               string `json:"id"`
	Field            string `json:"field"`
	Input            string `json:"input"`
	Values           string `json:"values,omitempty"`
	DataType         string `json:"data_type"`
	Description      string `json:"description,omitempty"`
	Examples         string `json:"examples,omitempty"`
	RequirementLevel string `json:"requirement_level,omitempty"`
}

type AttributeDetails struct {
	KeySetName     string              `json:"key_set_name"`
	Executor       string              `json:"executor"`
	Version        string              `json:"version,omitempty"`
	AttributesList []AttributeInfoResp `json:"attribute_list"`
}

type AttributeResponse struct {
	Protocol             string             `json:"protocol"`
	AttributeDetailsList []AttributeDetails `json:"attribute_details"`
}

type AttributeListResponse struct {
	AttributesList []AttributeResponse `json:"attributes_list"`
}

func ConvertAttributeDtoToAttributeResponse(data []AttributeDto) AttributeListResponse {
	var resp AttributeListResponse
	var protocolToKeySetToAttributeListMap = make(map[string][]AttributeDto)
	for _, v := range data {
		protocolToKeySetToAttributeListMap[v.Protocol] = append(protocolToKeySetToAttributeListMap[v.Protocol], v)
	}

	for protocol, keySetToAttributeListMap := range protocolToKeySetToAttributeListMap {
		attributesDetailsList := make([]AttributeDetails, 0)
		for _, v := range keySetToAttributeListMap {
			attributesList := make([]AttributeInfo, 0)
			_ = json.Unmarshal([]byte(v.Attributes), &attributesList)
			attributesListForFrontend := make([]AttributeInfoResp, 0)
			for _, attribute := range attributesList {
				if attribute.SendToFrontEnd == true {
					idParts := strings.Split(attribute.CommonId, ">")
					for i, part := range idParts {
						idParts[i] = strings.TrimSpace(part)
						idParts[i] = fmt.Sprintf("\"%s\"", idParts[i])
					}
					a := AttributeInfoResp{
						Id:               strings.Join(idParts, "."),
						Field:            attribute.Field,
						Input:            attribute.Input,
						Values:           attribute.Values,
						DataType:         attribute.DataType,
						Description:      attribute.Description,
						Examples:         attribute.Examples,
						RequirementLevel: attribute.RequirementLevel,
					}
					attributesListForFrontend = append(attributesListForFrontend, a)
				}
			}

			if len(attributesListForFrontend) == 0 {
				continue
			}

			var attributeDetails AttributeDetails
			attributeDetails.KeySetName = v.KeySet
			attributeDetails.Executor = attributesList[0].Executor
			attributeDetails.AttributesList = attributesListForFrontend
			attributesDetailsList = append(attributesDetailsList, attributeDetails)

		}

		var attributeResponse AttributeResponse
		attributeResponse.Protocol = protocol
		attributeResponse.AttributeDetailsList = attributesDetailsList
		resp.AttributesList = append(resp.AttributesList, attributeResponse)
	}

	return resp
}

type ExecutorAttributesResponse struct {
	ExecutorAttributesList ExecutorAttributesList `json:"executor_attributes_list"`
}

type ExecutorAttributesList struct {
	Attributes []ExecutorAttributes `json:"executor_attributes"`
	Version    int64                `json:"version"`
	Update     bool                 `json:"update"`
}

type ExecutorAttributes struct {
	Executor   string            `json:"executor"`
	Version    string            `json:"version"`
	Protocol   string            `json:"protocol"`
	Attributes map[string]string `json:"attributes"`
}

func ConvertAttributeDtoToExecutorAttributesResponse(data []AttributeDto) ExecutorAttributesResponse {
	var resp ExecutorAttributesResponse
	var protocolToKeySetToAttributeListMap = make(map[string][]AttributeDto)
	var maxUpdatedAt int64
	for _, v := range data {
		protocolToKeySetToAttributeListMap[v.Protocol] = append(protocolToKeySetToAttributeListMap[v.Protocol], v)
		v, _ := strconv.ParseInt(v.UpdatedAt, 10, 64)
		if v > maxUpdatedAt {
			maxUpdatedAt = v
		}
	}

	separator := "<-.,e>"

	finalMap := make(map[string]map[string]string)
	for protocol, keySetToAttributeListMap := range protocolToKeySetToAttributeListMap {
		for _, v := range keySetToAttributeListMap {
			attributesList := make([]AttributeInfo, 0)
			_ = json.Unmarshal([]byte(v.Attributes), &attributesList)
			attributesForExecutor := make(map[string]string)
			for _, attribute := range attributesList {
				if attribute.SendToFrontEnd == true {
					commonIdParts := strings.Split(attribute.CommonId, ">")
					for i, part := range commonIdParts {
						commonIdParts[i] = strings.TrimSpace(part)
						commonIdParts[i] = fmt.Sprintf("\"%s\"", commonIdParts[i])
					}

					versionIdParts := strings.Split(attribute.VersionId, ">")
					for i, part := range versionIdParts {
						versionIdParts[i] = strings.TrimSpace(part)
						versionIdParts[i] = fmt.Sprintf("\"%s\"", versionIdParts[i])
					}
					key := strings.Join([]string{protocol, attribute.Executor, v.Version}, separator)
					attributesForExecutor[strings.Join(commonIdParts, ".")] = strings.Join(versionIdParts, ".")
					val := finalMap[key]
					if val == nil {
						val = make(map[string]string)
					}
					val[strings.Join(commonIdParts, ".")] = strings.Join(versionIdParts, ".")
					finalMap[key] = val
				}
			}
		}
	}

	var executorAttributesList ExecutorAttributesList
	executorList := make([]ExecutorAttributes, 0)
	for k, v := range finalMap {
		e := ExecutorAttributes{
			Executor:   strings.Split(k, separator)[1],
			Version:    strings.Split(k, separator)[2],
			Protocol:   strings.Split(k, separator)[0],
			Attributes: v,
		}
		executorList = append(executorList, e)
	}
	executorAttributesList.Version = maxUpdatedAt
	executorAttributesList.Attributes = executorList
	resp.ExecutorAttributesList = executorAttributesList

	return resp
}
