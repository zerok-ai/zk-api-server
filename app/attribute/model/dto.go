package model

import (
	"encoding/json"
	"strings"
)

type AttributeDto struct {
	KeySet     string `json:"key_set"`
	Version    string `json:"version"`
	Attributes string `json:"attribute_list"`
}

type AttributeDtoList []AttributeDto

func ConvertAttributeInfoRequestToAttributeDto(req []AttributeInfoRequest) AttributeDtoList {
	myMap := make(map[string]map[string][]AttributeInfo)
	attributeDtoList := make(AttributeDtoList, 0)
	for _, v := range req {
		element := AttributeInfo{
			Attribute:        v.Attribute,
			Type:             v.Type,
			Description:      v.Description,
			Examples:         v.Examples,
			RequirementLevel: v.RequirementLevel,
		}

		keySetName := strings.Trim(v.KeySetName, " ")
		version := strings.Trim(v.Version, " ")
		if _, ok := myMap[keySetName]; ok {
			if _, ok := myMap[keySetName][version]; ok {
				myMap[keySetName][version] = append(myMap[keySetName][version], element)
			} else {
				myMap[keySetName][version] = []AttributeInfo{element}
			}
		} else {
			myMap[keySetName] = map[string][]AttributeInfo{version: {element}}
		}
	}

	for k, value := range myMap {
		for version, v := range value {
			attrStr, _ := json.Marshal(v)
			attributeDto := AttributeDto{
				KeySet:     k,
				Attributes: string(attrStr),
				Version:    version,
			}
			attributeDtoList = append(attributeDtoList, attributeDto)
		}
	}

	return attributeDtoList
}

func (t AttributeDto) GetAllColumns() []any {
	return []any{t.Version, t.KeySet, t.Attributes}
}
