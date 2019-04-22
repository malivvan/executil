package executil

import (
	"github.com/mitchellh/go-homedir"
	"os/exec"
	"path/filepath"
	"runtime"
)

var executilPath = func() string {
	dir, err := homedir.Dir()
	if err == nil {
		dir, err = homedir.Expand(dir)
		if err == nil {
			return filepath.Join(dir, ".executil")
		}
	}
	return ""
}()

// Register a new Package. If a command
func Register(pkg Package) {
	if pkg.OS == runtime.GOOS && pkg.Arch == runtime.GOARCH {
		for _, command := range pkg.Commands {
			existingPkg, _ := lookup(command)
			if existingPkg != nil {
				panic("error: command '" + command + "' already available in package '" + existingPkg.Name + "!")
			}
		}
	}
	packages = append(packages, &pkg)
}

func Ensure(name string) error {

	if IsAvailable(name) {
		return nil
	}

	if !IsInstalled(name) {
		pkg, err := lookup(name)
		if err != nil {
			return err
		}
		err = pkg.install(executilPath)
		if err != nil {
			return err
		}
	}

	err := Activate(name)
	if err != nil {
		return err
	}

	return nil
}

// IsAvailable returns true if the command is available within
// the PATH variable.
func IsAvailable(command string) bool {
	_, err := exec.LookPath(command)
	if err == nil {
		return true
	}
	return IsActive(command)
}

func IsActive(command string) bool {
	pkg, err := lookup(command)
	if err != nil {
		return false
	}
	return pkg.isActive()
}

// Activate makes a command available.
func Activate(command string) error {
	pkg, err := lookup(command)
	if err != nil {
		return err
	}
	if pkg.isActive() {
		return nil
	}
	activePackages = append(activePackages, pkg)
	return pkg.activate(executilPath)
}

func Install(name string) error {

	pkg, err := lookup(name)
	if err != nil {
		return err
	}
	err = pkg.install(executilPath)
	if err != nil {
		return err
	}
	return nil
}

func IsInstalled(name string) bool {
	pkg, err := lookup(name)
	if err != nil {
		return false
	}
	return pkg.isInstalled(executilPath)
}
