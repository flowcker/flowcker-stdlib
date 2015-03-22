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
	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
)
import "encoding/json"

// AccumAdd element recieves json unsigned integers in `input` port and adds them in an accumulator
// for each incoming number, it will send via `output` port the content of the accumulator after
// the adding the last incomming integer.
// Internally uses Golang uint64.
func AccumAdd(atom *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	var log = setupLogging(atom.ID)

	out = make(chan fc.IPOutbound)

	go func() {
		var acc uint64
		var count uint64

		defer log.Debug("AccumAdd element: exiting")
		defer close(out)
		for incoming := range in {
			switch incoming.GetTo().Name {
			case "input":
				log.Debug("AccumAdd element: recieved data")
				if len(incoming.GetData()) == 0 {
					log.Debug("AccumAdd element: recieved empty packet\n%s", spew.Sdump(incoming))
					count++
				} else {
					var value uint64
					err := json.Unmarshal(incoming.GetData(), &value)
					if err != nil {
						log.Error("Error unmarshalling incoming data: \n%s", spew.Sdump(incoming.GetData()))
						panic(err)
					}
					log.Debug("AccumAdd element: recieved value: %d", value)

					acc += value
					count++
					log.Debug("AccumAdd element: new accumulator value: %d", acc)
					log.Debug("AccumAdd element: count: %d", count)

					data, _ := json.Marshal(acc)
					out <- fc.NewIPOut(data, "output")
				}
			}
		}
	}()

	return out, nil
}
