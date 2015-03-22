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
	"fmt"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("flowcker_stdlib")

func setupTestLogging() {

	var logging_backend = logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(
			"%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}",
		),
	)

	logging.SetBackend(logging_backend)
}

func setupLogging(atomID uint32) *logging.Logger {

	var log = logging.MustGetLogger("atom_" + fmt.Sprint(atomID))

	var logging_backend = logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(
			"Atom "+fmt.Sprint(atomID)+": %{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}",
		),
	)

	logging.SetBackend(logging_backend)

	return log
}
