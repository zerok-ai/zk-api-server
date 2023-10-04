package model

import (
	"encoding/json"
	"fmt"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"strconv"
	"strings"
)

type AttributeInfo struct {
	AttributeId      string         `json:"attr_id"`
	AttributePath    string         `json:"attr_path"`
	KeySetName       string         `json:"key_set_name"`
	JsonField        bool           `json:"json_field"`
	Field            string         `json:"field"`
	Input            string         `json:"input"`
	Values           string         `json:"values"`
	DataType         string         `json:"data_type"`
	Description      string         `json:"description,omitempty"`
	Examples         string         `json:"examples,omitempty"`
	Executor         model.Executor `json:"executor,omitempty"`
	Protocol         model.Protocol `json:"protocol,omitempty"`
	SendToFrontEnd   bool           `json:"send_to_front_end,omitempty"`
	SupportedFormats *[]string      `json:"supported_formats,omitempty"`
}

type AttributeInfoResp struct {
	Id               string    `json:"id"`
	Field            string    `json:"field"`
	Input            string    `json:"input"`
	Values           string    `json:"values,omitempty"`
	DataType         string    `json:"data_type"`
	Description      string    `json:"description,omitempty"`
	Examples         string    `json:"examples,omitempty"`
	SupportedFormats *[]string `json:"supported_formats,omitempty"`
}

type AttributeDetails struct {
	KeySetName     string              `json:"key_set_name"`
	Executor       model.Executor      `json:"executor"`
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

func getResp(attributesList []AttributeDto) []AttributeDetails {
	keySetExecutorToAttributesListStringMap := make(map[string][]AttributeInfo)
	separator := "<-.,e>"

	for _, v := range attributesList {
		attributesList := make([]AttributeInfo, 0)
		_ = json.Unmarshal([]byte(v.Attributes), &attributesList)

		for _, x := range attributesList {
			key := strings.Join([]string{string(v.Executor), x.KeySetName}, separator)
			if _, ok := keySetExecutorToAttributesListStringMap[key]; !ok {
				keySetExecutorToAttributesListStringMap[key] = make([]AttributeInfo, 0)
			}
			keySetExecutorToAttributesListStringMap[key] = append(keySetExecutorToAttributesListStringMap[key], x)
		}
	}

	attributesInfoList := make([]AttributeDetails, 0)
	for key, attributesList := range keySetExecutorToAttributesListStringMap {
		splitsArr := strings.Split(key, separator)
		executor := splitsArr[0]
		keySetName := splitsArr[1]
		attributesListForFrontend := make([]AttributeInfoResp, 0)
		for _, attribute := range attributesList {
			if attribute.SupportedFormats == nil || len(*attribute.SupportedFormats) == 0 {
				attribute.SupportedFormats = nil
			}
			a := AttributeInfoResp{
				Id:               attribute.AttributeId,
				Field:            attribute.Field,
				Input:            attribute.Input,
				Values:           attribute.Values,
				DataType:         attribute.DataType,
				Description:      attribute.Description,
				Examples:         attribute.Examples,
				SupportedFormats: attribute.SupportedFormats,
			}
			attributesListForFrontend = append(attributesListForFrontend, a)
		}
		var attributeDetails AttributeDetails
		attributeDetails.KeySetName = keySetName
		attributeDetails.Executor = model.Executor(executor)
		attributeDetails.AttributesList = attributesListForFrontend
		attributesInfoList = append(attributesInfoList, attributeDetails)
	}

	return attributesInfoList
}

func ConvertAttributeDtoToAttributeResponse(data []AttributeDto) AttributeListResponse {
	var resp AttributeListResponse
	type KeySetToAttributesList struct {
		KeySet         string
		AttributesList AttributeInfo
	}

	mapProtocolToAttributeDtoList := make(map[string][]AttributeDto)
	for _, v := range data {
		key := string(v.Protocol)
		if _, ok := mapProtocolToAttributeDtoList[key]; !ok {
			mapProtocolToAttributeDtoList[key] = make([]AttributeDto, 0)
		}
		mapProtocolToAttributeDtoList[key] = append(mapProtocolToAttributeDtoList[key], v)
	}

	for protocol, attributeDtoList := range mapProtocolToAttributeDtoList {
		attributesInfoList := getResp(attributeDtoList)
		var attributeResponse AttributeResponse
		attributeResponse.Protocol = protocol
		attributeResponse.AttributeDetailsList = attributesInfoList
		resp.AttributesList = append(resp.AttributesList, attributeResponse)
	}

	return resp
}

type ExecutorAttributesResponse struct {
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
		key := string(v.Protocol)
		if _, ok := protocolToKeySetToAttributeListMap[key]; !ok {
			protocolToKeySetToAttributeListMap[key] = make([]AttributeDto, 0)
		}
		protocolToKeySetToAttributeListMap[key] = append(protocolToKeySetToAttributeListMap[key], v)
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
				pathParts := strings.Split(attribute.AttributePath, ">")
				for i, part := range pathParts {
					pathParts[i] = strings.TrimSpace(part)
					pathParts[i] = fmt.Sprintf("\"%s\"", pathParts[i])
				}
				key := strings.Join([]string{protocol, string(attribute.Executor), v.Version}, separator)
				attributesForExecutor[attribute.AttributeId] = strings.Join(pathParts, ".")
				val := finalMap[key]
				if val == nil {
					val = make(map[string]string)
				}
				val[attribute.AttributeId] = strings.Join(pathParts, ".")
				finalMap[key] = val
			}
		}
	}

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
	resp.Version = maxUpdatedAt
	resp.Attributes = executorList

	return resp
}
