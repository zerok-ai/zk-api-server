package zkpostgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	zkLogger "main/app/utils/logs"
	pgConfig "main/app/utils/postgres/config"
	"main/app/utils/zkerrors"
)

type DatabaseRepo[T any] interface {
	CreateConnection() *sql.DB
	Delete(stmt string, param []any, tx *sql.Tx, rollback bool) (int, error)
	Get(query string, param []any, args ...any) *zkerrors.ZkError
	GetAll(query string, param []any, rowsProcessor func(rows *sql.Rows, sqlErr error) (*[]T, *[]string, *zkerrors.ZkError)) (*[]T, *[]string, *zkerrors.ZkError)
}

type zkPostgresRepo[T any] struct {
}

func NewZkPostgresRepo[T any]() DatabaseRepo[T] {
	x := zkPostgresRepo[T]{}
	return &x
}

var config pgConfig.PostgresConfig

var LOG_TAG = "zkpostgres_db_repo"

func Init(c pgConfig.PostgresConfig) {
	config = c
}

func (zkPostgresService zkPostgresRepo[T]) CreateConnection() *sql.DB {

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

func (zkPostgresService zkPostgresRepo[T]) Get(query string, param []any, args ...any) *zkerrors.ZkError {
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

func (zkPostgresService zkPostgresRepo[T]) GetAll(query string, param []any, rowsProcessor func(rows *sql.Rows, sqlErr error) (*[]T, *[]string, *zkerrors.ZkError)) (*[]T, *[]string, *zkerrors.ZkError) {
	db := zkPostgresService.CreateConnection()
	defer db.Close()
	rows, err := db.Query(query, param...)
	defer rows.Close()
	return rowsProcessor(rows, err)
}

func (zkPostgresService zkPostgresRepo[T]) Delete(stmt string, param []any, tx *sql.Tx, rollback bool) (int, error) {
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
