// Copyright 30-Mar-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Utilities for HTML conections between client - server.
package cgi

import (
	"fmt"
	"github.com/dedeme/golib/cryp"
	"github.com/dedeme/golib/file"
	"github.com/dedeme/golib/json"
	"github.com/dedeme/golib/log"
	"math/rand"
	"path"
	"time"
)

const (
	// Standad length of passwords
	Klen          = 300
	tNoExpiration = 2592000 // seconds == 30 days
	demeKey       = "nkXliX8lg2kTuQSS/OoLXCk8eS4Fwmc+N7l6TTNgzM1vdKewO0cjok51vcdl" +
		"OKVXyPu83xYhX6mDeDyzapxL3dIZuzwyemVw+uCNCZ01WDw82oninzp88Hef" +
		"bn3pPnSMqEaP2bOdX+8yEe6sGkc3IO3e38+CqSOyDBxHCqfrZT2Sqn6SHWhR" +
		"KqpJp4K96QqtVjmXwhVcST9l+u1XUPL6K9HQfEEGMGcToMGUrzNQxCzlg2g+" +
		"Hg55i7iiKbA0ogENhEIFjMG+wmFDNzgjvDnNYOaPTQ7l4C8aaPsEfl3sugiw"
	noSessionKey = "nosession"
)

var homeV string
var tExpirationV int
var fkey = cryp.Key(demeKey, len(demeKey)) // File encryption key
var key string                             // Communication key

// Initializes a new interface of commnications.
//    home       : Aboslute path of application directory. For example:
//                    "/peter/wwwcgi/dmcgi/JsMon"
//                    or
//                    "/home/deme/.dmCApp/JsMon" ).
//    tExpiration: Time in seconds.
func Initialize(home string, tExpiration int) {
	rand.Seed(time.Now().UTC().UnixNano())
	homeV = home
	tExpirationV = tExpiration

	if !file.Exists(home) {
		file.Mkdirs(home)
	}

	fusers := path.Join(home, "users.db")
	if !file.Exists(fusers) {
		writeUsers([]*userT{})
		putUser("admin", demeKey, "0")
		writeSessions([]*sessionT{})
	}
}

// User ------------------------------------------------------------------------

type userT struct {
	id    string
	pass  string
	level string
}

func uToJson(us []*userT) json.T {
	var tmps []json.T
	for _, u := range us {
		jss := []json.T{json.Ws(u.id), json.Ws(u.pass), json.Ws(u.level)}
		tmps = append(tmps, json.Wa(jss))
	}
	return json.Wa(tmps)
}
func uFromJson(js json.T) (us []*userT) {
	jss := js.Ra()
	for _, ujs := range jss {
		var ujss []json.T
		ujss = ujs.Ra()
		var id, pass, level string
		id = ujss[0].Rs()
		pass = ujss[1].Rs()
		level = ujss[2].Rs()
		us = append(us, &userT{id, pass, level})
	}
	return
}

func writeUsers(us []*userT) {
	js := uToJson(us)
	file.WriteAll(path.Join(homeV, "users.db"), cryp.Cryp(fkey, string(js)))
}

func readUsers() []*userT {
	js := cryp.Decryp(fkey, file.ReadAll(path.Join(homeV, "users.db")))
	r := uFromJson(json.FromString(js))
	return r
}

func putUser(id, pass, level string) {
	pass = cryp.Key(pass, Klen)
	users := readUsers()

	var r *userT
	for _, u := range users {
		if u.id == id {
			r = u
			break
		}
	}
	if r == nil {
		users = append(users, &userT{id, pass, level})
	} else {
		r.pass = pass
		r.level = level
	}
	writeUsers(users)
}

// If check fails, returns "". Otherwise it returns user level.
func checkUser(id, pass string) string {
	pass = cryp.Key(pass, Klen)
	users := readUsers()
	for _, u := range users {
		if u.id == id && u.pass == pass {
			return u.level
		}
	}
	return ""
}

// Session ---------------------------------------------------------------------

type sessionT struct {
	id     string
	comKey string // Communication key
	conKey string // Connection key
	user   string // User id
	level  string
	time   int64 // time.Unix
	lapse  int
}

