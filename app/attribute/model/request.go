package model

type AttributeInfoRequest struct {
	Version          string `csv:"Version" json:"version"`
	Attribute        string `csv:"Attribute" json:"attribute"`
	KeySetName       string `csv:"keyset_name" json:"keyset_name"`
	Type             string `csv:"Type" json:"type"`
	Description      string `csv:"Description" json:"description"`
	Examples         string `csv:"Examples" json:"examples"`
	RequirementLevel string `csv:"Requirement Level" json:"requirementLevel"`
}
