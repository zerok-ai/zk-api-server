package model

import "encoding/json"

type AttributeInfo struct {
	Attribute        string `json:"attribute"`
	Type             string `json:"type"`
	Description      string `json:"description,omitempty"`
	Examples         string `json:"examples,omitempty"`
	RequirementLevel string `json:"requirementLevel,omitempty"`
}

type AttributeResponse struct {
	KeySet     string          `json:"key_set,omitempty"`
	Version    string          `json:"version,omitempty"`
	Attributes []AttributeInfo `json:"attribute_list"`
}

type AttributeListResponse struct {
	AttributesList []AttributeResponse `json:"attributesList"`
}

func ConvertAttributeDtoToAttributeResponse(data []AttributeDto) AttributeListResponse {
	var resp AttributeListResponse
	for _, v := range data {
		var attr AttributeResponse
		attr.KeySet = v.KeySet
		attr.Version = v.Version
		json.Unmarshal([]byte(v.Attributes), &attr.Attributes)
		resp.AttributesList = append(resp.AttributesList, attr)
	}
	return resp
}
