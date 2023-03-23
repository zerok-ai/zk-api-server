package handlerimplementation

import (
	"context"
	"io"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
)

type PodDetailsListMux struct {
	Table TablePrinterPodDetailsList
}

type TablePrinterPodDetailsList struct {
	Values []PodDetails
}

func (t *TablePrinterPodDetailsList) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterPodDetailsList) HandleRecord(ctx context.Context, r *types.Record) error {
	t.Values = append(t.Values, ConvertPixieDataToPodDetails(r))
	return nil
}

func (t *TablePrinterPodDetailsList) HandleDone(ctx context.Context) error {
	return nil
}

func (s *PodDetailsListMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *PodDetailsListMux) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {

	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
