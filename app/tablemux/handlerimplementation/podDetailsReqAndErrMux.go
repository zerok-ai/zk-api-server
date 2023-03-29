package handlerimplementation

import (
	"context"
	"io"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
)

type PodDetailsReqAndErrMux struct {
	Table TablePrinterPodDetailsReqAndErr
}

type TablePrinterPodDetailsReqAndErr struct {
	Values []PodDetailsErrAndReq
}

func (t *TablePrinterPodDetailsReqAndErr) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterPodDetailsReqAndErr) HandleRecord(ctx context.Context, r *types.Record) error {
	t.Values = append(t.Values, ConvertPixieDataToPodDetailsErrAndReq(r))
	return nil
}

func (t *TablePrinterPodDetailsReqAndErr) HandleDone(ctx context.Context) error {
	return nil
}

func (s *PodDetailsReqAndErrMux) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *PodDetailsReqAndErrMux) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {

	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
