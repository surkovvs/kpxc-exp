package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/surkovvs/kpxc-exp/internal/entity"
	"github.com/surkovvs/kpxc-exp/internal/process/args"
	"github.com/surkovvs/kpxc-exp/internal/process/layout"
	"github.com/surkovvs/kpxc-exp/internal/process/parsing"
	"github.com/surkovvs/kpxc-exp/internal/process/tools"
)

func main() {
	arguments, err := args.ScanArgs()
	if err != nil {
		log.Fatalf("import running: %s", err.Error())
	}

	imported, err := tools.RunImport(arguments.PathKDBX, arguments.PasswordKDBX)
	if err != nil {
		log.Fatalf("import running: %s", err.Error())
	}

	kpf := entity.KeePassFile{}
	if err := xml.Unmarshal(imported, &kpf); err != nil {
		log.Fatalf("unmarshaling import: %s", err.Error())
	}

	envGroup := parsing.ParseGroups(kpf.Root.Group, 0)

	entrys, toPrint := layout.GetLayout(envGroup)

	fmt.Printf("%s\nEnter entry tag to export: ", toPrint)

	if err := tools.RunExport(tools.EntryChoose(entrys)); err != nil {
		log.Fatalf("exporting: %s", err.Error())
	}
}
