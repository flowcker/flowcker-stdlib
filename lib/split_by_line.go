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
	"bufio"
	"io"

	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
)

func SplitByLine(atom *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	var log = setupLogging(atom.ID)

	out = make(chan fc.IPOutbound)
	rd, wr := io.Pipe()

	go func() {
		defer log.Notice("SplitByLine element incoming loop: exiting")
		defer wr.Close()

		var count uint64

		log.Debug("Starting SplitByLine element incoming loop")
		for incoming := range in {
			switch incoming.GetTo().Name {
			case "input":
				count++
				log.Info("SplitByLine element: received data count %d", count)
				log.Debug("SplitByLine element: received data \n%s\n%s", spew.Sdump(incoming.GetData()), spew.Sdump(incoming))
				n, err := wr.Write(incoming.GetData())
				if err != nil {
					log.Error("SplitByLine element: error writing data on pipe")
					panic(err)
				}
				if n < len(incoming.GetData()) {
					log.Error("SplitByLine element: error writing less data on pipe %d < %d", n, len(incoming.GetData()))
					panic(n)
				}
			}
		}
	}()

	sendLine := func(line []byte, count *uint64) {
		if len(line) <= 0 {
			return
		}

		out <- fc.NewIPOut(line, "output")
		*count++
		log.Debug("SplitByLine element outgoing loop: data %s", string(line))
		log.Info("SplitByLine element outgoing loop: count %d", *count)
	}

	go func() {
		defer log.Notice("SplitByLine element outgoing loop: exiting")
		defer close(out)

		r := bufio.NewReader(rd)
		scanner := bufio.NewScanner(r)
		log.Notice("SplitByLine element outgoing loop: starting")

		var count uint64
		var line []byte

		for {
			if !scanner.Scan() {
				if scanner.Err() == nil {
					break
				}
				panic(scanner.Err())
			}
			line = scanner.Bytes()
			sendLine(line, &count)
		}

	}()

	return out, nil
}
