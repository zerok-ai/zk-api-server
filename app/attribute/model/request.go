package model

type AttributeInfoRequest struct {
	Version          string `csv:"version" json:"version"`
	CommonId         string `csv:"common_id" json:"common_id"`
	VersionId        string `csv:"version_id" json:"version_id"`
	Field            string `csv:"field" json:"field"`
	Input            string `csv:"input" json:"input"`
	Values           string `csv:"values" json:"values"`
	KeySetName       string `csv:"key_set_name" json:"key_set_name"`
	DataType         string `csv:"data_type" json:"data_type"`
	Description      string `csv:"description" json:"description"`
	Examples         string `csv:"examples" json:"examples"`
	RequirementLevel string `csv:"requirement_level" json:"requirement_level"`
	Protocol         string `csv:"protocol" json:"protocol"`
	Executor         string `csv:"executor" json:"executor"`
	SendToFrontEnd   bool   `csv:"send_to_front_end" json:"send_to_front_end"`
}
