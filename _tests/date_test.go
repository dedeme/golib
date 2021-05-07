// Copyright 07-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package _tests

import (
	"github.com/dedeme/golib/date"
	"sort"
	"testing"
  "fmt"
)

func TestToFrom(t *testing.T) {
	d := date.Now()
	if r := eqi(len(d.String()), 8); r != "" {
		t.Fatal(r)
	}

	d = date.New(2, 4, 2010)
	if r := eq(d.String(), "20100402"); r != "" {
		t.Fatal(r)
	}
	if r := eq(d.ToIso(), "02/04/2010"); r != "" {
		t.Fatal(r)
	}
	if r := eq(d.ToEn(), "04/02/2010"); r != "" {
		t.Fatal(r)
	}

	d1 := date.FromEn("12/31/1988", "/")
	d2 := date.FromIso("31/12/1988", "/")
	if !d1.Eq(d2) {
		t.Fatal(fail)
	}
	if r := eq(d1.String(), "19881231"); r != "" {
		t.Fatal(r)
	}
	if r := eq(d1.ToIso(), "31/12/1988"); r != "" {
		t.Fatal(r)
	}
	if r := eq(d2.ToEn(), "12/31/1988"); r != "" {
		t.Fatal(r)
	}

	d1 = date.NewTime(2, 4, 2010, 12, 40, 15)
	if r := eq(d.String(), "20100402"); r != "" {
		t.Fatal(r)
	}
	if r := eq(d1.Format("%t"), "12:40:15"); r != "" {
		t.Fatal(r)
	}
	if r := eq(d1.Format("%T in %%"), "12:40:15.000 in %"); r != "" {
		t.Fatal(r)
	}
  if !d1.EqTime(date.FromJs(d1.ToJs())) {
    fmt.Println(d1)
    fmt.Println(date.FromJs(d1.ToJs()))
    t.Fatal(fail)
  }
}

func TestOperations(t *testing.T) {
	d1 := date.New(29, 2, 2013)
	d2 := date.New(6, 3, 2013)
	d3 := date.New(30, 4, 2013)

	if r := eqi(d1.Day(), 1); r != "" {
		t.Fatal(r)
	}
	if r := eqi(d1.Month(), 3); r != "" {
		t.Fatal(r)
	}
	if r := eqi(d1.Year(), 2013); r != "" {
		t.Fatal(r)
	}
	if r := eqi(d1.Hour(), 12); r != "" {
		t.Fatal(r)
	}
	if r := eqi(d1.Minute(), 0); r != "" {
		t.Fatal(r)
	}
	if r := eqi(d1.Second(), 0); r != "" {
		t.Fatal(r)
	}

	if d1.Eq(d2) {
		t.Fatal(fail)
	}
	if d1.Compare(d2) >= 0 {
		t.Fatal(fail)
	}
	if r := eqi(d1.Df(d2), -5); r != "" {
		t.Fatal(r)
	}

	if d3.Eq(d2) {
		t.Fatal(fail)
	}
	if d3.Compare(d2) <= 0 {
		t.Fatal(fail)
	}
	if r := eqi(d3.Df(d2), 55); r != "" {
		t.Fatal(r)
	}

	if r := eqi(d1.Add(25).Add(-25).Df(d1), 0); r != "" {
		t.Fatal(r)
	}
	if r := eqi(d1.Add(25).Add(-30).Df(d1), -5); r != "" {
		t.Fatal(r)
	}
	if r := eqi(d1.Add(25).Add(-20).Df(d1), 5); r != "" {
		t.Fatal(r)
	}

	if r := eqi(d1.AddSeconds(25).AddSeconds(-25).DfTime(d1), 0); r != "" {
		t.Fatal(r)
	}
	if r := eqi(d1.AddSeconds(25).AddSeconds(-30).DfTime(d1), -5000); r != "" {
		t.Fatal(r)
	}
	if r := eqi(d1.AddSeconds(25).AddSeconds(-20).DfTime(d1), 5000); r != "" {
		t.Fatal(r)
	}

	ds := []date.T{d3, d1, d2}
	sort.Sort(date.Sorter(ds))
	if !ds[0].Eq(d1) {
		t.Fatal(fail)
	}
	if !ds[1].Eq(d2) {
		t.Fatal(fail)
	}
	if !ds[2].Eq(d3) {
		t.Fatal(fail)
	}

	ds = []date.T{d3, d1, d2}
	sort.Sort(date.SorterTime(ds))
	if !ds[0].Eq(d1) {
		t.Fatal(fail)
	}
	if !ds[1].Eq(d2) {
		t.Fatal(fail)
	}
	if !ds[2].Eq(d3) {
		t.Fatal(fail)
	}

}
