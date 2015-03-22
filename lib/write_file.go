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
	"os"

	fc "github.com/flowcker/flowcker/common"
)

type writeFileConfig struct {
	path string
}

func WriteFile(atom *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	var log = setupLogging(atom.ID)

	out = make(chan fc.IPOutbound)
	var config writeFileConfig

	config.path = "./output.txt"
	json.Unmarshal(*atom.Config, &config)

	go func() {
		f, _ := os.Create(config.path)
		defer f.Close()

		defer log.Debug("WriteFile element: exiting")
		defer close(out)

		log.Debug("Starting WriteFile element")
		for incoming := range in {
			switch incoming.GetTo().Name {
			case "input":
				log.Debug("WriteFile element: received data")
				f.Write(incoming.GetData())
				f.Write([]byte("\n"))
				log.Debug("WriteFile element: data saved")
			}
		}
	}()

	return out, nil
}
