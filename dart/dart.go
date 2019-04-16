package dart

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/xerrors"
)

func DartFileName(name string) string {
	return fmt.Sprintf("%s.dart", name)
}

type APIParam struct {
	HTTPMethod string
	APIName    string
	Path       string
	Body       string
	FileName   string
	Request    Request
	Response   Response
}

func FileNames(as []*APIParam) []string {
	var names []string
	m := make(map[string]struct{})

	for i := range as {
		a := as[i]
		_, ok := m[a.FileName]
		if !ok {
			m[a.FileName] = struct{}{}
			names = append(names, a.FileName)
		}

		_, ok = m[a.Response.FileName]
		if !ok {
			m[a.Response.FileName] = struct{}{}
			names = append(names, a.Response.FileName)
		}

		_, ok = m[a.Request.FileName]
		if !ok {
			m[a.Request.FileName] = struct{}{}
			names = append(names, a.Request.FileName)
		}
	}
	return names
}

type Request struct {
	Name     string
	FileName string
}

type Response struct {
	Name     string
	FileName string
}

type GenerateDart struct {
	File *os.File
}

func NewGenerateDart() (*GenerateDart, error) {
	file, err := os.OpenFile("test.pb.dart", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, xerrors.Errorf("failed to get absolute path: %w", err)
	}

	return &GenerateDart{File: file}, nil
}

func WriteImports(g *GenerateDart, apiParams []*APIParam, project, path string) error {
	_, err := fmt.Fprint(g.File, "import 'package:http/http.dart' as http;\n")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	//TODO ここはちゃんと関数に分けて切り出す
	dartProject := strings.ReplaceAll(project, "-", "_")
	files := FileNames(apiParams)
	for i := range files {
		file := files[i]
		sliceFile := strings.Split(file, "/")
		dartFile := strings.ReplaceAll(sliceFile[len(sliceFile)-1], "proto", "pb")
		if _, err := fmt.Fprintf(g.File, "import 'package:%s%s%s.dart';\n", dartProject, path, dartFile); err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	return nil
}

//https://github.com/iancoleman/strcase

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

func toCamel(s string) string {
	return toCamelInitCase(s, true)
}

func WriteClass(g *GenerateDart, apiParams []*APIParam, project string) error {
	camelProject := toCamel(project)
	_, err := fmt.Fprintf(g.File, "class %sClient {\n", camelProject)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	_, err = fmt.Fprint(g.File, "\tString baseUrl;\n")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	_, err = fmt.Fprintf(g.File, "\t%sClient(String baseUrl) {this.baseUrl = baseUrl;}\n", camelProject)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	for i := range apiParams {
		apiParam := apiParams[i]
		_, err = fmt.Fprintf(g.File,
			"\tFuture<%s> %c%s(%s body, Map<String, String> headers) async {\n"+
				"\t\tfinal response = await http.%s(\n"+
				"\t\t\tthis.baseUrl + \"%s\",\n"+
				"\t\t\tbody: body,\n"+
				"\t\t\theaders: headers);\n\n"+
				"\t\tfinal %s res = %s.fromBuffer(response.bodyBytes);\n"+
				"\t\treturn res;\n\t}\n\n",
			apiParam.Response.Name,
			strings.ToLower(apiParam.APIName)[0],
			apiParam.APIName[1:],
			apiParam.Request.Name,
			strings.ToLower(apiParam.HTTPMethod),
			apiParam.Path,
			apiParam.Response.Name,
			apiParam.Response.Name,
		)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	_, err = fmt.Fprint(g.File, "}\n")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func Build(apiParams []*APIParam, project, path string) (*GenerateDart, error) {
	g, err := NewGenerateDart()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if err := WriteImports(g, apiParams, project, path); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	if err := WriteClass(g, apiParams, project); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return g, nil
}
