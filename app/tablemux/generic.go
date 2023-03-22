package tablemux

import (
	"px.dev/pxapi"
)

type TableMux[C pxapi.TableMuxer] struct {
	Muxer C
}
