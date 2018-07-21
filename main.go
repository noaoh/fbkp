package main 
import (
        "fmt"
        "flag"
        "os"
        "path/filepath"
)

func main() {
        file := flag.String("f", "", "The file(s) to backup")
        ext := flag.String("e", "bak", "The name of the backup extension")
        restore := flag.Bool("r", false, "Restores the file from the backup if true, otherwise it backs the file up")
        verbose := flag.Bool("v", false, "Verbosely prints output")

        flag.Parse()

        if flag.NFlag() == 0 {
                fmt.Println("Usage of fbkp:")
                flag.PrintDefaults()
                os.Exit(0)
        }

        if *file == "" {
                fmt.Println("No file(s) passed as an argument")
                os.Exit(1)
        }

        files, err := filepath.Glob(*file)
        if err != nil {
                fmt.Println(err)
        }


        for _, file := range files {
                bak_filename := file + "." + *ext
                if *restore {
                        err := RestoreFile(file, *ext)
                        if err != nil {
                                fmt.Println(err)
                                os.Exit(1)
                        }
                        
                        if *verbose {
                                fmt.Printf("%q -> %q\n", bak_filename, file)
                        }
                } else {
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
