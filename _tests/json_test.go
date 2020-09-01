// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package _tests

import (
	"fmt"
	"github.com/dedeme/golib/json"
	"testing"
)

func TestNull(t *testing.T) {
	if json.FromString("xxx").IsNull() {
		t.Fatal(fail)
	}
	if !json.FromString("null").IsNull() {
		t.Fatal(fail)
	}
	if r := eq(json.Wn().String(), "null"); r != "" {
		t.Fatal(r)
	}
}

func TestBool(t *testing.T) {
	var test = func(value string) {
		v := json.FromString(value).Rb()
		if r := eq(json.Wb(v).String(), value); r != "" {
			t.Fatal(r)
		}
	}

	test("true")
	test("false")
}

func TestInt(t *testing.T) {
	var test = func(value string) {
		v := json.FromString(value).Ri()
		if r := eq(json.Wi(v).String(), value); r != "" {
			t.Fatal(r)
		}
	}

	test("0")
	test("123")
	test("-2500")
}

func TestLong(t *testing.T) {
	var test = func(value string) {
		v := json.FromString(value).Rl()
		if r := eq(json.Wl(v).String(), value); r != "" {
			t.Fatal(r)
		}
	}

	test("0")
	test("123")
	test("-2500")
}

func TestFloat(t *testing.T) {
	var test = func(value string) {
		v := json.FromString(value).Rf()
		if r := eq(json.Wf(v).String(), value); r != "" {
			t.Fatal(r)
		}
	}

	test("0")
	test("123.456")
	test("-2500.02")
}

func TestDouble(t *testing.T) {
	var test = func(value string) {
		v := json.FromString(value).Rd()
		if r := eq(json.Wd(v).String(), value); r != "" {
			t.Fatal(r)
		}
	}

	test("0")
	test("123.456")
	test("-2500.02")
}

func TestString(t *testing.T) {
	var test = func(value string) {
		v := json.FromString(value).Rs()
		if r := eq(json.Ws(v).String(), value); r != "" {
			t.Fatal(r)
		}
	}

	test("\"\"")
	test("\"abc\"")
	rs := "\" \\tcañón\\\"ŧ\\\"\\n\""
	test(rs)
	if r := eq(json.Ws(" \tcañón\"ŧ\"\n").String(), rs); r != "" {
		t.Fatal(r)
	}
}

func TestArray(t *testing.T) {
	mkErr := func(js json.T) (err string) {
		defer func() {
			r := recover()
			switch r.(type) {
			case string:
				err = r.(string)
			default:
				err = fmt.Sprintf("No errors found in '%v'", js)
			}
		}()
		js.Ra()
		return
	}

	var test = func(value string) {
		v := json.FromString(value).Ra()
		if r := eq(json.Wa(v).String(), value); r != "" {
			t.Fatal(r)
		}
	}

	test("[]")
	test("[1]")
	test("[1,2,3]")
	test("[1,\"abc\",2]")
	test("[1,\"abñc\",2]")
	test("[1,\"a\\\"b\\\"ñc\",2]")
	test("[1,[2,[3,4]],[2,3]]")

	err := mkErr(json.FromString("[4"))
	msg := "Array does not end with ']' in\n'[4'"
	if r := eq(err, msg); r != "" {
		t.Fatal(r)
	}

	err = mkErr(json.FromString("4]"))
	msg = "Array does not start with '[' in\n'4]'"
	if r := eq(err, msg); r != "" {
		t.Fatal(r)
	}

	err = mkErr(json.FromString("[,]"))
	msg = "Missing elements in\n'[,]'"
	if r := eq(err, msg); r != "" {
		t.Fatal(r)
	}
	err = mkErr(json.FromString("[1,,2]"))
	msg = "Missing elements in\n'[1,,2]'"
	if r := eq(err, msg); r != "" {
		t.Fatal(r)
	}

	err = mkErr(json.FromString("[,1,2]"))
	msg = "Missing elements in\n'[,1,2]'"
	if r := eq(err, msg); r != "" {
		t.Fatal(r)
	}

	err = mkErr(json.FromString("[1,2,]"))
	msg = "Missing elements in\n'[1,2,]'"
	if r := eq(err, msg); r != "" {
		t.Fatal(r)
	}

	jss := json.FromString("[345]").Ra()
	jss[0].Ri()
	msg = "invalid character 'a' looking for beginning of value in\n'a4'"

	// jss = json.Ra(json.FromString("[a4]"))
	// json.Ri(jss[0]) // error

	jss = json.FromString("[1,\"a\\\"b\\\"ñc\",2]").Ra()
	jss[1].Rs()

	jss = json.FromString("[1,\"a\\\"b\\\"ñc,2]").Ra() // error
	//jss[1].Rs() // error

	jss = json.FromString("[1,\"a\\\"b\\\"ñc\",2]").Ra()
	//jss[2].Rb() // error

}

