package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	ext := flag.String("e", "bak", "The name of the backup extension")
	restore := flag.Bool("r", false, "Restores the file from the backup if true, otherwise it backs the file up")
	verbose := flag.Bool("v", false, "Verbosely prints output")

	flag.Parse()

	if flag.NFlag() == 0 && flag.NArg() == 0 {
		fmt.Println("Usage of fbkp:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *restore {
		for _, file := range flag.Args() {
			bak_filename := CreateBackupName(file, *ext)
			err := RestoreFile(file, *ext)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if *verbose {
				fmt.Printf("%q -> %q\n", bak_filename, file)
			}
		}
	} else {
		for _, file := range flag.Args() {
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
	}
}
