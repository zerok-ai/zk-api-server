package model

import (
	"encoding/json"
	"strings"
)

type AttributeDto struct {
	Version    string `json:"version"`
	KeySet     string `json:"key_set"`
	Protocol   string `json:"protocol"`
	Executor   string `json:"executor"`
	UpdatedAt  string `json:"updated_at"`
	Attributes string `json:"attribute_list"`
}

type AttributeDtoList []AttributeDto

func ConvertAttributeInfoRequestToAttributeDto(req []AttributeInfoRequest) AttributeDtoList {
	keySetNameToAttributesInfoMap := make(map[string]map[string][]AttributeInfo)
	attributeDtoList := make(AttributeDtoList, 0)
	for _, v := range req {
		element := AttributeInfo{
			CommonId:         v.CommonId,
			VersionId:        v.VersionId,
			Field:            v.Field,
			Input:            v.Input,
			Values:           v.Values,
			DataType:         v.DataType,
			Description:      v.Description,
			Examples:         v.Examples,
			RequirementLevel: v.RequirementLevel,
			Protocol:         v.Protocol,
			Executor:         v.Executor,
			SendToFrontEnd:   v.SendToFrontEnd,
		}

		keySetName := strings.Trim(v.KeySetName, " ")
		version := strings.Trim(v.Version, " ")
		if _, ok := keySetNameToAttributesInfoMap[keySetName]; ok {
			if _, ok := keySetNameToAttributesInfoMap[keySetName][version]; ok {
				keySetNameToAttributesInfoMap[keySetName][version] = append(keySetNameToAttributesInfoMap[keySetName][version], element)
			} else {
				keySetNameToAttributesInfoMap[keySetName][version] = []AttributeInfo{element}
			}
		} else {
			keySetNameToAttributesInfoMap[keySetName] = map[string][]AttributeInfo{version: {element}}
		}
	}

	for k, value := range keySetNameToAttributesInfoMap {
		for version, v := range value {
			executor := getExecutor(v)
			protocol := getProtocol(v)
			if executor == "" || protocol == "" {
				return nil
			}
			attrStr, _ := json.Marshal(v)
			attributeDto := AttributeDto{
				Protocol:   protocol,
				Executor:   executor,
				KeySet:     k,
				Version:    version,
				Attributes: string(attrStr),
			}
			attributeDtoList = append(attributeDtoList, attributeDto)
		}
	}

	return attributeDtoList
}

func getProtocol(attributes []AttributeInfo) string {
	if len(attributes) == 0 {
		return ""
	}

	protocol := attributes[0].Protocol
	for _, v := range attributes {
		if v.Protocol != protocol {
			return ""
		}
	}
	return protocol
}

func getExecutor(attributes []AttributeInfo) string {
	if len(attributes) == 0 {
		return ""
	}

	executor := attributes[0].Executor
	for _, v := range attributes {
		if v.Executor != executor {
			return ""
		}
	}
	return executor
}

func (t AttributeDto) GetAllColumns() []any {
	return []any{t.Version, t.KeySet, t.Protocol, t.Executor, t.Attributes}
}
