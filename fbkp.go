package main

import (
	"io"
        "io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CreateBackupName(filename, ext string) string {
	return filename + "." + ext
}

func CopyFileContents(src, dst string) error {
        in, err := os.Open(src)
        if err != nil {
                return err
        }
        defer in.Close()

        out, err := os.Create(dst)
        if err != nil {
                return err
        }
        defer out.Close()

        if _, err = io.Copy(out, in); err != nil {
                return err
        }

        err = out.Sync()
        return err
}

func BackupFile(name, ext string) error {
        bkp_name := CreateBackupName(name, ext)
        err := CopyFileContents(name, bkp_name)
	if err != nil {
		return err
	}
        return err
}

func RestoreFile(name, ext string) error {
	bkp_name := CreateBackupName(name, ext)
        err := CopyFileContents(bkp_name, name)
	if err != nil {
		return err
	}
        return err
}

func RestoreDir(dir, ext string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	real_ext := "." + ext
	for _, file := range files {
		fname := file.Name()

		if filepath.Ext(fname) == real_ext {
			restore_file := strings.TrimSuffix(fname, real_ext)

			err := RestoreFile(restore_file, ext)
			if err != nil {
				return err
			}
		}
	}
	return err
}
