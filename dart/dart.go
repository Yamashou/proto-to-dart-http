package dart

import (
	"fmt"
	"os"
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

	files := FileNames(apiParams)
	for i := range files {
		file := files[i]
		if _, err := fmt.Fprintf(g.File, "import 'package:%s%s%s.dart';\n", project, path, file); err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}

	return nil
}

func WriteClass(g *GenerateDart, apiParams []*APIParam, project string) error {
	_, err := fmt.Fprintf(g.File, "class %sClient {\n", project)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	_, err = fmt.Fprint(g.File, "\tString baseUrl;\n")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	_, err = fmt.Fprintf(g.File, "\t%sClient(String baseUrl) { this.baseUrl = baseUrl}\n", project)
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
