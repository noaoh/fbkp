package main 
import (
        "errors"
        "fmt"
        "flag"
        "os"
)

func main() {
        file := flag.String("f", "", "The file(s) to backup")
        bak := flag.String("b", "bak", "The name of the backup extension")
        restore := flag.Bool("r", false, "Restores the file from the backup if true, otherwise it backs the file up")
        verbose := flag.Bool("v", false, "Verbosely prints output")

        flag.Parse()

        if flag.NFlag() == 0 {
                fmt.Println("Usage of fbkp:")
                flag.PrintDefaults()
                os.Exit(0)
        }

        if *file == "" {
                fmt.Println(errors.New("No file passed as an argument"))
                os.Exit(1)
        }

        bak_filename := *file + "." + *bak
        if *restore {
                err := RestoreFile(*file, *bak)
                if err != nil {
                        fmt.Println(err)
                        os.Exit(1)
                }
                
                if *verbose {
                        fmt.Printf("%q -> %q\n", bak_filename, *file)
                }
        } else {
                err := BackupFile(*file, *bak)
                if err != nil {
                        fmt.Println(err)
                        os.Exit(1)
                }

                if *verbose {
                        fmt.Printf("%q -> %q\n", *file, bak_filename)
                }
        }
}
