package args

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/term"
)

type Args struct {
	PasswordKDBX string
	PathKDBX     string
}

var pathDB, password *string

func init() {
	password = flag.String("p", "", "password to kdbx")
	pathDB = flag.String("d", "", "path to kdbx")
	flag.Parse()
}

func ScanArgs() (Args, error) {
	var args Args

	if *pathDB != "" {
		args.PathKDBX = *pathDB
	} else {
		fmt.Println("enter path to *.kdbx file")
		if _, err := fmt.Scanln(&args.PathKDBX); err != nil {
			log.Fatal("lable 1 ", err)
		}
	}

	if *password != "" {
		args.PasswordKDBX = *password
	} else {
		fmt.Println("enter password")
		rawPass, err := term.ReadPassword(0)
		if err != nil {
			log.Fatal("lable 1 ", err)
		}
		args.PasswordKDBX = string(rawPass)
	}

	return args, nil
}
