// Copyright 04-Apr-2020 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

// Intialization and system utilities.
package sys

import (
	"fmt"
	"github.com/dedeme/golib/file"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"time"
)

var homeV string

// Creates home directory and intialize random generator.
//    home: Aplication subdirectory in ~/.dmGoApp
func Initialize(home string) {
	rand.Seed(time.Now().UTC().UnixNano())
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("It is not possible get home dir")
	}
	p := path.Join(homeDir, ".dmGoApp", home)
	file.Mkdirs(p)
	homeV = p
}

// Returns application directory
func Home() string {
	return homeV
}

// Run a system command and returns is stdout and stderr responses.
// Example:
//    sys.Cmd("ls", "-l", "/")
func Cmd(c string, pars ...string) (rpOut []byte, rpError []byte) {
	cmd := exec.Command(c, pars...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		rpError = []byte(fmt.Sprintf("Pipe of stdout failed in cmd '%v'", c))
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		rpError = []byte(fmt.Sprintf("Pipe of stderror failed in cmd '%v'", c))
	}
	if err := cmd.Start(); err != nil {
		rpError = []byte(fmt.Sprintf("Fail starting cmd '%v'", c))
	}

	rpOut, err = ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
	}
	rpError, err = ioutil.ReadAll(stderr)
	if err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		if len(rpError) == 0 {
			rpError = []byte(fmt.Sprintf("'%v' in cmd '%v'", err, c))
		}
	}
	return
}

// Stop the current thread 'milliseconds' milliseconds.
func Sleep(milliseconds int) {
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}
