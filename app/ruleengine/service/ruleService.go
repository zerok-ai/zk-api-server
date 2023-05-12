package service

import (
	"encoding/json"
	"log"
	"main/app/ruleengine/model"
	"main/app/ruleengine/repository"
	"main/app/ruleengine/transformer"
	"main/app/utils"
	"main/app/utils/zkerrors"
)

type RuleService interface {
	GetAllRules() (*transformer.RulesResponse, *zkerrors.ZkError)
}

type ruleService struct {
	repo repository.RulesRepo
}

func NewRuleService(repo repository.RulesRepo) RuleService {
	return &ruleService{repo: repo}
}
func (r ruleService) GetAllRules() (*transformer.RulesResponse, *zkerrors.ZkError) {
	filterStringArr, err := r.repo.GetAllRules()
	var retVal []model.FilterRule
	for _, v := range filterStringArr {
		js, _ := json.Marshal(v)
		var d model.FilterRule
		err := json.Unmarshal(js, &d)
		if err != nil || d.Rules == nil {
			log.Println(err)
			return nil, utils.ToPtr[zkerrors.ZkError](zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil))
		}

		retVal = append(retVal, d)
	}

	if retVal == nil && err != nil {
		return nil, utils.ToPtr[zkerrors.ZkError](zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil))
	}

	return transformer.FromFilterRuleArrayToRulesResponse(retVal), nil
}
