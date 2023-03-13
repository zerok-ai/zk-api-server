package tablemux

import (
	"context"
	"io"
	"main/app/cluster/models"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
)

type TablePrinterServiceStat struct {
	Values []models.ServiceStat
}

type ServiceStatMux struct {
	Table TablePrinterServiceStat
}

func (t *TablePrinterServiceStat) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterServiceStat) HandleRecord(ctx context.Context, r *types.Record) error {
	t.Values = append(t.Values, models.ConvertPixieDataToServiceStat(r))
	return nil
}

func (t *TablePrinterServiceStat) HandleDone(ctx context.Context) error {
	return nil
}

func (s *ServiceStatMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *ServiceStatMux) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {

	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