func TestObject(t *testing.T) {
	var test = func(value string, r map[string]json.T) {
		v := json.FromString(value).Ro()

		jss := json.Wo(v)
		mjs := jss.Ro()
		if len(mjs) != len(r) {
			t.Fatal(eq(jss.String(), value))
		}

		for k, v := range mjs {
			val := r[k]
			if v.String() != val.String() {
				t.Fatal(eq(jss.String(), value))
			}
		}
	}

	test("{}", map[string]json.T{})
	test("{\"one\":1}", map[string]json.T{"one": json.Wi(1)})
	test("{\"one\":1,\"two\":2,\"three\":3}",
		map[string]json.T{
			"one":   json.Wi(1),
			"two":   json.Wi(2),
			"three": json.Wi(3),
		},
	)
	test("{\"one\":1,\"two\":[2,3],\"three\":3}",
		map[string]json.T{
			"one":   json.Wi(1),
			"two":   json.Wa([]json.T{json.Wi(2), json.Wi(3)}),
			"three": json.Wi(3),
		},
	)
	/*
		test("{\"one\":1,\"two\":{\"a\":2,\"b\":3},\"three\":3}",
			map[string]json.T{
				"one": json.Wi(1),
				"two":   json.Wo(map[string]json.T{"a":json.Wi(2),"b:":json.Wi(3)}),
				"three": json.Wi(3),
			},
		)
	*/
	test("{\"one\":1,\"two\":\"abc\",\"three\":3}",
		map[string]json.T{
			"one":   json.Wi(1),
			"two":   json.Ws("abc"),
			"three": json.Wi(3),
		},
	)
	test("{\"one\":1,\"two\":\"abñc\",\"three\":3}",
		map[string]json.T{
			"one":   json.Wi(1),
			"two":   json.Ws("abñc"),
			"three": json.Wi(3),
		},
	)
	test("{\"one\":1,\"two\":\"a\\\"b\\\"ñc\",\"three\":3}",
		map[string]json.T{
			"one":   json.Wi(1),
			"two":   json.Ws("a\"b\"ñc"),
			"three": json.Wi(3),
		},
	)

	//json.FromString("{\"o\":4").Ro()
	//---> Error "Object does not end with '}' in\n{\"o\":4"

	//son.FromString("\"o\":4}").Ro()
	//---> Error "Object does not start with '{' in\n\"o\":4}"

	//json.FromString("{\"a\":1,\"b\":2,}").Ro()
	//---> Error "Expected ':' in\n{\"a\":1,\"b\":2,}"

	//json.FromString("{\"a\":1,\"b\":2,\"c\":}").Ro()
	//---> Error "Value missing in\n{\"a\":1,\"b\":2,\"c\":}"

	jss := json.FromString("{\"a\":345}").Ro()
	jss["a"].Ri() // Ok

	//jss = json.FromString("{\"a\":a345}").Ro()
	//jss["a"].Ri()
	//---> Error invalid character 'a' looking for beginning of value in

	jss = json.FromString("{\"a\":1,\"b\":\"a\\\"b\\\"ñc\",\"c\":2}").Ro()
	jss["b"].Rs() // Ok

	//jss = json.FromString("{\"a\":1,\"b\":\"a\\\"b\\\"ñc,\"c\":2}").Ro()
	//jss["b"].Rs()
	//---> Error invalid character 'c' after top-level value in

	json.FromString("{\"a\":1,\"b\":\"a\\\"b\\\"ñc\",\"c\":2}").Ro()
	//jss["c"].Rb()
	//---> Error json: cannot unmarshal number into Go value of type bool in

}
