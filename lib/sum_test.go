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

package stdlib_test

import (
	"testing"

	fc "github.com/flowcker/flowcker/common"
	"github.com/flowcker/flowcker/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	atom := new(fc.Atom)
	atom.ID = 0

	in := make(chan fc.IPInbound)

	out, err := stdlib.Sum(atom, in)
	if err != nil {
		t.Fatalf("Error creating element: %s", err)
	}

	go func() {
		in <- fc.NewIPIn("10", fc.Port{Name: "input", Index: 0})
		in <- fc.NewIPIn("20", fc.Port{Name: "input", Index: 1})
		in <- fc.NewIPIn("5", fc.Port{Name: "input", Index: 2})
		in <- fc.NewIPIn("5", fc.Port{Name: "input", Index: 0})
		in <- fc.NewIPIn([]byte{}, fc.Port{Name: "input", Index: 1})
		in <- fc.NewIPIn("1000", fc.Port{Name: "input", Index: 2})
	}()

	var res fc.IPOutbound

	res = <-out
	assert.Equal(t, res.GetFrom().Name, "output")
	assert.False(t, res.GetIndexSelected())
	assert.False(t, res.GetAll())
	assert.Equal(t, string(res.GetData()), "10")
	res = <-out
	assert.Equal(t, res.GetFrom().Name, "output")
	assert.False(t, res.GetIndexSelected())
	assert.False(t, res.GetAll())
	assert.Equal(t, string(res.GetData()), "30")
	res = <-out
	assert.Equal(t, res.GetFrom().Name, "output")
	assert.False(t, res.GetIndexSelected())
	assert.False(t, res.GetAll())
	assert.Equal(t, string(res.GetData()), "35")
	res = <-out
	assert.Equal(t, res.GetFrom().Name, "output")
	assert.False(t, res.GetIndexSelected())
	assert.False(t, res.GetAll())
	assert.Equal(t, string(res.GetData()), "30")
	res = <-out
	assert.Equal(t, res.GetFrom().Name, "output")
	assert.False(t, res.GetIndexSelected())
	assert.False(t, res.GetAll())
	assert.Equal(t, string(res.GetData()), "1025")
}
