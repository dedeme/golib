// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// JSON utilities.
package json

import (
	gson "encoding/json"
	"fmt"
	"strings"
)

type T string

// json.T converter
func (js T) String() string {
	return string(js)
}

// Creates a new json.T from string ('s' must be a valid JSON value, although
// it is not checked).
func FromString(s string) T {
	return T(strings.TrimSpace(s))
}

func Wn() T {
	return "null"
}

func (js T) IsNull() bool {
	return string(js) == "null"
}

// bool -> json.T
func Wb(v bool) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> bool
func (js T) Rb() (v bool) {
	err := gson.Unmarshal([]byte(js), &v)
	if err != nil {
		panic(fmt.Sprintf("%v in\n'%v'", err.Error(), string(js)))
	}
	return
}

// int -> json.T
func Wi(v int) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> int
func (js T) Ri() (v int) {
	err := gson.Unmarshal([]byte(js), &v)
	if err != nil {
		panic(fmt.Sprintf("%v in\n'%v'", err.Error(), string(js)))
	}
	return
}

// int64 -> json.T
func Wl(v int64) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> int64
func (js T) Rl() (v int64) {
	err := gson.Unmarshal([]byte(js), &v)
	if err != nil {
		panic(fmt.Sprintf("%v in\n'%v'", err.Error(), string(js)))
	}
	return
}

// float32 -> json.T
func Wf(v float32) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> float64
func (js T) Rf() (v float32) {
	err := gson.Unmarshal([]byte(js), &v)
	if err != nil {
		panic(fmt.Sprintf("%v in\n'%v'", err.Error(), string(js)))
	}
	return
}

// float64 -> json.T
func Wd(v float64) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> float64
func (js T) Rd() (v float64) {
	err := gson.Unmarshal([]byte(js), &v)
	if err != nil {
		panic(fmt.Sprintf("%v in\n'%v'", err.Error(), string(js)))
	}
	return
}

// string -> json.T
func Ws(v string) T {
	js, _ := gson.Marshal(v)
	return T(js)
}

// json.T -> string
func (js T) Rs() (v string) {
	err := gson.Unmarshal([]byte(js), &v)
	if err != nil {
		panic(fmt.Sprintf("%v in\n'%v'", err.Error(), string(js)))
	}
	return
}

// []json.T -> json.T
func Wa(v []T) T {
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
func (js T) Ra() (v []T) {
	s := string(js)
	if !strings.HasPrefix(s, "[") {
		panic(fmt.Sprintf("Array does not start with '[' in\n'%v'", s))
	}
	if !strings.HasSuffix(s, "]") {
		panic(fmt.Sprintf("Array does not end with ']' in\n'%v'", s))
	}
	s2 := strings.TrimSpace(s[1 : len(s)-1])
	l := len(s2)
	if l == 0 {
		return
	}
	i := 0
	var e string
	for {
		if i2, ok := nextByte(s2, ',', i); ok {
      e = strings.TrimSpace(s2[i:i2])
      if e == "" {
        panic(fmt.Sprintf("Missing elements in\n'%v'", s))
      }
      v = append(v, T(e))
      i = i2 + 1
      continue
    }
    e = strings.TrimSpace(s2[i:l])
    if e == "" {
      panic(fmt.Sprintf("Missing elements in\n'%v'", s))
    }
    v = append(v, T(e))
    break
	}
	return
}

// map[string]json.T -> json.T
func Wo(v map[string]T) T {
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
func (js T) Ro() (v map[string]T) {
	v = make(map[string]T)
	s := string(js)
	if !strings.HasPrefix(s, "{") {
		panic(fmt.Sprintf("Object does not start with '{' in\n'%v'", s))
	}
	if !strings.HasSuffix(s, "}") {
		panic(fmt.Sprintf("Object does not end with '}' in\n'%v'", s))
	}
	s2 := strings.TrimSpace(s[1 : len(s)-1])
	l := len(s2)
  if l == 0 {
    return
  }
	i := 0
	var kjs string
	var k string
	var val string
	for {
		i2, ok := nextByte(s2, ':', i)
    if !ok {
			panic(fmt.Sprintf("Expected ':' in\n'%v'", s2))
		}
		kjs = strings.TrimSpace(s2[i:i2])
		if kjs == "" {
			panic(fmt.Sprintf("Key missing in\n'%v'", s))
		}
		k = T(kjs).Rs()

		i = i2 + 1

		if i2, ok := nextByte(s2, ',', i); ok {
      val = strings.TrimSpace(s2[i:i2])
      if val == "" {
        panic(fmt.Sprintf("Value missing in\n'%v'", s))
      }
      v[k] = T(val)
      i = i2 + 1
      continue
    }
    val = strings.TrimSpace(s2[i:l])
    if val == "" {
      panic(fmt.Sprintf("Value missing in\n'%v'", s))
    }
    v[k] = T(val)
    break
	}
	return
}
