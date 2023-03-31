package handlerimplementation

import (
	"context"
	"io"
	"main/app/utils"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
)

type PixieTraceDataListMux struct {
	Table TablePrinterPixieTraceDataList
}

type TablePrinterPixieTraceDataList struct {
	Values []PixieTraceData
}

func (t *TablePrinterPixieTraceDataList) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterPixieTraceDataList) HandleRecord(ctx context.Context, r *types.Record) error {
	d := ConvertPixieDataToPixieTraceData(r)
	if utils.IsEmpty(d.TraceId) || utils.IsEmpty(d.SpanId) {
		return nil
	}
	t.Values = append(t.Values, d)
	return nil
}

func (t *TablePrinterPixieTraceDataList) HandleDone(ctx context.Context) error {
	return nil
}

func (s *PixieTraceDataListMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *PixieTraceDataListMux) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {

	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
