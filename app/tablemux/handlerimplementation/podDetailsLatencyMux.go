package handlerimplementation

import (
	"context"
	"io"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
)

type PodDetailsLatencyMux struct {
	Table TablePrinterPodDetailsLatency
}

type TablePrinterPodDetailsLatency struct {
	Values []PodDetailsLatency
}

func (t *TablePrinterPodDetailsLatency) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterPodDetailsLatency) HandleRecord(ctx context.Context, r *types.Record) error {
	t.Values = append(t.Values, ConvertPixieDataToPodDetailsLatency(r))
	return nil
}

func (t *TablePrinterPodDetailsLatency) HandleDone(ctx context.Context) error {
	return nil
}

func (s *PodDetailsLatencyMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *PodDetailsLatencyMux) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {

	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
