// Copyright 31-Aug-2017 ºDeme
// GNU General Public License - V3 <http://www.gnu.org/licenses/>

package _tests

import (
	"github.com/dedeme/golib/file"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

func TestFile(t *testing.T) {
	u, _ := user.Current()
	dir := filepath.Join(u.HomeDir, ".dmGoApp", "dmGoLib")
	tmp := filepath.Join(dir, "tmp.txt")
	tmp1 := filepath.Join(dir, "tmp1.txt")
	tmp2 := filepath.Join(dir, "tmp2.txt")
  cpdir := filepath.Join(dir, "cpdir")
  cptmp := filepath.Join(cpdir, "tmp.txt")
  cptmpx := filepath.Join(cpdir, "tmpx.txt")
  cpdir2 := filepath.Join(dir, "cpdir2")
  cptmp2 := filepath.Join(cpdir2, "cpdir", "tmp.txt")
  cptmpx2 := filepath.Join(cpdir2, "cpdir", "tmpx.txt")

	file.Mkdirs(dir)

	ftmp := file.OpenWrite(tmp)
	file.Write(ftmp, "Una\n")
	file.Write(ftmp, "\n")
	file.Write(ftmp, "Dos...\ny Tres")
	ftmp.Close()

  file.Mkdir(cpdir)
  file.Copy(tmp, cpdir)
  file.Copy(tmp, cptmpx)
  if file.ReadAll(cptmp) != file.ReadAll(cptmpx) {
    t.Fatal("Fail copying file")
  }
  file.Mkdir(cpdir2)
  file.Copy(cpdir, cpdir2)
  if file.ReadAll(cptmp2) != file.ReadAll(cptmpx2) {
    t.Fatal("Fail copying directory")
  }
  file.Remove(cpdir)
  file.Remove(cpdir2)

	ftmp = file.OpenAppend(tmp)
	file.WriteBin(ftmp, []byte("\nY un añadido"))
	ftmp.Close()

	tx := file.ReadAll(tmp)
	t.Log(tx)
	if tx != "Una\n\nDos...\ny Tres\nY un añadido" {
		t.Fatal("LineReader gives:\n" + tx)
	}

	ftmp1 := file.OpenWrite(tmp1)
	file.Write(ftmp1, tx)
	ftmp1.Close()

	file.WriteAll(tmp2, tx)

	tx1 := file.ReadAll(tmp1)
	if tx != tx1 {
		t.Fatal("LineReader gives:\n" + tx1)
	}
	ftmp1.Close()

	tx2 := ""
	file.Lines(tmp2, func(l string) bool {
		tx2 += l + "\n"
		return false
	})

	if tx != strings.TrimSpace(tx2) {
		t.Fatal("LineReader gives:\n" + tx2)
	}

	file.Zip(dir, filepath.Join(dir, "dir.zip"))
	file.Unzip(filepath.Join(dir, "dir.zip"), dir)

	zipDir := filepath.Join(dir, "dmGoLib")
	if !file.Exists(zipDir) {
		t.Fatal("Decompressed directory 'dmGoLib' does not exist")
	}
	if len(file.List(zipDir)) != 4 {
		t.Fatal("Number of files in decompressed directory 'dmGoLib' is not 4")
	}

	file.Remove(dir)
}
