package main 

import (
	"fmt"
	"os"

        "github.com/noaoh/fbkp"
	"github.com/spf13/pflag"
)

type fbkpError struct {
	path, message string
	err           error
}

func (e *fbkpError) Error() string {
	switch e.err {
	case nil:
		return fmt.Sprintf("%s: %q.\n", e.message, e.path)
	default:
		return fmt.Sprintf("%s %q: %s.\n", e.message, e.path, e.err)
	}
}

func main() {
	ext := pflag.StringP("ext", "e", "bak", "The name of the backup extension.")
	recur := pflag.BoolP("recursive", "r", false, "Recursively backs up files in directories.")
	backup := pflag.BoolP("backup", "b", false, "Backs up the file(s) if this flag is passed, otherwise it restores the file(s).")
	verbose := pflag.BoolP("verbose", "v", false, "Verbosely prints output.")

	pflag.Parse()

	if pflag.NFlag() == 0 && pflag.NArg() == 0 {
		fmt.Println("Usage of fbkp:")
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if *backup {
		for _, path := range pflag.Args() {
			info, err := os.Stat(path)

			if os.IsNotExist(err) {
				fmt.Print((&fbkpError{message: "Path does not exist", path: path, err: nil}).Error())
                                continue
			}

			if info.IsDir() {
				err := fbkp.BackupDir(path, *ext, *verbose, *recur)
				if err != nil {
					fmt.Print((&fbkpError{message: "Backup failed for directory", path: path, err: err}).Error())
				}
			} else {
				bkp_path := fbkp.CreateBackupName(path, *ext)

				err := fbkp.BackupFile(path, *ext)
				if err != nil {
					fmt.Print((&fbkpError{message: "Backup failed for file", path: path, err: err}).Error())
				}

				if *verbose && err == nil {
					fmt.Printf("%q -> %q\n", path, bkp_path)
				}
			}
		}
	} else {
		for _, path := range pflag.Args() {
			info, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Print((&fbkpError{message: "Path does not exist", path: path, err: nil}).Error())
                                continue
			}

			if info.IsDir() {
				err := fbkp.RestoreDir(path, *ext, *verbose, *recur)
				if err != nil {
					fmt.Print((&fbkpError{message: "Restore failed for directory", path: path, err: err}).Error())
				}
			} else {
				orig_path := fbkp.CreateOriginalName(path)
				err := fbkp.RestoreFile(path)
				if err != nil {
					fmt.Print((&fbkpError{message: "Restore failed for file", path: path, err: err}).Error())
				}

				if *verbose && err == nil {
					fmt.Printf("%q -> %q\n", path, orig_path)
				}
			}
		}
	}
}
