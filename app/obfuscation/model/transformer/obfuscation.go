package transformer

import (
	"encoding/json"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkObfuscation "github.com/zerok-ai/zk-utils-go/obfuscation/model"
	"time"
	dto "zk-api-server/app/obfuscation/model/dto"
)

var LogTag = "obfuscation_transformer"

type ObfuscationResponseOperator struct {
	Obfuscations []zkObfuscation.RuleOperator `json:"obfuscations"`
	DeletedIds   []string                     `json:"deleted_obfuscations"`
	DisabledIds  []string                     `json:"disabled_obfuscations"`
}

type ObfuscationListResponse struct {
	Response []zkObfuscation.Rule `json:"obfuscations"`
}

type ObfuscationResponse struct {
	Response *zkObfuscation.Rule `json:"obfuscation"`
}

func ToObfuscationResponseOperator(obj dto.Obfuscation) (zkObfuscation.RuleOperator, error) {
	var rule zkObfuscation.RuleOperator
	err := json.Unmarshal(obj.RuleDef, &rule)
	if err != nil {
		zkLogger.Error(LogTag, "Error while unmarshalling the obfuscation rule: ", err)
		return zkObfuscation.RuleOperator{}, err
	}
	rule.Id = obj.ID
	rule.UpdatedAt = obj.UpdatedAt
	return rule, nil
}

func ToObfuscationListResponseOperator(oArr []dto.Obfuscation) ObfuscationResponseOperator {
	active := []zkObfuscation.RuleOperator{}
	deleted := []string{}
	disabled := []string{}
	for _, o := range oArr {
		rule, err := ToObfuscationResponseOperator(o)
		if err != nil {
			zkLogger.Error(LogTag, "Error while converting dto.obfuscation the obfuscation response model: ", err)
			continue
		}
		if o.Deleted {
			deleted = append(deleted, rule.Id)
		} else if o.Disabled {
			disabled = append(disabled, rule.Id)
		} else {
			active = append(active, rule)
		}
	}
	return ObfuscationResponseOperator{Obfuscations: active, DeletedIds: deleted, DisabledIds: disabled}
}

func ToObfuscationListResponse(oArr []dto.Obfuscation) ObfuscationListResponse {
	obfuscations := []zkObfuscation.Rule{}
	for _, o := range oArr {
		rule, err := ToObfuscationResponse(o)
		if err != nil {
			zkLogger.Error(LogTag, "Error while converting dto.obfuscation the obfuscation response model: ", err)
			continue
		}
		obfuscations = append(obfuscations, *rule.Response)
	}
	return ObfuscationListResponse{Response: obfuscations}
}

func ToObfuscationResponse(obj dto.Obfuscation) (ObfuscationResponse, error) {
	var rule zkObfuscation.Rule
	err := json.Unmarshal(obj.RuleDef, &rule)
	if err != nil {
		zkLogger.Error(LogTag, "Error while unmarshalling the obfuscation rule: ", err)
		return ObfuscationResponse{}, err
	}
	rule.Id = obj.ID
	rule.CreatedAt = obj.CreatedAt
	rule.UpdatedAt = obj.UpdatedAt
	return ObfuscationResponse{Response: &rule}, nil
}

func FromObfuscationRequestToObfuscationDto(oReq zkObfuscation.Rule, orgId string, id string) *dto.Obfuscation {
	currentTime := time.Now().Unix()
	ruleDef, err := json.Marshal(oReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while marshalling the obfuscation rule: ", err)
		return nil
	}
	return &dto.Obfuscation{
		ID:        id,
		OrgID:     orgId,
		RuleName:  oReq.Name,
		RuleType:  "obfuscation_rule",
		RuleDef:   ruleDef,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Deleted:   false,
		Disabled:  !oReq.Enabled,
	}
}
