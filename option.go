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
	packagePath string
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
	outPutPass := flag.String("o", "", `out put pass`)
	packagePath := flag.String("pp", "/", `project name`)
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

	path := *packagePath
	if *packagePath != "/" {
		if path[0] != '/' {
			path = "/" + path
		}
		last := path[len(path)-1]
		if last != '/' {
			path += "/"
		}
	}

	return &Option{
		importPaths: protoImportPaths,
		projectName: *projectName,
		outPutPass:  *outPutPass,
		packagePath: path,
	}, flag.Args(), nil
}
func (o *Option) ImportRootPath() string {
	ss := strings.Split(o.outPutPass, "/lib")
	return fmt.Sprintf("package:%s.dart", o.projectName) + ss[1]
}
