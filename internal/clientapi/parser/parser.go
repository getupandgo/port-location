package parser

import (
	"io"

	"github.com/bcicen/jstream"

	"port-location/internal/clientapi/model"
)

func ReadPortInfo(r io.Reader) (<-chan model.Port, <-chan error) {
	portCh := make(chan model.Port)
	errCh := make(chan error)

	dec := jstream.NewDecoder(r, 1).EmitKV()

	go func() {
		for entry := range dec.Stream() {
			d := entry.Value.(jstream.KV)
			locode, portInfo := d.Key, d.Value.(map[string]interface{})

			portCh <- toModelPort(locode, portInfo)
		}

		close(portCh)
	}()

	return portCh, errCh
}

//func SaveUnprocessedPort(w io.Writer, port model.Port) error {
//	return nil
//}
