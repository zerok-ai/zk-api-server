package handlerimplementation

import (
	"context"
	"io"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
)

type ServiceListMux struct {
	Table TablePrinterServiceList
}

type TablePrinterServiceList struct {
	Values []Service
}

func (t *TablePrinterServiceList) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterServiceList) HandleRecord(ctx context.Context, r *types.Record) error {
	t.Values = append(t.Values, ConvertPixieDataToService(r))
	return nil
}

func (t *TablePrinterServiceList) HandleDone(ctx context.Context) error {
	return nil
}

func (s *ServiceListMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *ServiceListMux) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {

	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
