package handlerimplementation

import (
	"context"
	"encoding/json"
	"io"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
	"strconv"
)

type StringListMux struct {
	Table TablePrinterStringList
}

type TablePrinterStringList struct {
	Values []string
}

func (t *TablePrinterStringList) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterStringList) HandleRecord(ctx context.Context, r *types.Record) error {
	tempObj := make(map[string]any)
	for k, d := range r.Data {
		t.Values = append(t.Values, d.String())
		var colName string = r.TableMetadata.ColInfo[k].Name
		var value = d.String()

		var bufferSingleMap map[string]interface{}
		json.Unmarshal([]byte(value), &bufferSingleMap)
		if bufferSingleMap != nil {
			tempObj[colName] = bufferSingleMap
		} else {
			if num, err := strconv.Atoi(value); err == nil {
				tempObj[colName] = num
			} else {
				tempObj[colName] = value
			}
		}
	}
	return nil
}

func (t *TablePrinterStringList) HandleDone(ctx context.Context) error {
	return nil
}

func (s *StringListMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *StringListMux) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {
	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
