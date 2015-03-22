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

func Identity(atom *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	var log = setupLogging(atom.ID)

	out = make(chan fc.IPOutbound)

	go func() {
		defer log.Notice("Identity element: exiting")
		defer close(out)
		log.Notice("Starting Identity element")

		for incoming := range in {
			switch incoming.GetTo().Name {
			case "input":
				log.Debug("Identity element: received data \n%s", spew.Sdump(incoming))
				out <- fc.NewIPOut(incoming.GetData(), "output")
				log.Info("Identity element: data sent")
			}
		}
	}()

	return out, nil
}

var _ fc.Element = Identity
