package shtool

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	packages       = []*Package{}
	activePackages = []*Package{}
)

type Package struct {
	Name      string
	OS        string
	Arch      string
	Commands  []string
	Download  string
	Checksum  string
	Installer func(path string) error
	Activator func(path string) error
}

func lookup(name string) (*Package, error) {
	availablePackages := []*Package{}
	for _, pkg := range packages {
		if pkg.OS == runtime.GOOS && pkg.Arch == runtime.GOARCH {
			availablePackages = append(availablePackages, pkg)
		}
	}

	for _, pkg := range availablePackages {
		for _, command := range pkg.Commands {
			if command == name {
				return pkg, nil
			}
		}
	}
	return nil, errors.New("error: command '" + name + "' not available")
}

func (pkg *Package) isActive() bool {
	for _, activePackage := range activePackages {
		if activePackage == pkg {
			return true
		}
	}
	return false
}

func (pkg *Package) isInstalled(path string) bool {
	_, err := os.Stat(filepath.Join(path, pkg.Name, ".installed"))
	return err == nil
}

func (pkg *Package) activate(path string) error {
	return pkg.Activator(filepath.Join(path, pkg.Name))
}

func (pkg *Package) install(path string) error {
	info("installing package " + pkg.Name)

	pkgPath := filepath.Join(path, pkg.Name)

	err := os.RemoveAll(pkgPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(pkgPath, os.ModePerm)
	if err != nil {
		return err
	}

	info("downloading from " + pkg.Download)
	archiveName := filepath.Base(pkg.Download)
	archivePath := filepath.Join(pkgPath, archiveName)
	err = download(archivePath, pkg.Download)
	if err != nil {
		return err
	}

	archive, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	h := sha256.New()
	if _, err := io.Copy(h, archive); err != nil {
		return err
	}
	checksum := hex.EncodeToString(h.Sum(nil))
	archive.Close()
	if pkg.Checksum != checksum {
		os.Remove(archivePath)
		return errors.New("expected sha256sum " + pkg.Checksum + " but got " + checksum)
	}
	info("sha256sum verified " + pkg.Checksum)

	info("installing " + filepath.Join(pkg.Name, archiveName))
	err = pkg.Installer(archivePath)
	if err != nil {
		return err
	}
	err = os.Remove(archivePath)
	if err != nil {
		return err
	}
	installedFile, err := os.Create(filepath.Join(pkgPath, ".installed"))
	if err != nil {
		return err
	}
	defer installedFile.Close()
	_, err = installedFile.WriteString(time.Now().String())
	if err != nil {
		return err
	}
	info("package '" + pkg.Name + "' installed successfully")

	return nil
}
