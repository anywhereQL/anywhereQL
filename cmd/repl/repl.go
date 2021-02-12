package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anywhereQL/anywhereQL/common/logger"
	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/anywhereQL/anywhereQL/runtime"
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
			rt := runtime.New()
			rs, err := rt.Start(line)
			if err != nil {
				logger.Errorf("%#+v", err)
				continue
			}
			for _, line := range rs {
				for i, col := range line {
					switch col.Type {
					case value.INTEGER:
						fmt.Fprintf(out, "%d", col.Int)
					case value.FLOAT:
						fmt.Fprintf(out, "%f", col.Float)
					case value.STRING:
						fmt.Fprintf(out, "%s", col.String)
					case value.NULL:
						fmt.Fprintf(out, "NULL")
					case value.BOOL:
						if col.Bool.True {
							fmt.Fprintf(out, "TRUE")
						} else if col.Bool.False {
							fmt.Fprintf(out, "FALSE")
						}
					}
					if i != (len(line) - 1) {
						fmt.Fprintf(out, ",")
					}
				}
				fmt.Fprintf(out, "\n")
			}
		}
		fmt.Fprintf(out, "\n")
	}
}

func (repl *REPL) parseCommand(line string, out io.Writer) bool {
	parsedLine := strings.SplitN(strings.Trim(line, "\n"), " ", 2)
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
