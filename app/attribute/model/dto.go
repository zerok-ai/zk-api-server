package model

import (
	"encoding/json"
	"github.com/zerok-ai/zk-utils-go/scenario/model"
	"strings"
)

type AttributeDto struct {
	Version    string         `json:"version"`
	Protocol   model.Protocol `json:"protocol"`
	Executor   model.Executor `json:"executor"`
	UpdatedAt  string         `json:"updated_at"`
	Attributes string         `json:"attribute_list"`
}

type AttributeDtoList []AttributeDto

func ConvertAttributeInfoRequestToAttributeDto(req []AttributeInfoRequest) AttributeDtoList {
	version := strings.Trim(req[0].Version, " ")
	executor := strings.Trim(string(req[0].Executor), " ")
	if version != "common" {
		for _, v := range req {
			v.SupportedFormats = nil
			v.Field = nil
			v.DataType = nil
			v.Input = nil
			v.Values = nil
			v.Examples = nil
			v.KeySetName = nil
			v.Description = nil
			v.KeySetName = nil
		}
	} else {
		for _, v := range req {
			v.AttributePath = ""
		}
	}

	protocolToAttributesInfoRequestListMap := getProtocolToAttributesMap(req)
	attributeDtoList := make(AttributeDtoList, 0)
	for protocol, attributesInfoRequestList := range protocolToAttributesInfoRequestListMap {
		attrStr, _ := json.Marshal(attributesInfoRequestList)
		attributeDto := AttributeDto{
			Protocol:   model.Protocol(protocol),
			Executor:   model.Executor(executor),
			Version:    version,
			Attributes: string(attrStr),
		}

		attributeDtoList = append(attributeDtoList, attributeDto)
	}

	return attributeDtoList
}

func getProtocolToAttributesMap(reqInfoList []AttributeInfoRequest) map[string][]AttributeInfoRequest {
	protocolToAttributesInfoRequestListMap := make(map[string][]AttributeInfoRequest)
	for _, v := range reqInfoList {
		protocol := strings.Trim(string(v.Protocol), " ")
		if _, ok := protocolToAttributesInfoRequestListMap[protocol]; ok {
			protocolToAttributesInfoRequestListMap[protocol] = append(protocolToAttributesInfoRequestListMap[protocol], v)
		} else {
			protocolToAttributesInfoRequestListMap[protocol] = []AttributeInfoRequest{v}
		}
	}

	return protocolToAttributesInfoRequestListMap
}

func (t AttributeDto) GetAllColumns() []any {
	return []any{t.Version, t.Protocol, t.Executor, t.Attributes}
}
