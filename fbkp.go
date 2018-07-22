package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func BackupFile(name string, suffix string) error {
	og_file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer og_file.Close()

	buf, err := ioutil.ReadAll(og_file)
	if err != nil {
		return err
	}

	bkp_name := name + "." + suffix
	_, err = os.Create(bkp_name)
	if err != nil {
		return err
	}

	err = os.Chmod(bkp_name, 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(bkp_name, buf, 0755)
	if err != nil {
		return err
	}

	return err
}

func RestoreFile(name string, suffix string) error {
	bkp_name := name + "." + suffix
	bkp_file, err := os.Open(bkp_name)
	if err != nil {
		return err
	}
	defer bkp_file.Close()

	buf, err := ioutil.ReadAll(bkp_file)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(name, buf, 0755)
	if err != nil {
		return err
	}
	return err
}

func RestoreDir(dir string, suffix string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	real_suffix := "." + suffix
	for _, file := range files {
		fname := file.Name()

		if filepath.Ext(fname) == real_suffix {
			restore_file := strings.TrimSuffix(fname, real_suffix)

			err := RestoreFile(restore_file, suffix)
			if err != nil {
				return err
			}
		}
	}
	return err
}

