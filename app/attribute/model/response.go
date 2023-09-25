package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type AttributeInfo struct {
	Id               string `json:"id"`
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
					idParts := strings.Split(attribute.Id, ">")
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
