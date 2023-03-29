// Copyright 2022 The KCL Authors. All rights reserved.
package modfile

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"kusionstack.io/kpm/pkg/opt"
	"kusionstack.io/kpm/pkg/reporter"
)

const (
	MOD_FILE      = "kcl.mod"
	MOD_LOCK_FILE = "kcl.mod.lock"
)

// 'Package' is the kcl package section of 'kcl.mod'.
type Package struct {
	Name    string `toml:"name,omitempty"`    // kcl package name
	Edition string `toml:"edition,omitempty"` // kclvm version
	Version string `toml:"version,omitempty"` // kcl package version
}

// 'ModFile' is kcl package file 'kcl.mod'.
type ModFile struct {
	HomePath string  `toml:"-"`
	Pkg      Package `toml:"package,omitempty"`
	Dependencies
}

// 'Dependencies' is dependencies section of 'kcl.mod'.
type Dependencies struct {
	Deps map[string]Dependency `toml:"dependencies,omitempty"`
}

type Dependency struct {
	Name    string `toml:"name,omitempty"`
	Version string `toml:"version,omitempty"`
	Sum     string `toml:"sum,omitempty"`
	Source
}

type Source struct {
	*Git
}

type Git struct {
	Url    string `toml:"url,omitempty"`
	Branch string `toml:"branch,omitempty"`
	Commit string `toml:"commit,omitempty"`
	Tag    string `toml:"tag,omitempty"`
}

// ModFileExists returns whether a 'kcl.mod' file exists in the path.
func ModFileExists(path string) (bool, error) {
	return exists(filepath.Join(path, MOD_FILE))
}

// ModLockFileExists returns whether a 'kcl.mod.lock' file exists in the path.
func ModLockFileExists(path string) (bool, error) {
	return exists(filepath.Join(path, MOD_LOCK_FILE))
}

// LoadModFile load the contents of the 'kcl.mod' file in the path.
func LoadModFile(homePath string) (*ModFile, error) {
	modFile := new(ModFile)
	err := loadFile(homePath, MOD_FILE, modFile)
	if err != nil {
		return nil, err
	}

	modFile.HomePath = filepath.Join(homePath, MOD_FILE)

	if modFile.Dependencies.Deps == nil {
		modFile.Dependencies.Deps = make(map[string]Dependency)
	}

	return modFile, nil
}

// Write the contents of 'ModFile' to 'kcl.mod' file
func (mfile *ModFile) Store() error {
	fullPath := filepath.Join(mfile.HomePath, MOD_FILE)
	return storeToFile(fullPath, mfile)
}

// Write the contents of dependencies 'ModFile' to 'kcl.mod.lock' file
func (mfile *ModFile) StoreLockFile() error {
	fullPath := filepath.Join(mfile.HomePath, MOD_LOCK_FILE)
	return storeToFile(fullPath, mfile.Dependencies)
}

func (mfile *ModFile) GetModFilePath() string {
	return filepath.Join(mfile.HomePath, MOD_FILE)
}

func (mfile *ModFile) GetModLockFilePath() string {
	return filepath.Join(mfile.HomePath, MOD_LOCK_FILE)
}

const defaultVerion = "0.0.1"
const defaultEdition = "0.0.1"

func NewModFile(opts *opt.InitOptions) *ModFile {
	return &ModFile{
		HomePath: opts.InitPath,
		Pkg: Package{
			Name:    opts.Name,
			Version: defaultVerion,
			Edition: defaultEdition,
		},
	}
}

func storeToFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		reporter.ExitWithReport("kpm: failed to create file: ", filePath, err)
		return err
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(data); err != nil {
		reporter.ExitWithReport("kpm: failed to encode TOML:", err)
		return err
	}
	return nil
}

func loadFile(homePath string, fileName string, file interface{}) error {
	readFile, err := os.OpenFile(filepath.Join(homePath, fileName), os.O_RDWR, 0644)
	if err != nil {
		reporter.Report("kpm: failed to load", fileName)
		return err
	}
	defer readFile.Close()

	_, err = toml.NewDecoder(readFile).Decode(file)
	if err != nil {
		return err
	}

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}