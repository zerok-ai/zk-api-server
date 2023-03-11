package utils

import (
	"main/app/utils"
	"px.dev/pxapi/types"
	"strconv"
)

func GetString(key string, r *types.Record) string {
	return r.GetDatum(key).String()
}

func GetFloat(key string, r *types.Record, bitSize int) (float64, error) {
	return strconv.ParseFloat(GetString(key, r), bitSize)
}

func GetInteger(key string, r *types.Record) (int, error) {
	return strconv.Atoi(GetString(key, r))
}

func GetStringPtr(key string, r *types.Record) *string {
	return utils.StringToPtr(GetString(key, r))
}

func GetFloat64Ptr(key string, r *types.Record) *float64 {
	v, err := GetFloat(key, r, 64)
	if err != nil {
		return nil
	} else {
		return &v
	}
}

func GetFloat32Ptr(key string, r *types.Record) *float64 {
	v, err := GetFloat(key, r, 32)
	if err != nil {
		return nil
	} else {
		return &v
	}
}

func GetIntegerPtr(key string, r *types.Record) *int {
	v, err := strconv.Atoi(GetString(key, r))
	if err != nil {
		return nil
	} else {
		return &v
	}
}
