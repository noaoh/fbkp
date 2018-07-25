package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

type Env struct{ directory, dir_prefix, ext string }
type testCase struct{ filename, real_filename, backup_filename, ext string }

var testCases []testCase
var benchCases []testCase

func init() {
	TestEnv := Env{
		directory:  "./assets/test",
		dir_prefix: "./assets/test/",
		ext:        "bak",
	}

	BenchEnv := Env{
		directory:  "./assets/bench",
		dir_prefix: "./assets/bench/",
		ext:        "bak",
	}

	test_files, _ := ioutil.ReadDir(TestEnv.directory)
	for _, test_file := range test_files {
		fname := test_file.Name()
		rname := TestEnv.dir_prefix + fname
		bname := CreateBackupName(rname, TestEnv.ext)
		bak_ext := TestEnv.ext

		testCases = append(testCases, testCase{
			filename:        fname,
			real_filename:   rname,
			backup_filename: bname,
			ext:             bak_ext,
		})
	}

	bench_files, _ := ioutil.ReadDir(BenchEnv.directory)
	for _, bench_file := range bench_files {
		fname := bench_file.Name()
		rname := BenchEnv.dir_prefix + fname
		bname := CreateBackupName(rname, BenchEnv.ext)
		bak_ext := BenchEnv.ext

		benchCases = append(benchCases, testCase{
			filename:        fname,
			real_filename:   rname,
			backup_filename: bname,
			ext:             bak_ext,
		})
	}
}

func EqualFiles(file1, file2 string) bool {
	sf, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}

	df, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}

	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)
	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return false
		}
	}
	return true
}

func TestBackupFile(t *testing.T) {
	for _, test := range testCases {
		err := BackupFile(test.real_filename, test.ext)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if !EqualFiles(test.real_filename, test.backup_filename) {
			t.Logf("Backup file %q is not equivalent to %q", test.backup_filename, test.real_filename)
			t.Fail()
		}
	}
}

func TestRestoreFile(t *testing.T) {
	for _, test := range testCases {
		err := RestoreFile(test.backup_filename)
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		if !EqualFiles(test.real_filename, test.backup_filename) {
			t.Logf("Restored file %q is not equivalent to %q", test.real_filename, test.backup_filename)
			t.Fail()
		}
	}
}

func TestBackupDir(t *testing.T) {
	dir := "./assets/test"
	dir_prefix := "./assets/test/"
	err := BackupDir(dir, "bak")
	if err != nil {
		t.Log(err)
	}

	test_files, _ := ioutil.ReadDir(dir)
	for _, file := range test_files {
		fname := dir_prefix + file.Name()
		if filepath.Ext(fname) == ".txt" {
			bkp_file := CreateBackupName(fname, "bak")
			if !EqualFiles(fname, bkp_file) {
                                t.Logf("Backup file %q is not equivalent to %q", bkp_file, fname)
				t.Fail()
			}
		}
	}
}

func TestRestoreDir(t *testing.T) {
	dir := "./assets/test"
	dir_prefix := "./assets/test/"
	err := RestoreDir(dir, "bak")
	if err != nil {
		t.Log(err)
	}

	test_files, _ := ioutil.ReadDir(dir)
	for _, file := range test_files {
		fname := dir_prefix + file.Name()
		if filepath.Ext(fname) == ".txt" {
			bkp_file := CreateBackupName(fname, "bak")
			if !EqualFiles(fname, bkp_file) {
                                t.Logf("Restored file %q is not equivalent to %q", fname, bkp_file)
				t.Fail()
			}
		}
	}
}

func TestMain(m *testing.M) {
	runTests := m.Run()

        for _, file := range testCases {
                os.Remove(file.backup_filename)
        }
	os.Exit(runTests)
}
