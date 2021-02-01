package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anywhereQL/anywhereQL/common/debug"
	"github.com/anywhereQL/anywhereQL/common/result"
	"github.com/anywhereQL/anywhereQL/compiler/lexer"
	"github.com/anywhereQL/anywhereQL/compiler/parser"
	"github.com/anywhereQL/anywhereQL/compiler/planner"
	"github.com/anywhereQL/anywhereQL/runtime/vm"
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
			tokens := lexer.Lex(line)
			debug.PrintToken(out, tokens)

			ast, err := parser.Parse(tokens)
			if err != nil {
				fmt.Fprintf(out, "%v", err)
				continue
			}
			debug.PrintAST(out, ast)

			vc := planner.Translate(ast)
			debug.PrintVC(out, vc)

			rs, err := vm.Run(vc)
			if err != nil {
				fmt.Fprintf(out, "%v", err)
				continue
			}
			for i, col := range rs {
				switch col.Type {
				case result.Integral:
					fmt.Fprintf(out, "%d", col.Integral)
				case result.Float:
					fmt.Fprintf(out, "%f", col.Float)
				}
				if i != (len(rs) - 1) {
					fmt.Fprintf(out, ",")
				}
			}
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
