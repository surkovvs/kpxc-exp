package tools

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/surkovvs/kpxc-exp/internal/entity"
	"golang.org/x/term"
)

const separator string = `|&;`

func RunImport(path, password string) ([]byte, error) {
	if password == "" {
		dbRawPass, err := term.ReadPassword(0)
		if err != nil {
			return nil, fmt.Errorf("reading pass error: %w", err)
		}
		password = string(dbRawPass)
	}

	kpxc := exec.Command(
		"keepassxc-cli",
		"export",
		"-q",
		path,
	)

	in, err := kpxc.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("piping stdin keepassxc-cli error: %w", err)
	}
	go func() {
		defer in.Close()
		if _, err := io.WriteString(in, password); err != nil {
			log.Fatalf("write to executed process: %s", err.Error())
		}
	}()

	out, err := kpxc.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("running keepassxc-cli error: %w; message:%s", err, string(out))
	}
	return out, nil
}

func EntryChoose(entrys map[string]entity.EnvEntry) string {
	var (
		ok      bool
		toSetup string
		entry   entity.EnvEntry
	)
	for !ok {
		fmt.Fscan(os.Stdin, &toSetup)
		fmt.Printf("%s", toSetup)
		entry, ok = entrys[toSetup]
		if !ok {
			fmt.Printf("Incorrect input, retry pls.\n")
		}
	}

	envs := make([]string, 0, len(entry.Envs))
	for _, env := range entry.Envs {
		envs = append(envs, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}

	return strings.Join(envs, separator)
}

func RunExport(output string) error {
	var err error
	fd3 := os.NewFile(3, "fd3")
	defer func() {
		if err = fd3.Close(); err != nil {
			err = fmt.Errorf("closing fd3: %w", err)
		}
	}()

	if _, err = fd3.WriteString(output); err != nil {
		return fmt.Errorf("write to fd3: %w", err)
	}

	return err
}
