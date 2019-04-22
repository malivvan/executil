package executil

import (
	"github.com/mholt/archiver"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"runtime"
)

type Logger func(...interface{})

var logger Logger = nil

func SetLogger(l Logger) {
	logger = l
}

func info(i ...interface{}) {
	if logger != nil {
		logger(i...)
	}
}

func PathActivator(path string) error {
	seperator := ":"
	if runtime.GOOS == "windows" {
		seperator = ";"
	}
	return os.Setenv("PATH", os.Getenv("PATH")+seperator+path)
}

func ArchiveInstaller(prefix string) func(path string) error {
	return func(path string) error {
		dir := filepath.Dir(path)
		err := archiver.Unarchive(path, dir)
		if err != nil {
			return err
		}

		files, err := ioutil.ReadDir(filepath.Join(dir, prefix))
		if err != nil {
			return err
		}
		for _, f := range files {
			err = os.Rename(filepath.Join(dir, prefix, f.Name()), filepath.Join(dir, f.Name()))
			if err != nil {
				return err
			}
		}

		segs := strings.Split(prefix, string(os.PathSeparator))
		if len(segs) > 0 {
			err := os.Remove(filepath.Join(dir, segs[0]))
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func download(path string, url string) error {

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
