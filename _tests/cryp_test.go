// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package _tests

import (
	"fmt"
	"github.com/dedeme/golib/cryp"
	"testing"
)

func TestGenK(t *testing.T) {
	ac := fmt.Sprintf("%d", len(cryp.GenK(12)))
	if r := eq(ac, "12"); r != "" {
		t.Fatal(r)
	}
	ac = fmt.Sprintf("%d", len(cryp.GenK(6)))
	if r := eq(ac, "6"); r != "" {
		t.Fatal(r)
	}
}

func TestKey(t *testing.T) {
	if r := eq(cryp.Key("deme", 6), "wiWTB9"); r != "" {
		t.Fatal(r)
	}
	if r := eq(cryp.Key("Generaro", 5), "Ixy8I"); r != "" {
		t.Fatal(r)
	}
	if r := eq(cryp.Key("Generara", 5), "0DIih"); r != "" {
		t.Fatal(r)
	}
}

func TestCryp(t *testing.T) {
	r := eq(cryp.Cryp("deme", "Cañón€%ç"), "v12ftuzYeq2Xz7q7tLe8tNnHtqY=")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("deme", cryp.Cryp("deme", "Cañón€%ç")), "Cañón€%ç")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("deme", cryp.Cryp("deme", "1")), "1")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("deme", cryp.Cryp("deme", "")), "")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("", cryp.Cryp("", "Cañón€%ç")), "Cañón€%ç")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("", cryp.Cryp("", "1")), "1")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("", cryp.Cryp("", "")), "")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("abc", cryp.Cryp("abc", "01")), "01")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("abcd", cryp.Cryp("abcd", "11")), "11")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("abc", cryp.Cryp("abc", "")), "")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("c", cryp.Cryp("c", "a")), "a")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("xxx", cryp.Cryp("xxx", "ab c")), "ab c")
	if r != "" {
		t.Fatal(r)
	}
	r = eq(cryp.Decryp("abc", cryp.Cryp("abc", "\n\ta€b c")), "\n\ta€b c")
	if r != "" {
		t.Fatal(r)
	}
}
