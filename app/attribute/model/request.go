package model

type AttributeInfoRequest struct {
	Version       string  `csv:"version" json:"version"`
	AttributeId   string  `csv:"attr_id" json:"attr_id"`
	JsonField     *bool   `csv:"JSON" json:"json,omitempty"`
	AttributePath string  `csv:"attr_path" json:"attr_path,omitempty"`
	Field         *string `csv:"field" json:"field,omitempty"`
	DataType      *string `csv:"data_type" json:"data_type,omitempty"`
	Input         *string `csv:"input" json:"input,omitempty"`
	Values        *string `csv:"values" json:"values,omitempty"`
	Protocol      string  `csv:"protocol" json:"protocol"`
	Examples      *string `csv:"examples" json:"examples,omitempty"`
	KeySetName    *string `csv:"key_set_name" json:"key_set_name,omitempty"`
	Description   *string `csv:"description" json:"description,omitempty"`
	Executor      string  `csv:"executor" json:"executor"`
}
