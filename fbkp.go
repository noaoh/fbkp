package main 

import (
        "log"
        "os"
        "io/ioutil"
)

func BackupFile(name string, suffix string) error {
        og_file, err := os.Open(name)
        if err != nil {
                log.Fatal(err)
        }
        defer og_file.Close()

        buf, err := ioutil.ReadAll(og_file)
        if err != nil {
                log.Fatal(err)
        }

        bkp_name := name + "." + suffix
        _, err = os.Create(bkp_name)
        if err != nil {
                log.Fatal(err)
        }

        err = os.Chmod(bkp_name, 0644)
        if err != nil {
                log.Fatal(err)
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
                log.Fatal(err)
        }
        defer bkp_file.Close()

        buf, err := ioutil.ReadAll(bkp_file)
        if err != nil {
                log.Fatal(err)
        }

        err = ioutil.WriteFile(name, buf, 0755)
        if err != nil {
                log.Fatal(err)
        }
        return err
}
