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
	"os"

	fc "github.com/flowcker/flowcker/common"
)

func ReadFile(atom *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	var log = setupLogging(atom.ID)

	out = make(chan fc.IPOutbound)

	go func() {
		log.Debug("ReadFile element: starting")
		defer log.Debug("ReadFile element: exiting")
		defer close(out)

		f, _ := os.Open("./input.txt")
		defer f.Close()
		r := bufio.NewReader(f)
		for {
			// TODO configurable delimiter
			data, err := r.ReadBytes('\n')
			if err == io.EOF {
				if len(data) > 0 {
					out <- fc.NewIPOut(data, "output")

					log.Debug("ReadFile element: data %s", string(data))
				}

				// Send empty package
				out <- fc.NewIPOut([]byte{}, "output")
				log.Debug("ReadFile element: sending empty package")
				return
			} else if err != nil {
				// TODO tell control atom about error
				return
			}

			out <- fc.NewIPOut(data, "output")
			log.Debug("ReadFile element: data %s", string(data))
		}
	}()

	return out, nil
}
