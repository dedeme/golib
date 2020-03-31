// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Text management.
package json

func nextByte (s string, ch byte, ix int) (pos int) {
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
			}
		}
		pos++
	}
	return
}
