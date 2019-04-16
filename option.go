package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

type Option struct {
	importPaths []string
	projectName string
	outPutPass  string
}

func parseOption() (*Option, []string, error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
   %s [OPTIONS] [pb files...]
Options\n`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	dir, err := os.Getwd()
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to get current directory: %v", err)
	}

	projectName := flag.String("p", "", `project name`)
	outPutPass := flag.String("o", "", `project name`)
	protoImportOpt := flag.String("i", dir, `pb files import directory`)
	flag.Parse()

	protoImportPaths := strings.Split(*protoImportOpt, ",")
	for i := range protoImportPaths {
		protoImportPath, err := filepath.Abs(protoImportPaths[i])
		if err != nil {
			return nil, nil, xerrors.Errorf("failed to get absolute path: %v", err)
		}
		protoImportPaths[i] = protoImportPath
	}

	return &Option{
		importPaths: protoImportPaths,
		projectName: *projectName,
		outPutPass:  *outPutPass,
	}, flag.Args(), nil
}
func (o *Option) ImportRootPath() string {
	ss := strings.Split(o.outPutPass, "/lib")
	return fmt.Sprintf("package:%s.dart", o.projectName) + ss[1]
}
