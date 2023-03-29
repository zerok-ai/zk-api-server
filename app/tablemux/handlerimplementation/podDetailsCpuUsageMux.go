package handlerimplementation

import (
	"context"
	"io"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
)

type PodDetailsCpuUsageMux struct {
	Table TablePrinterPodDetailsCpuUsage
}

type TablePrinterPodDetailsCpuUsage struct {
	Values []PodDetailsCpuUsage
}

func (t *TablePrinterPodDetailsCpuUsage) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterPodDetailsCpuUsage) HandleRecord(ctx context.Context, r *types.Record) error {
	t.Values = append(t.Values, ConvertPixieDataToPodDetailsCpuUsage(r))
	return nil
}

func (t *TablePrinterPodDetailsCpuUsage) HandleDone(ctx context.Context) error {
	return nil
}

func (s *PodDetailsCpuUsageMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *PodDetailsCpuUsageMux) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {

	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
