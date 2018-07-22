package main 
import (
        "bufio"
        "bytes"
        "path/filepath"
        "io/ioutil"
        "log"
        "os"
        "testing"
)

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
        dir_prefix := "./assets/"
        test_files, _ := ioutil.ReadDir("./assets") 
        for _, file := range test_files {
                fname := dir_prefix + file.Name()
                bkp_file := CreateBackupName(fname, "bak")
                err := BackupFile(fname, "bak")
                if err != nil {
                        t.Log(err)
                        t.Fail()
                }

                if !EqualFiles(fname, bkp_file) {
                        t.Logf("%q != %q", fname, bkp_file)
                        t.Fail()
                }
                
        }
}

func TestRestoreFile(t *testing.T) {
        dir := "./assets"
        dir_prefix := "./assets/"
        test_files, _ := ioutil.ReadDir(dir) 
        for _, file := range test_files {
                fname := dir_prefix + file.Name()
                if filepath.Ext(fname) == ".txt" {
                        bkp_file := CreateBackupName(fname, "bak")
                        err := RestoreFile(fname, "bak")
                        if err != nil {
                                t.Log(err)
                                t.Fail()
                        }

                        if !EqualFiles(fname, bkp_file) {
                                t.Logf("%q != %q", fname, bkp_file)
                                t.Fail()
                        }
                }
        }
}

func TestRestoreDir(t *testing.T) {
        dir := "./assets"
        dir_prefix := "./assets/"
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
                                t.Logf("%q != %q", fname, bkp_file)
                                t.Fail()
                        }                        
                }
        }

        for _, file := range test_files {
                fname := dir_prefix + file.Name()
                if filepath.Ext(fname) == ".txt" {
                        bkp_file := CreateBackupName(fname, "bak")
                        os.Remove(bkp_file)
                }
        }
}
