package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zerok-ai/zk-utils-go/rules/model"
	zkUtilsGo "github.com/zerok-ai/zk-utils-go/utils"
	"log"
	"main/app/utils"
	zkLogger "main/app/utils/logs"
	zkPostgres "main/app/utils/postgres"
	"main/app/utils/zkerrors"
	"sort"
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
	GetAllRules(filters *RuleQueryFilter) (*[]model.FilterRule, *zkerrors.ZkError)
}

//
//type rulesFromFileRepo struct {
//}
//
//var filePath = "data.json"
//
//func NewRulesFromFileRepo() RulesRepo {
//	return &rulesFromFileRepo{}
//}
//
//func (r rulesFromFileRepo) GetAllRules(filters *RuleQueryFilter) (*[]model.NewRuleSchema, *zkerrors.ZkError) {
//	var err error
//	file, err := os.Open(filePath)
//	if err != nil {
//		log.Printf("unable to open file, err: %s", err.Error())
//		return nil, utils.ToPtr[zkerrors.ZkError](zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil))
//	}
//	defer file.Close()
//
//	scanner := bufio.NewScanner(file)
//	var filterStringArr []map[string]interface{}
//	for scanner.Scan() {
//		var data map[string]interface{}
//		b := scanner.Bytes()
//		err = json.Unmarshal(b, &data)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//
//		filterStringArr = append(filterStringArr, data)
//	}
//
//	var retVal []model.NewRuleSchema
//	for _, v := range filterStringArr {
//		js, _ := json.Marshal(v)
//		var d model.NewRuleSchema
//		err := json.Unmarshal(js, &d)
//		if err != nil || d.Workloads == nil {
//			log.Println(err)
//			return nil, utils.ToPtr[zkerrors.ZkError](zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil))
//		}
//
//		retVal = append(retVal, d)
//	}
//	return &retVal, nil
//}

type zkPostgresRepo struct {
}

func NewZkPostgresRepo() RulesRepo {
	return &zkPostgresRepo{}
}

func (zkPostgresService zkPostgresRepo) GetAllRules(filters *RuleQueryFilter) (*[]model.FilterRule, *zkerrors.ZkError) {
	query := GetAllRulesSqlStatement
	zkPostgresRepo := zkPostgres.NewZkPostgresRepo[model.FilterRule]()

	params := []any{filters.ClusterId, filters.Version, filters.Deleted, filters.Version, filters.Limit, filters.Offset}
	return zkPostgresRepo.GetAll(query, params, Processor)
}

func Processor(rows *sql.Rows, sqlErr error) (*[]model.FilterRule, *zkerrors.ZkError) {
	defer rows.Close()

	switch sqlErr {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_NOT_FOUND, nil)
		return nil, &zkError
	case nil:
		break
	default:
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
		zkLogger.Debug(LOG_TAG, "unable to scan rows", zkError)
		return nil, &zkError
	}

	if rows == nil {
		zkErr := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
		return nil, &zkErr
	}

	var ruleString string
	var ruleStringArr []string
	var rulesList []model.FilterRule

	for rows.Next() {

		// Scan the values from the current row into variables
		err := rows.Scan(&ruleString)
		if err != nil {
			log.Fatal(err)
		}

		// Print the retrieved values
		fmt.Printf("Filter: %s\n", ruleString)
		ruleStringArr = append(ruleStringArr, ruleString)
	}

	// Check for any errors occurred during iteration
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for _, js := range ruleStringArr {
		var d model.FilterRule
		err := json.Unmarshal([]byte(js), &d)
		if err != nil || d.Workloads == nil {
			log.Println(err)
			return nil, utils.ToPtr[zkerrors.ZkError](zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil))
		}

		rulesList = append(rulesList, d)
		for _, v := range d.Workloads {
			WorkLoadUUID(v)
		}
	}

	return &rulesList, nil
}

func WorkLoadUUID(w model.WorkloadRule) uuid.UUID {
	sort.Sort(model.Rules(w.ConditionalRule.RuleSet))
	jStr, _ := json.Marshal(w)
	id := zkUtilsGo.CalculateHash(string(jStr))
	fmt.Println(id)
	return id
}

const GetAllRulesSqlStatement = `SELECT filters FROM FilterRule WHERE is_default = TRUE UNION SELECT filters FROM FilterRule WHERE cluster_id=$1 AND version>$2 AND (deleted=$3 OR deleted_at>$4) LIMIT $5 OFFSET $6`
