package main

import (
        "github.com/spf13/pflag"
	"fmt"
        "path/filepath"
	"os"
)

func main() {
	ext := pflag.StringP("ext", "e", "bak", "The name of the backup extension")
        // recur := pflag.Bool("r", false, "Recursively backs up files in directories.")
	backup := pflag.BoolP("backup", "b", false, "Backs up the file(s) if this flag is passed, otherwise it restores the file(s)")
	verbose := pflag.BoolP("verbose", "v", false, "Verbosely prints output")

	pflag.Parse()

	if pflag.NFlag() == 0 && pflag.NArg() == 0 {
		fmt.Println("Usage of fbkp:")
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if *backup {
		for _, file := range pflag.Args() {
			bak_filename := CreateBackupName(file, *ext)
			err := BackupFile(file, *ext)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if *verbose {
				fmt.Printf("%q -> %q\n", file, bak_filename)
			}
		}
	} else {
                real_ext := "." + *ext
		for _, file := range pflag.Args() {
                        if filepath.Ext(file) == real_ext {
                                og_filename := CreateOriginalName(file)
                                err := RestoreFile(file)

                                if err != nil {
                                        fmt.Println(err)
                                        os.Exit(1)
                                }

                                if *verbose {
                                        fmt.Printf("%q -> %q\n", file, og_filename)
                                }
                        }

		}
        }
}
