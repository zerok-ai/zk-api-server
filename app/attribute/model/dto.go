package model

import (
	"encoding/json"
	"strings"
)

type AttributeDto struct {
	Protocol   string `json:"protocol"`
	KeySet     string `json:"key_set"`
	Version    string `json:"version"`
	Attributes string `json:"attribute_list"`
}

type AttributeDtoList []AttributeDto

func ConvertAttributeInfoRequestToAttributeDto(req []AttributeInfoRequest) AttributeDtoList {
	keySetNameToAttributesInfoMap := make(map[string]map[string][]AttributeInfo)
	attributeDtoList := make(AttributeDtoList, 0)
	for _, v := range req {
		element := AttributeInfo{
			Id:               v.Id,
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
			attrStr, _ := json.Marshal(v)
			attributeDto := AttributeDto{
				KeySet:     k,
				Version:    version,
				Attributes: string(attrStr),
			}
			attributeDtoList = append(attributeDtoList, attributeDto)
		}
	}

	return attributeDtoList
}

func (t AttributeDto) GetAllColumns() []any {
	return []any{t.Version, t.KeySet, t.Attributes}
}
