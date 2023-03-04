package fbkp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateBackupPath takes a filepath and an extension,
// such as bak, and returns a  filepath with that
// extension appended.
// For example: CreateBackupPath("test.txt", "bak")
// will return "/path/to/test.txt.bak"
func CreateBackupPath(filename, ext string) (string, error) {
	backup_name := filename + "." + ext
	backup_path, err := filepath.Abs(backup_name)
	if err != nil {
		return "", err
	} else {
		return backup_path, nil
	}
}

// CreateOriginalPath takes a filename and returns a
// filepath with the last extension removed.
// For example: CreateOriginalPath("test.txt.bak")
// will return "/path/to/test.txt"
func CreateOriginalPath(filename string) (string, error) {
	original_name := strings.TrimSuffix(filename, filepath.Ext(filename))
	original_path, err := filepath.Abs(original_name)
	if err != nil {
		return "", err
	} else {
		return original_path, nil
	}
}

// CopyFileContents takes a source filepath and a
// destination filepath, to copy the contents
// and file permissions of the source filepath
// to the destination filepath, creating the destination
// filepath if it does not exist.
// For example, BackupFile("test.txt", "test.txt.bak") copies "test.txt"
// to "test.txt.bak"
func CopyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	permissions := info.Mode().Perm()
	err = os.WriteFile(dst, buf, permissions)
	if err != nil {
		return err
	}

	return err
}

// BackupFile takes a source filepath and an extension
// and copies the contents and permissions of the source filepath
// to a destination filepath which is the source filepath appended
// with the extension.  BackupFile returns any error that occurred.
// For example, BackupFile("test.txt", "bak") copies "test.txt"
// to "test.txt.bak"
func BackupFile(path, ext string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		bkp_path, err := CreateBackupPath(path, ext)
		if err != nil {
			return err
		}

		err = CopyFileContents(path, bkp_path)
		if err != nil {
			return err
		}
	}
	return err
}

// RestoreFile takes a source filepath assuming it has a backup extension
// (i.e. .bak), and copies the source filepath to a destination filepath
// with that backup extension removed.  RestoreFile returns any error that occurred.
// For example, RestoreFile("test.txt.bak") copies "test.txt.bak"
// to "test.txt".
func RestoreFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		og_path, err := CreateOriginalPath(path)
		if err != nil {
			return err
		}

		err = CopyFileContents(path, og_path)
		if err != nil {
			return err
		}
	}
	return err
}

// BackupDir takes a directory, an extension, a verbose flag, and a recursive
// flag.  BackupDir backs up the dir, ignoring files ending with the
// ext string, verbosely showing the copying if verbose is set to True, and
// backing up directories in that dir if recursive is set to True.  BackupDir
// returns any error that occurred.
func BackupDir(dir string, ext string, verbose bool, recursive bool) error {
	real_ext := "." + ext
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("skipping a directory without errors: %q\n", info.Name())
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
			return err
		}

		fname := info.Name()
		if !info.IsDir() && filepath.Ext(fname) != real_ext {
			bkp_path, err := CreateBackupPath(path, ext)
			if err != nil {
				return err
			}

			err = CopyFileContents(path, bkp_path)
			if err != nil {
				return err
			}

			if verbose {
				fmt.Printf("%q -> %q\n", path, bkp_path)
			}
		} else if info.IsDir() && !recursive {
			fmt.Printf("skipping a directory without errors: %q\n", info.Name())
			return filepath.SkipDir
		}
		return err
	})
	return err
}

// RestoreDir takes a directory, an extension, a verbose flag, and a recursive
// flag.  RestoreDir restores the dir from files ending with the ext string,
// verbosely showing the copying if verbose is set to True, and restoring
// directories in that dir if recursive is set to True.  RestoreDir returns any
// error that occurred.
func RestoreDir(dir string, ext string, verbose bool, recursive bool) error {
	real_ext := "." + ext
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
			return err
		}

		fname := info.Name()
		if !info.IsDir() && filepath.Ext(fname) == real_ext {
			og_path, err := CreateOriginalPath(path)
			if err != nil {
				return err
			}

			err = CopyFileContents(path, og_path)
			if err != nil {
				return err
			}

			if verbose {
				fmt.Printf("%q -> %q\n", path, og_path)
			}
		} else if info.IsDir() && !recursive {
			fmt.Printf("skipping a directory without errors: %q\n", info.Name())
			return filepath.SkipDir
		}
		return nil
	})
	return err
}
