package models

import (
	"main/app/utils/zkerrors"
	"px.dev/pxapi"
)

type PixieResponse struct {
	Result       interface{}
	ResultsStats *pxapi.ResultsStats
	Error        *zkerrors.ZkError
}
