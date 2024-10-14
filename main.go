package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/sylabs/singularity/v4/pkg/build/types/parser"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run() error {
	var (
		in  io.Reader = os.Stdin
		out io.Writer = os.Stdout
	)

	if len(os.Args) > 1 && os.Args[1] != "-" {
		f, err := os.Open(os.Args[1])
		if err != nil {
			return err
		}

		defer f.Close()

		in = f
	}

	if len(os.Args) == 3 && os.Args[2] != "-" {
		f, err := os.Create(os.Args[1])
		if err != nil {
			return err
		}

		defer f.Close()

		out = f
	}

	return processFile(in, out)
}

func processFile(in io.Reader, out io.Writer) error {
	defsPreBuildArgs, err := parser.All(in)
	if err != nil {
		return err
	}

	for n, buildStage := range defsPreBuildArgs {
		if n > 0 {
			fmt.Fprintln(out)
		}

		bootstrap := buildStage.Header["bootstrap"]
		from := buildStage.Header["from"]
		stage := buildStage.Header["stage"]

		if bootstrap != "docker" && bootstrap != "" {
			return errors.New("unsupported bootstrap")
		}

		if stage != "" {
			fmt.Fprintf(out, "FROM %s AS %s\n\n", from, stage)
		} else {
			fmt.Fprintf(out, "FROM %s\n\n", from)
		}

		for _, files := range buildStage.BuildData.Files {
			stage := files.Stage()

			for _, file := range files.Files {
				if stage != "" {
					fmt.Fprintf(out, "COPY --from=%s %s %s\n", stage, file.Src, file.Dst)
				} else {
					fmt.Fprintf(out, "COPY %s %s\n", file.Src, file.Dst)
				}
			}
		}

		EOF := "EOF"

		for n := 0; strings.Contains(buildStage.BuildData.Post.Script, EOF); n++ {
			EOF = "EOF" + strconv.Itoa(n)
		}

		fmt.Fprintf(out, "\nRUN <<%s\n%s\n%s\n", EOF, "#!/bin/bash\nset -eo pipefail\ndeclare SINGULARITY_ENVIRONMENT=/etc/profile.d/99-docker-env.sh\n"+buildStage.BuildData.Post.Script, EOF)
	}

	fmt.Fprintln(out, "\nENTRYPOINT [\"/bin/bash\", \"-l\"]")

	return nil
}
