package parser

import (
	"io"

	"github.com/bcicen/jstream"

	"port-location/internal/common/model"
)

// ReadPortInfo is based on io.Reader for easier usage (thus we can provide request.Body instead of file if necessary)
func ReadPortInfo(r io.Reader) (<-chan model.Port, <-chan error) {
	portCh := make(chan model.Port)
	errCh := make(chan error)

	// jstream lib allows to stream json rows as raw <string,interface> pairs. Used due to limited memory limit
	dec := jstream.NewDecoder(r, 1).EmitKV()

	go func() {
		for entry := range dec.Stream() {
			d := entry.Value.(jstream.KV)
			locode, portInfo := d.Key, d.Value.(map[string]interface{})

			p, err := toModelPort(locode, portInfo)
			if err != nil {
				errCh <- err
			}

			portCh <- p
		}

		close(portCh)
	}()

	return portCh, errCh
}
