package handlerimplementation

import (
	"context"
	"io"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
)

type TablePrinterServiceMap struct {
	Values []ServiceMap
}

type ServiceMapMux struct {
	Table TablePrinterServiceMap
}

func (t *TablePrinterServiceMap) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterServiceMap) HandleRecord(ctx context.Context, r *types.Record) error {
	t.Values = append(t.Values, ConvertPixieDataToServiceMap(r))
	return nil
}

func (t *TablePrinterServiceMap) HandleDone(ctx context.Context) error {
	return nil
}

func (s *ServiceMapMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *ServiceMapMux) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {

	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