func sToJson(ss []*sessionT) json.T {
	var tmps []json.T
	for _, s := range ss {
		jss := []json.T{
			json.Ws(s.id), json.Ws(s.comKey), json.Ws(s.conKey),
			json.Ws(s.user), json.Ws(s.level),
			json.Wl(s.time), json.Wi(s.lapse),
		}
		tmps = append(tmps, json.Wa(jss))
	}
	return json.Wa(tmps)
}
func sFromJson(js json.T) (ss []*sessionT, err error) {
	jss := js.Ra()
	for _, sjs := range jss {
		var sjss []json.T
		sjss = sjs.Ra()
		var id, comKey, conKey, user, level string
		var time int64
		var lapse int
		id = sjss[0].Rs()
		comKey = sjss[1].Rs()
		conKey = sjss[2].Rs()
		user = sjss[3].Rs()
		level = sjss[4].Rs()
		time = sjss[5].Rl()
		lapse = sjss[6].Ri()
		ss = append(ss, &sessionT{id, comKey, conKey, user, level, time, lapse})
	}
	return
}

// Returns 'false' if 's' is out of date.
func (s *sessionT) update() bool {
	now := time.Now().Unix()
	if now > s.time+int64(s.lapse) {
		return false
	}
	s.time = now + int64(s.lapse)
	return true
}

func writeSessions(ss []*sessionT) {
	js := sToJson(ss)
	file.WriteAll(path.Join(homeV, "sessions.db"), cryp.Cryp(fkey, string(js)))
}

func readSessions() []*sessionT {
	js := cryp.Decryp(fkey, file.ReadAll(path.Join(homeV, "sessions.db")))
	r, err := sFromJson(json.FromString(js))
	if err != nil {
		log.Fatal(err)
	}
	return r
}

// Adds session and purge sessions.
func addSession(sessionId, comKey, conKey, user, level string, lapse int) {
	now := time.Now().Unix()
	ss := readSessions()
	var newSs []*sessionT
	for _, s := range ss {
		if now <= s.time+int64(s.lapse) {
			newSs = append(newSs, s)
		}
	}
	newSs = append(newSs,
		&sessionT{sessionId, comKey, conKey, user, level, now, lapse})
	writeSessions(newSs)
}

// Replace a session with a new date and a new connection key
func replaceSession(mdss *sessionT) {
	ss := readSessions()
	var newSs []*sessionT
	for _, s := range ss {
		if s.id != mdss.id {
			newSs = append(newSs, s)
		}
	}
	newSs = append(newSs, mdss)
	writeSessions(newSs)
}

// Public interface ------------------------------------------------------------

// Root application directory.
func Home() string {
	return homeV
}

// Sends to client 'communicationKey', 'userId' and 'userLevel'. If conection
// fails every one is "".
//    sessionId: Session identifier.
//    return   : {key: String, conKey:String, user: String, level: String}.
func Connect(sessionId string) string {
	var r *sessionT
	for _, s := range readSessions() {
		if s.id == sessionId && s.update() {
			r = s
			break
		}
	}

	comKey := ""
	conKey := ""
	user := ""
	level := ""
	if r != nil {
		comKey = r.comKey
		user = r.user
		level = r.level
		conKey = cryp.GenK(Klen)
		r.conKey = conKey
		replaceSession(r)
	}
	return Rp(sessionId, map[string]json.T{
		"key":    json.Ws(comKey),
		"conKey": json.Ws(conKey),
		"user":   json.Ws(user),
		"level":  json.Ws(level),
	})
}

// Sends to client 'sessionId', 'communicationKey' and 'userLevel'. If
// conection fails every one is "".
//   user          : User id.
//    key           : User password.
//    withExpiration: If is set to false, session will expire after 30 days.
//    return        : {sessionId: String, key: String, conKey: String,
//                     level: String}.
func Authentication(key, user, pass string, withExpiration bool) string {
	sessionId := ""
	comKey := ""
	conKey := ""
	level := checkUser(user, pass)
	if level != "" {
		sessionId = cryp.GenK(Klen)
		comKey = cryp.GenK(Klen)
		conKey = cryp.GenK(Klen)

		lapse := tNoExpiration
		if withExpiration {
			lapse = tExpirationV
		}
		addSession(sessionId, comKey, conKey, user, level, lapse)
	}

	return Rp(key, map[string]json.T{
		"sessionId": json.Ws(sessionId),
		"key":       json.Ws(comKey),
		"conKey":    json.Ws(conKey),
		"level":     json.Ws(level),
	})
}

