// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Text management.
package json

import (
	"errors"
)

func nextByte(s string, ch byte, ix int) (pos int, err error) {
	pos = ix
	l := len(s)
	quotes := false
	bar := false
	var c byte
	for {
		if pos == l {
			break
		}
		c = s[pos]
		if quotes {
			if bar {
				bar = false
			} else {
				if c == '\\' {
					bar = true
				} else if c == '"' {
					quotes = false
				}
			}
		} else {
			if c == ch {
				break
			} else if c == '"' {
				quotes = true
			} else if c == '[' {
				n := 1
				var open, close int
				for n > 0 {
					pos++
					open, err = nextByte(s, '[', pos)
					if err != nil {
						return
					}
					close, err = nextByte(s, ']', pos)
					if err != nil {
						return
					}
					if open < close {
						pos = open
						n++
					} else {
						pos = close
						n--
					}
				}
				if pos == l {
					err = errors.New("'[' not closed")
					return
				}
			} else if c == '{' {
				n := 1
				var open, close int
				for n > 0 {
					pos++
					open, err = nextByte(s, '{', pos)
					if err != nil {
						return
					}
					close, err = nextByte(s, '}', pos)
					if err != nil {
						return
					}
					if open < close {
						pos = open
						n++
					} else {
						pos = close
						n--
					}
				}
				if pos == l {
					err = errors.New("'{' not closed")
					return
				}
			}
		}
		pos++
	}
	return
}
