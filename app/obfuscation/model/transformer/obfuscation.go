package transformer

import (
	"encoding/json"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkObfuscation "github.com/zerok-ai/zk-utils-go/obfuscation/model"
	"time"
	dto "zk-api-server/app/obfuscation/model/dto"
)

var LogTag = "obfuscation_transformer"

type ObfuscationListResponse struct {
	Response []zkObfuscation.Rule `json:"obfuscations"`
}

type ObfuscationResponse struct {
	Response zkObfuscation.Rule `json:"obfuscation"`
}

func ToObfuscationListResponse(oArr []dto.Obfuscation) ObfuscationListResponse {
	var obfuscations []zkObfuscation.Rule
	for _, o := range oArr {
		var rule zkObfuscation.Rule
		err := json.Unmarshal(o.RuleDef, &rule)
		if err != nil {
			zkLogger.Error(LogTag, "Error while unmarshalling the obfuscation rule: ", err)
			continue
		}
		obfuscations = append(obfuscations, rule)
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
	return ObfuscationResponse{Response: rule}, nil
}

func FromObfuscationRequestToObfuscationDto(oReq zkObfuscation.Rule, orgId string, id string) *dto.Obfuscation {
	currentTime := time.Now()
	ruleDef, err := json.Marshal(oReq)
	if err != nil {
		zkLogger.Error(LogTag, "Error while marshalling the obfuscation rule: ", err)
		return nil
	}
	return &dto.Obfuscation{
		ID:        &id,
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
