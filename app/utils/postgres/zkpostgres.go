package zkpostgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	zkLogger "main/app/utils/logs"
	pgConfig "main/app/utils/postgres/config"
	"main/app/utils/zkerrors"
)

type DatabaseRepo interface {
	CreateConnection() *sql.DB
	Delete(stmt string, param []any, tx *sql.Tx, rollback bool) (int, error)
	Get(query string, param []any, args ...any) *zkerrors.ZkError
	GetAll(query string, param []any) (*sql.Rows, error)
}

type zkPostgresRepo struct {
}

func NewZkPostgresRepo() DatabaseRepo {
	return &zkPostgresRepo{}
}

var config pgConfig.PostgresConfig

var LOG_TAG = "zkpostgres_db_repo"

func Init(c pgConfig.PostgresConfig) {
	config = c
}

func (zkPostgresService zkPostgresRepo) CreateConnection() *sql.DB {
	config.Host = "localhost"
	config.Port = 5432
	config.Password = "pl"
	config.Dbname = "pl"
	config.User = "pl"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)
	//zkLogger.Debug(LOG_TAG, "psqlInfo==", psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func (zkPostgresService zkPostgresRepo) Get(query string, param []any, args ...any) *zkerrors.ZkError {
	db := zkPostgresService.CreateConnection()
	defer db.Close()
	row := db.QueryRow(query, param...)
	err := row.Scan(args...)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_NOT_FOUND, nil)
		return &zkError
	case nil:
		return nil
	default:
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZK_ERROR_INTERNAL_SERVER, nil)
		zkLogger.Debug(LOG_TAG, "unable to scan rows", err)
		return &zkError
	}
}

func (zkPostgresService zkPostgresRepo) GetAll(query string, param []any) (*sql.Rows, error) {
	db := zkPostgresService.CreateConnection()
	defer db.Close()
	return db.Query(query, param...)
}

func (zkPostgresService zkPostgresRepo) Delete(stmt string, param []any, tx *sql.Tx, rollback bool) (int, error) {
	res, err := tx.Exec(stmt, param...)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil {
			return int(count), nil
		} else {
			return 0, err
		}
	}

	zkLogger.Debug(LOG_TAG, err.Error())
	return 0, err
}
