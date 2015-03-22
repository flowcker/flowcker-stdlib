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
	"strings"

	"github.com/flowcker/flowcker/atom"
	fc "github.com/flowcker/flowcker/common"
)

var elements = map[string]fc.Element{
	"identity":       Identity,
	"accumadd":       AccumAdd,
	"readfile":       ReadFile,
	"writefile":      WriteFile,
	"readhttp":       ReadHTTP,
	"updaterediskey": UpdateRedisKey,
	"splitbybyte":    SplitByByte,
	"splitbyline":    SplitByLine,
	"sum":            Sum,
}

// Main is the entry point to execute different elements of the Stdlib from the
// command line. It runs the element passed as argument.
func Main() {
	if len(os.Args) != 2 {
		fmt.Println("Must pass element name")
		implementedElements()
		return
	}
	elementName := strings.ToLower(os.Args[1])
	element, ok := elements[elementName]
	if !ok {
		fmt.Println("Element not found " + elementName)
		implementedElements()
		return
	}

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if port == "" {
		port = "0"
	}
	if host == "" {
		host = "0.0.0.0"
	}

	atom, portsServer := atom.LaunchTCP(element, host+":"+port)
	fmt.Println(portsServer.GetAddr().String())

	<-atom.CloseChannel()
}

func implementedElements() {
	fmt.Println("Implemented elements:")
	for e := range elements {
		fmt.Println("  " + e)
	}
}
