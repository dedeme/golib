// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// JSON utilities.
package json

import (
	"fmt"
	gson "encoding/json"
	"strings"
	"errors"
)

type T string

// json.T converter
func (js T) String () string {
	return string(js)
}

// Creates a new json.T from string ('s' must be a valid JSON value, although
// it is not checked).
func FromString (s string) T {
	return T(strings.TrimSpace(s))
}

func Wn () T {
	return "null"
}

func IsNull (js T) bool {
	return string(js) == "null"
}

// bool -> json.T
func Wb (v bool) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> bool
func Rb (js T) (v bool, err error) {
	err = gson.Unmarshal([]byte(js), &v)
	if err != nil {
		err = fmt.Errorf("%v in\n%v", err.Error(), string(js))
	}
	return
}

// int -> json.T
func Wi (v int) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> int
func Ri (js T) (v int, err error) {
	err = gson.Unmarshal([]byte(js), &v)
	if err != nil {
		err = fmt.Errorf("%v in\n%v", err.Error(), string(js))
	}
	return
}

// int64 -> json.T
func Wl (v int64) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> int64
func Rl (js T) (v int64, err error) {
	err = gson.Unmarshal([]byte(js), &v)
	if err != nil {
		err = fmt.Errorf("%v in\n%v", err.Error(), string(js))
	}
	return
}

// float32 -> json.T
func Wf (v float32) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> float64
func Rf (js T) (v float32, err error) {
	err = gson.Unmarshal([]byte(js), &v)
	if err != nil {
		err = fmt.Errorf("%v in\n%v", err.Error(), string(js))
	}
	return
}

// float64 -> json.T
func Wd (v float64) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> float64
func Rd (js T) (v float64, err error) {
	err = gson.Unmarshal([]byte(js), &v)
	if err != nil {
		err = fmt.Errorf("%v in\n%v", err.Error(), string(js))
	}
	return
}

// string -> json.T
func Ws (v string) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> string
func Rs (js T) (v string, err error) {
	err = gson.Unmarshal([]byte(js), &v)
	if err != nil {
		err = fmt.Errorf("%v in\n%v", err.Error(), string(js))
	}
	return
}

// []json.T -> json.T
func Wa (v []T) T {
	var b strings.Builder
	b.WriteByte('[')
	for i, js := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(string(js))
	}
	b.WriteByte(']')
	return T(b.String())
}

// json.T -> []json.T
func Ra (js T) (v []T, err error) {
	s := string(js)
	if !strings.HasPrefix(s, "[") {
		err = errors.New(fmt.Sprintf("Array does not start with '[' in\n%v", s))
		return
	}
	if !strings.HasSuffix(s, "]") {
		err = errors.New(fmt.Sprintf("Array does not end with ']' in\n%v", s))
		return
	}
	s2 := strings.TrimSpace(s[1:len(s)-1])
	l := len(s2)
	i := 0
	var e string
	var i2 int
	for i2 < l {
		i2 = nextByte(s2, ',', i)
		e = strings.TrimSpace(s2[i : i2])
		if e == "" {
			err = errors.New(fmt.Sprintf("Empty elements in\n%v", s))
			break
		}
		v = append(v, T(e))
		i = i2 + 1
	}
	return
}

// map[string]json.T -> json.T
func Wo (v map[string]T) T {
	var b strings.Builder
	b.WriteByte('{')
	more := false
	for k, js := range v {
		if more {
			b.WriteByte(',')
		} else {
			more = true
		}
		b.WriteString(string(Ws(k)))
		b.WriteByte(':')
		b.WriteString(string(js))
	}
	b.WriteByte('}')
	return T(b.String())
}

// json.T -> map[string]json.T
func Ro (js T) (v map[string]T, err error) {
	v = make(map[string]T)
	s := string(js)
	if !strings.HasPrefix(s, "{") {
		err = errors.New(fmt.Sprintf("Object does not start with '{' in\n%v", s))
		return
	}
	if !strings.HasSuffix(s, "}") {
		err = errors.New(fmt.Sprintf("Object does not end with '}' in\n%v", s))
		return
	}
	s2 := strings.TrimSpace(s[1:len(s)-1])
	l := len(s2)
	i := 0
	var kjs string
	var k string
	var val string
	var i2 int
	for i2 < l {
		i2 = nextByte(s2, ':', i)
		kjs = strings.TrimSpace(s2[i : i2])
		if kjs == "" {
			err = errors.New(fmt.Sprintf("Key missing in\n%v", s))
			break
		}
		k, err = Rs(T(kjs))
		if err != nil {
			break
		}

		i = i2 + 1
		i2 = nextByte(s2, ',', i)
		val = strings.TrimSpace(s2[i : i2])
		if val == "" {
			err = errors.New(fmt.Sprintf("Value missing in\n%v", s))
			break
		}

		v[k] = T(val)
		i = i2 + 1
	}
	return
}

