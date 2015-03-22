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
	"time"

	fc "github.com/flowcker/flowcker/common"
)
import "encoding/json"

func debounceChannel(interval time.Duration, output chan fc.IPOutbound) chan fc.IPOutbound {
	input := make(chan fc.IPOutbound)

	go func() {
		var buffer fc.IPOutbound
		var ok bool

		buffer, ok = <-input
		if !ok {
			return
		}
		for {
			select {
			case buffer, ok = <-input:
				if !ok {
					return
				}
			case <-time.After(interval):
				output <- buffer
				buffer, ok = <-input
				if !ok {
					return
				}
			}
		}
	}()

	return input
}

// Sum element recieves json unsigned integers in `input` port and save then in a internal table
// using the port index as index, with every number in the same port index overriding the
// previous value.
// For each incoming number, it will send via `output` port the sumatory of its table values
// Internally uses Golang uint64.
func Sum(atom *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	var log = setupLogging(atom.ID)

	out = make(chan fc.IPOutbound)
	var values = make(map[uint32]uint64)

	//debouncedOut := debounceChannel(time.Second, out)

	go func() {
		defer log.Debug("Sum element: exiting")
		defer close(out)
		for incoming := range in {
			switch incoming.GetTo().Name {
			case "input":
				log.Debug("Sum element: recieved data")
				if len(incoming.GetData()) == 0 {
					log.Debug("Sum element: recieved empty packet")
					continue
				}

				var value uint64
				var acc uint64
				json.Unmarshal(incoming.GetData(), &value)
				log.Debug("Sum element: recieved value: %d to port index %d", value, incoming.GetTo().Index)
				values[incoming.GetTo().Index] = value
				for _, member := range values {
					acc += member
				}
				log.Debug("Sum element: new sumatory value: %d", acc)

				data, _ := json.Marshal(acc)
				out <- fc.NewIPOut(data, "output")
			}
		}
	}()

	return out, nil
}
