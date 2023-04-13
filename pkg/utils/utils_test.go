package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDataDir = "test_data"

func getTestDir(subDir string) string {
	pwd, _ := os.Getwd()
	testDir := filepath.Join(pwd, testDataDir)
	testDir = filepath.Join(testDir, subDir)

	return testDir
}

func TestParseRepoNameFromGitUrl(t *testing.T) {
	assert.Equal(t, ParseRepoNameFromGitUrl("test"), "test", "ParseRepoNameFromGitUrl failed.")
	assert.Equal(t, ParseRepoNameFromGitUrl("test.git"), "test", "ParseRepoNameFromGitUrl failed.")
	assert.Equal(t, ParseRepoNameFromGitUrl("test.git.git"), "test.git", "ParseRepoNameFromGitUrl failed.")
	assert.Equal(t, ParseRepoNameFromGitUrl("https://test.git"), "test", "ParseRepoNameFromGitUrl failed.")
	assert.Equal(t, ParseRepoNameFromGitUrl("https://test.git.git"), "test.git", "ParseRepoNameFromGitUrl failed.")
	assert.Equal(t, ParseRepoNameFromGitUrl("httfsdafps://test.git.git"), "test.git", "ParseRepoNameFromGitUrl failed.")
}

type TestPath struct {
	FilePath string
}

func (tp *TestPath) TestStore() error {
	return StoreToFile(tp.FilePath, "test")
}

func TestCreateFileIfNotExist(t *testing.T) {
	test_path := getTestDir("test_exist.txt")
	isExist, _ := Exists(test_path)
	assert.Equal(t, isExist, false)

	tp := TestPath{
		FilePath: test_path,
	}
	err := CreateFileIfNotExist(tp.FilePath, tp.TestStore)
	assert.Equal(t, err, nil)

	isExist, _ = Exists(test_path)
	assert.Equal(t, isExist, true)

	_ = os.Remove(test_path)
	isExist, _ = Exists(test_path)
	assert.Equal(t, isExist, false)
}

func TestHashDir(t *testing.T) {
	test_path := filepath.Join(getTestDir("test_hash"), "test_hash.txt")
	tp := TestPath{
		FilePath: test_path,
	}

	_ = CreateFileIfNotExist(tp.FilePath, tp.TestStore)
	res, err := HashDir(filepath.Dir(tp.FilePath))
	assert.Equal(t, err, nil)
	assert.Equal(t, res, "n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg=")
}

func TestTarDir(t *testing.T) {
	testDir := getTestDir("test_tar")
	tarPath := filepath.Join(testDir, "test.tar")

	_, err := os.Stat(tarPath)
	if !os.IsNotExist(err) {
		os.Remove(tarPath)
	}

	err = TarDir(filepath.Join(testDir, "test_src"), tarPath)
	assert.Equal(t, err, nil)

	_, err = os.Stat(tarPath)
	assert.Equal(t, err, nil)
	os.Remove(tarPath)
}

func TestGetAbsKpmHome(t *testing.T) {
	os.Setenv("KPM_HOME", ".")
	got, err := GetAbsKpmHome()
	expect, _ := filepath.Abs(".")
	assert.Equal(t, err, nil)
	assert.Equal(t, got, expect)

	os.Setenv("KPM_HOME", "test_subdir")
	got, err = GetAbsKpmHome()
	assert.Equal(t, got, filepath.Join(expect, "test_subdir"))
	assert.Equal(t, err, nil)
}

func TestCreateSymbolLink(t *testing.T) {
	testDir := getTestDir("test_link")
	need_linked := filepath.Join(testDir, "need_be_linked_v1")
	linkPath := filepath.Join(testDir, "linked")

	_ = os.Remove(linkPath)
	err := CreateSymlink(need_linked, linkPath)

	linkTarget, _ := os.Readlink(linkPath)
	assert.Equal(t, err, nil)
	assert.Equal(t, linkTarget, need_linked)
	_ = os.Remove(linkPath)
}