// Returns the session communication key.
//		ssId  : Session identifier.
//		conKey: Connection key. If its value is "", this parameter is not used.
func GetComKey(ssId, conKey string) (comKey string, ok bool) {
	ss := readSessions()
	for _, s := range ss {
		if s.id == ssId && (conKey == "" || conKey == s.conKey) && s.update() {
			comKey = s.comKey
			ok = true
			return
		}
	}
	return
}

// Changes user password.
//		ck    : Communication key
//    user  : User name to change password.
//    old   : Old password.
//    new   : New password.
//    return: A boolean field {ok:true|false}, sets to true if operation
//            succeeded. A fail can come up if 'user' authentication fails.
func ChangePass(ck, user, old, new string) (rp string) {
	rp = Rp(ck, map[string]json.T{"ok": json.Wb(false)})

	us := readUsers()
	var u *userT
	for _, u0 := range us {
		if u0.id == user {
			u = u0
			break
		}
	}
	if u == nil {
		return
	}

	old2 := cryp.Key(old, Klen)
	if old2 != u.pass {
		return
	}

	u.pass = cryp.Key(new, Klen)
	writeUsers(us)
	rp = Rp(ck, map[string]json.T{"ok": json.Wb(true)})
	return
}

// Deletes 'sessionId' and returns an empty response.
func DelSession(ck string, sessionId string) string {
	ss := readSessions()
	var newss []*sessionT
	for _, s := range ss {
		if s.id != sessionId {
			newss = append(newss, s)
		}
	}
	writeSessions(newss)
	return RpEmpty(ck)
}

// Messages --------------------------------------------------------------------

// Returns a response to send to client.
//	 ck: Communication key.
//	 rp: Response.
func Rp(ck string, rp map[string]json.T) string {
	js := json.Wo(rp)
	return cryp.Cryp(ck, string(js))
}

// Returns an empty response.
//	 ck: Communication key.
func RpEmpty(ck string) string {
	return Rp(ck, map[string]json.T{})
}

// Returns a message with an only field "error" with value 'msg'.
//	 ck: Communication key.
func RpError(ck, msg string) string {
	return Rp(ck, map[string]json.T{"error": json.Ws(msg)})
}

// Returns a message with an only field "expired" with value 'true',
// codified with the key 'noSessionKey' ("nosession")
func RpExpired() string {
	return Rp(noSessionKey, map[string]json.T{"expired": json.Wb(true)})
}

// Requests --------------------------------------------------------------------

// Reads a bool value
func RqBool(rq map[string]json.T, key string) (v bool) {
	js, ok := rq[key]
	if !ok {
		panic(fmt.Sprintf("Key '%v' not found in request", key))
	}
	v = js.Rb()
	return
}

// Reads a int value
func RqInt(rq map[string]json.T, key string) (v int) {
	js, ok := rq[key]
	if !ok {
		panic(fmt.Sprintf("Key '%v' not found in request", key))
	}
	v = js.Ri()
	return
}

// Reads a int64 value
func RqLong(rq map[string]json.T, key string) (v int64) {
	js, ok := rq[key]
	if !ok {
		panic(fmt.Sprintf("Key '%v' not found in request", key))
	}
	v = js.Rl()
	return
}

// Reads a float32 value
func RqFloat(rq map[string]json.T, key string) (v float32) {
	js, ok := rq[key]
	if !ok {
		panic(fmt.Sprintf("Key '%v' not found in request", key))
	}
	v = js.Rf()
	return
}

// Reads a float64 value
func RqDouble(rq map[string]json.T, key string) (v float64) {
	js, ok := rq[key]
	if !ok {
		panic(fmt.Sprintf("Key '%v' not found in request", key))
	}
	v = js.Rd()
	return
}

// Reads a string value
func RqString(rq map[string]json.T, key string) (v string) {
	js, ok := rq[key]
	if !ok {
		panic(fmt.Sprintf("Key '%v' not found in request", key))
	}
	v = js.Rs()
	return
}
