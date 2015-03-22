/*  Copyright (C) 2015 Leopoldo Lara Vazquez.

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package stdlib

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
)

type readHTTPConfig struct {
	URL string
}

// ReadHTTP reads a URL content and sends it via `output` port
// Configuration:
//  - url string // the url
func ReadHTTP(atom *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	var log = setupLogging(atom.ID)
	var count uint32
	out = make(chan fc.IPOutbound)

	var config readHTTPConfig
	config.URL = "CONFIGURE"
	json.Unmarshal(*atom.Config, &config)

	sendBuffer := func(buff []byte) {
		if len(buff) <= 0 {
			return
		}
		count++
		sendbuff := make([]byte, len(buff))
		copy(sendbuff, buff)
		out <- fc.NewIPOut(sendbuff, "output")
		log.Info("ReadHTTP element: sent data out count %d", count)
		log.Debug("ReadHTTP element: data %s", spew.Sdump(sendbuff))
	}

	go func() {
		log.Notice("ReadHTTP element: starting")
		log.Info("ReadHTTP url is " + config.URL)
		defer close(out)

		// Do HTTP GET request
		resp, err := http.Get(config.URL)
		if err != nil {
			log.Fatal("Error getting file by HTTP")
			panic(err)
		}
		defer resp.Body.Close()

		var buff [1024 * 4]byte

		r := resp.Body

		for {
			// Read from the stream and send it in a IP
			n, err := r.Read(buff[:])
			sendBuffer(buff[:n])

			if err == io.EOF {
				break
			} else if err != nil {
				// TODO tell control atom about error
				panic(err)
			}
		}

		log.Notice("ReadHTTP element: exiting")
	}()

	return out, nil
}
