package handlerimplementation

import (
	"context"
	"io"
	"px.dev/pxapi"
	"px.dev/pxapi/types"
)

type ItemType interface {
	string | ServiceMap | Service | PodDetails | ServiceState | PodDetailsErrAndReq | PodDetailsLatency | PodDetailsCpuUsage | PixieTraceData
}

func New[itemType ItemType]() *ItemMapMux[itemType] {
	return &ItemMapMux[itemType]{}
}

type ItemMapMux[itemType ItemType] struct {
	Table TablePrinterItemMap[itemType]
}

type TablePrinterItemMap[itemType ItemType] struct {
	Values []itemType
}

func (t *TablePrinterItemMap[itemType]) HandleInit(ctx context.Context, metadata types.TableMetadata) error {
	return nil
}

func (t *TablePrinterItemMap[itemType]) HandleRecord(ctx context.Context, r *types.Record) error {
	itemMap := ConvertPixieDataToItemStore[itemType](r)
	t.Values = append(t.Values, itemMap)
	return nil
}

func (t *TablePrinterItemMap[itemType]) HandleDone(ctx context.Context) error {
	return nil
}

func (s *ItemMapMux[itemType]) AcceptTable(ctx context.Context, metadata types.TableMetadata) (pxapi.TableRecordHandler, error) {
	return &s.Table, nil
}

func (s *ItemMapMux[itemType]) ExecutePxlScript(ctx context.Context, vz *pxapi.VizierClient, pxl string) (*pxapi.ScriptResults, error) {
	resultSet, err := vz.ExecuteScript(ctx, pxl, s)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return resultSet, nil
}
