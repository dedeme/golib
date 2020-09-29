// Copyright 07-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Date - time management
package date

import (
	"strings"
	"time"
)

type T time.Time

func Now() T {
	return T(time.Now())
}

// 'day' is between 1 and 31 both inclusive. It it is out of range, date rolls
// to the corresponding valid date.
// 'month' is between 1 and 12 both inclusive. If month is out of range, a
// 'panic' is raised.
func New(day, month, year int) T {
	return NewTime(day, month, year, 12, 0, 0)
}

// 'day' is between 1 and 31 both inclusive. It it is out of range, date rolls
// to the corresponding valid date.
// 'month' is between 1 and 12 both inclusive. If month is out of range, a
// 'panic' is raised.
func NewTime(day, month, year, hour, min, sec int) T {
	lc, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	return T(time.Date(year, time.Month(month), day, hour, min, sec, 0, lc))
}

// 's' is in format "yyyymmdd"
// If there are some errors, raise 'panic'
func FromString(s string) T {
	d, err := time.Parse("20060102", s)
	if err != nil {
		panic(err)
	}
	return T(d)
}

// For 'sep' = "/", 's' is in format "dd/mm/yyyy"
// If there are some errors, raise 'panic'
func FromIso(s, sep string) T {
	d, err := time.Parse("02"+sep+"01"+sep+"2006", s)
	if err != nil {
		panic(err)
	}
	return T(d)
}

// For 'sep' = "/", 's' is in format "mm/dd/yyyy"
// If there are some errors, raise 'panic'
func FromEn(s, sep string) T {
	d, err := time.Parse("01"+sep+"02"+sep+"2006", s)
	if err != nil {
		panic(err)
	}
	return T(d)
}

func (d T) Day() int {
	return time.Time(d).Day()
}

// Returns the week day. 0 -> Sunday ... 6 -> Saturday
func (d T) Weekday() int {
  return int(time.Time(d).Weekday())
}

// In the range [1-12]
func (d T) Month() int {
	return int(time.Time(d).Month())
}

func (d T) Year() int {
	return time.Time(d).Year()
}

func (d T) Hour() int {
	return time.Time(d).Hour()
}

func (d T) Minute() int {
	return time.Time(d).Minute()
}

func (d T) Second() int {
	return time.Time(d).Second()
}

func (d T) Millisecond() int {
	return time.Time(d).Nanosecond() / 1000000
}

func (d T) Add(days int) T {
	return T(time.Time(d).Add(time.Hour * 24 * time.Duration(days)))
}

func (d T) AddSeconds(seconds int) T {
	return T(time.Time(d).Add(time.Second * time.Duration(seconds)))
}

func (d T) AddMilliseconds(millis int) T {
	return T(time.Time(d).Add(time.Millisecond * time.Duration(millis)))
}

// Milliseconds difference. d - other.
func (d T) DfTime(other T) int {
	return int(time.Time(d).Sub(time.Time(other)).Milliseconds())
}

// Days difference. d - other.
func (d T) Df(other T) int {
	df := d.DfTime(other)
	dv := df / 86400000
	if df >= 0 {
		if df%86400000 >= 43200000 {
			dv++
		}
	} else {
		if df%86400000 <= -43200000 {
			dv--
		}
	}

	return dv
}

// Compare up to day.
func (d T) Eq(other T) bool {
	return d.Day() == other.Day() &&
		d.Month() == other.Month() &&
		d.Year() == other.Year()
}

// Compare up to milliseconds.
func (d T) EqTime(other T) bool {
	return d.Day() == other.Day() &&
		d.Month() == other.Month() &&
		d.Year() == other.Year() &&
		d.Hour() == other.Hour() &&
		d.Minute() == other.Minute() &&
		d.Second() == other.Second() &&
		d.Millisecond() == other.Millisecond()
}

// Compare up to day
func (d T) Compare(other T) int {
	if d.Year() > other.Year() {
		return 1
	}
	if d.Year() < other.Year() {
		return -1
	}
	if d.Month() > other.Month() {
		return 1
	}
	if d.Month() < other.Month() {
		return -1
	}
	if d.Day() > other.Day() {
		return 1
	}
	if d.Day() < other.Day() {
		return -1
	}
	return 0
}

// Compare up to milliseconds.
func (d T) CompareTime(other T) int {
	cmp := d.Compare(other)
	if cmp == 0 {
		if d.Hour() > other.Hour() {
			return 1
		}
		if d.Hour() < other.Hour() {
			return -1
		}
		if d.Minute() > other.Minute() {
			return 1
		}
		if d.Minute() < other.Minute() {
			return -1
		}
		if d.Second() > other.Second() {
			return 1
		}
		if d.Second() < other.Second() {
			return -1
		}
		if d.Millisecond() < other.Millisecond() {
			return -1
		}
		return 0
	}
	return cmp
}

// Interface to sort dates (Compare up to day). Example:
//    sort.Sort(date.Sorter(sliceOfDates))
type Sorter []T

func (a Sorter) Len() int {
	return len(a)
}
func (a Sorter) Less(i, j int) bool {
	return a[i].Compare(a[j]) < 0
}
func (a Sorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Interface to sort dates (Compare up to milliseconds). Example:
//    sort.Sort(date.Sorter(sliceOfDates))
type SorterTime []T

func (a SorterTime) Len() int {
	return len(a)
}
func (a SorterTime) Less(i, j int) bool {
	return a[i].CompareTime(a[j]) < 0
}
func (a SorterTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Returns a string that represents the date.
// 'template' is a kind 'printf' with next sustitution variables:
//    %d  Day in number 06 -> 6.
//    %D  Day with tow digits 06 -> 06.
//    %m  Month in number 03 -> 3.
//    %M  Month with two digits 03 -> 03.
//    %y  Year with two digits 2010 -> 10.
//    %Y  Year with four digits 2010 -> 2010.
//    %t  Time without milliseconds -> 15:03:55
//    %T  Time with milliseconds -> 15:03:55.345
//    %%  The sign '%'.
func (d T) Format(tmp string) string {
	tmp = strings.ReplaceAll(tmp, "%d", time.Time(d).Format("2"))
	tmp = strings.ReplaceAll(tmp, "%D", time.Time(d).Format("02"))
	tmp = strings.ReplaceAll(tmp, "%m", time.Time(d).Format("1"))
	tmp = strings.ReplaceAll(tmp, "%M", time.Time(d).Format("01"))
	tmp = strings.ReplaceAll(tmp, "%y", time.Time(d).Format("06"))
	tmp = strings.ReplaceAll(tmp, "%Y", time.Time(d).Format("2006"))
	tmp = strings.ReplaceAll(tmp, "%t", time.Time(d).Format("15:04:05"))
	tmp = strings.ReplaceAll(tmp, "%T", time.Time(d).Format("15:04:05.000"))
	tmp = strings.ReplaceAll(tmp, "%%", "%")
	return tmp
}

// Returns a string in format 'yyyymmdd'.
func (d T) String() string {
	return d.Format("%Y%M%D")
}

// Returns a string in format 'dd/mm/yyyy'.
func (d T) ToIso() string {
	return d.Format("%D/%M/%Y")
}

// Returns a string in format 'mm/dd/yyyy'.
func (d T) ToEn() string {
	return d.Format("%M/%D/%Y")
}
