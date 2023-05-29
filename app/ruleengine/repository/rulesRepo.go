package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/zerok-ai/zk-utils-go/rules/model"
	"log"
	rulesResponseModel "main/app/ruleengine/model"
	"main/app/utils"
	zkLogger "main/app/utils/logs"
	zkPostgres "main/app/utils/postgres"
	"main/app/utils/zkerrors"
)

var LOG_TAG = "zkpostgres_db_repo"

type RuleQueryFilter struct {
	ClusterId string
	Deleted   bool
	Version   int64
	Limit     int
	Offset    int
}

type RulesRepo interface {
	GetAllRules(filters *RuleQueryFilter) (*[]model.Scenario, *[]string, *zkerrors.ZkError)
}

type zkPostgresRepo struct {
}

func NewZkPostgresRepo() RulesRepo {
	return &zkPostgresRepo{}
}

func (zkPostgresService zkPostgresRepo) GetAllRules(filters *RuleQueryFilter) (*[]model.Scenario, *[]string, *zkerrors.ZkError) {
	query := GetAllRulesSqlStatement
	zkPostgresRepo := zkPostgres.NewZkPostgresRepo[model.Scenario]()

	params := []any{filters.ClusterId, filters.Version, filters.Limit, filters.Offset}
	return zkPostgresRepo.GetAll(query, params, Processor)
}

func Processor(rows *sql.Rows, sqlErr error) (*[]model.Scenario, *[]string, *zkerrors.ZkError) {
	defer rows.Close()

	switch sqlErr {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_NOT_FOUND, nil)
		return nil, nil, &zkError
	case nil:
		break
	default:
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
		zkLogger.Debug(LOG_TAG, "unable to scan rows", zkError)
		return nil, nil, &zkError
	}

	if rows == nil {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
		return nil, nil, &zkErr
	}

	var rulesResponse rulesResponseModel.RulesDbResponse
	var rulesResponseArr []rulesResponseModel.RulesDbResponse

	for rows.Next() {

		// Scan the values from the current row into variables
		err := rows.Scan(&rulesResponse.Filters, &rulesResponse.Deleted)
		if err != nil {
			log.Fatal(err)
		}

		// Print the retrieved values
		//fmt.Printf("Filter: %s\n", rulesResponseModel)
		rulesResponseArr = append(rulesResponseArr, rulesResponse)
	}

	// Check for any errors occurred during iteration
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	var rulesList []model.Scenario
	var deletedRulesList []string
	for _, rs := range rulesResponseArr {
		var d model.Scenario
		err := json.Unmarshal([]byte(rs.Filters), &d)
		if err != nil || d.Workloads == nil {
			log.Println(err)
			return nil, nil, utils.ToPtr[zkerrors.ZkError](zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil))
		}

		if rs.Deleted == false {
			rulesList = append(rulesList, d)
			//for oldId, v := range d.Workloads {
			//	id := model.WorkLoadUUID(v)
			//	delete(d.Workloads, oldId)
			//	d.Workloads[id.String()] = v
			//}
		} else {
			deletedRulesList = append(deletedRulesList, d.ScenarioId)
		}
	}

	return &rulesList, &deletedRulesList, nil
}

const GetAllRulesSqlStatement = `SELECT filters, deleted FROM Scenario WHERE (cluster_id=$1 OR cluster_id IS NULL) AND version>$2 LIMIT $3 OFFSET $4`
