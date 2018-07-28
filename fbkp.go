package main

import (
	"io/ioutil"
        "log"
	"os"
	"path/filepath"
	"strings"
)

func CreateBackupName(filename, ext string) string {
	return filename + "." + ext
}

func CreateOriginalName(filename string) string {
        return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func CopyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, buf, 0755)
	if err != nil {
		return err
	}

	return err
}

func BackupFile(path, ext string) error {
        info, err := os.Stat(path)
        if err != nil {
                return err
        }

        if !info.IsDir() {
                bkp_path := CreateBackupName(path, ext)
                err := CopyFileContents(path, bkp_path)
                if err != nil {
                        return err
                }
        }
        return err
}

func RestoreFile(path, ext string) error {
        info, err := os.Stat(path)
        if !info.IsDir() {
                og_path := CreateOriginalName(path)
                err := CopyFileContents(path, og_path)
                if err != nil {
                        return err
                }
        }
	return err
}

func BackupDir(dir string, ext string, verbose bool, recursive bool) error {
        real_ext := "." + ext
        err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                        log.Printf("skipping a directory without errors: %q\n", info.Name())
                        log.Printf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
                        return err
                }

                fname := info.Name()
                if !info.IsDir() && filepath.Ext(fname) != real_ext {
                        bkp_path := CreateBackupName(path, ext)
                        err := CopyFileContents(path, bkp_path)
                        if err != nil {
                                return err
                        }
                        log.Printf("%q -> %q\n", path, bkp_path)
                } else if !recursive && path != dir {
                        return filepath.SkipDir
                }
                return err
        })
        return err
}

func RestoreDir(dir string, ext string, verbose bool, recursive bool) error {
        real_ext := "." + ext
        err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
                if err != nil {
                        log.Printf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
                        return err
                }

                fname := info.Name()
                if !info.IsDir() && filepath.Ext(fname) == real_ext {
                        og_path := CreateOriginalName(path)
                        err := CopyFileContents(path, og_path)
                        if err != nil {
                                return err
                        } 
                        log.Printf("%q -> %q\n", path, og_path)
                } else if !recursive && path != dir {
                        log.Printf("skipping a directory without errors: %q\n", info.Name())
                        return filepath.SkipDir
                }
                return err
        })
        return err
}


func main() {
        log.Println("Backing up directories")
        BackupDir("./assets", "bak", true, true)
        log.Println("Restoring directories")
        RestoreDir("./assets", "bak", true, true)
}
