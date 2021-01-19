package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const PROMPT = ">>"

func Start(cmd *cobra.Command, args []string) error {
	r := &REPL{
		DebugMode: false,
	}
	return r.Start(os.Stdin, os.Stdout)
}

type REPL struct {
	DebugMode bool
}

func (repl *REPL) Start(in io.ReadCloser, out io.Writer) error {
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprintf(out, PROMPT)

		bb := bytes.NewBuffer([]byte(""))
		for i := 0; ; i++ {
			buf, cont, err := r.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			if i == 0 {
				bb = bytes.NewBuffer([]byte(""))
			}
			_, err = bb.Write(buf)
			if err != nil {
				return err
			}
			if cont {
				continue
			}
			if len(buf) == 0 {
				break
			} else {
				bb.Write([]byte("\n"))
			}
		}

		line := bb.String()

		if strings.HasPrefix(line, ".") {
			isExit := repl.parseCommand(line, out)
			if isExit {
				return nil
			}
		} else {
		}
		fmt.Fprintf(out, "\n")
	}
}

func (repl *REPL) parseCommand(line string, out io.Writer) bool {
	parsedLine := strings.SplitN(line, " ", 2)
	cmd := parsedLine[0]
	switch strings.ToLower(cmd) {
	case ".exit":
		return true
	case ".debug":
		repl.DebugMode = !repl.DebugMode
		return false
	default:
		fmt.Fprintf(out, "Unknown Command\n")
	}
	return false
}
