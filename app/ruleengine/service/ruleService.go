package service

import (
	"log"
	"main/app/ruleengine/repository"
	"main/app/ruleengine/transformer"
	"main/app/utils/zkerrors"
)

type RuleService interface {
	GetAllRules(clusterId string, version int64, deleted bool, limit, offset int) (*transformer.RulesResponse, *zkerrors.ZkError)
}

type ruleService struct {
	repo repository.RulesRepo
}

func NewRuleService(repo repository.RulesRepo) RuleService {
	return &ruleService{repo: repo}
}

func (r ruleService) GetAllRules(clusterId string, version int64, deleted bool, limit, offset int) (*transformer.RulesResponse, *zkerrors.ZkError) {
	filter := repository.RuleQueryFilter{
		ClusterId: clusterId,
		Deleted:   deleted,
		Version:   version,
		Limit:     limit,
		Offset:    offset,
	}

	retVal, zkErr := r.repo.GetAllRules(&filter)
	if zkErr != nil {
		log.Println(zkErr)
		return nil, zkErr
	}

	return transformer.FromFilterRuleArrayToRulesResponse(*retVal), nil
}
