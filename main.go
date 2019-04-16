package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Yamashou/proto-to-dart-http/dart"

	"golang.org/x/xerrors"

	"github.com/jhump/protoreflect/desc/protoparse"
)

func main() {
	opt, paths, err := parseOption()
	if err != nil || len(paths) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if err := run(paths, opt.importPaths, opt.packagePath, opt.projectName); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)

	}

}

func run(files []string, importPaths []string, packagePath, projectName string) error {
	p := protoparse.Parser{
		ImportPaths: importPaths,
	}

	fds, err := p.ParseFiles(files...)
	if err != nil {
		return xerrors.Errorf("Unable to parse pb file: %v \n", err)
	}

	apiParamBuilder := NewAPIParamsBuilder()
	apiParams, err := apiParamBuilder.Build(fds)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	g, err := dart.Build(apiParams, projectName, packagePath)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	if err := g.File.Close(); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}
