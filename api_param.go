package main

import (
	"github.com/Yamashou/proto-to-dart-http/dart"
	pbdesc "github.com/Yamashou/proto-to-dart-http/pbdsec"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"golang.org/x/xerrors"
	"google.golang.org/genproto/googleapis/api/annotations"
)

type apiParamsBuilder struct {
	pbdesc *pbdesc.ProtoDescriptor
}

func NewAPIParamsBuilder() *apiParamsBuilder {
	return &apiParamsBuilder{
		pbdesc: &pbdesc.ProtoDescriptor{},
	}
}

func (a *apiParamsBuilder) Build(fds []*desc.FileDescriptor) ([]*dart.APIParam, error) {
	var apiParams []*dart.APIParam
	for _, fd := range fds {
		for _, service := range fd.GetServices() {
			for _, method := range service.GetMethods() {
				params, err := a.build(method, service)
				if err != nil {
					return nil, xerrors.Errorf(": %w", err)
				}

				apiParams = append(apiParams, setFileName(params, fd.GetName())...)
			}
		}
	}

	return apiParams, nil
}

func setFileName(as []*dart.APIParam, name string) []*dart.APIParam {
	for i := range as {
		a := as[i]
		a.FileName = name
	}

	return as
}

func (a *apiParamsBuilder) build(method *desc.MethodDescriptor, service *desc.ServiceDescriptor) ([]*dart.APIParam, error) {
	opts := method.GetOptions()

	if !proto.HasExtension(opts, annotations.E_Http) {
		return []*dart.APIParam{}, nil
	}

	ext, err := proto.GetExtension(opts, annotations.E_Http)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	rule, ok := ext.(*annotations.HttpRule)
	if !ok {
		return nil, xerrors.New("annotation extension assertion error")
	}

	apiParams, err := a.apiParamsByHTTPRule(rule, method.GetInputType(), method.GetOutputType(), method.GetName())
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return apiParams, nil
}

func (a *apiParamsBuilder) apiParamsByHTTPRule(rule *annotations.HttpRule, inputType *desc.MessageDescriptor, outputType *desc.MessageDescriptor, name string) ([]*dart.APIParam, error) {
	var apiParams []*dart.APIParam

	apiParam, err := a.apiParamByHTTPRule(rule, inputType, outputType, name)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	apiParams = append(apiParams, apiParam)

	for _, r := range rule.GetAdditionalBindings() {
		apiParam, err := a.apiParamByHTTPRule(r, inputType, outputType, name)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		apiParams = append(apiParams, apiParam)
	}

	return apiParams, nil
}

func (a *apiParamsBuilder) apiParamByHTTPRule(rule *annotations.HttpRule, inputType *desc.MessageDescriptor, outputType *desc.MessageDescriptor, name string) (*dart.APIParam, error) {
	endpoint, err := newEndpoint(rule)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	bodyMsgType, err := a.pbdesc.BodyMsgTypeNameByHTTPRule(inputType, rule)
	bodyNotFound := xerrors.Is(err, pbdesc.ErrBodyNotFound)
	if err != nil && !bodyNotFound {
		return nil, xerrors.Errorf(": %w", err)
	}

	var jsonBody string
	if !bodyNotFound {
		jsonBody, err = a.pbdesc.JSONBody(bodyMsgType)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
	}

	request := dart.Request{
		Name:     inputType.GetName(),
		FileName: inputType.GetFile().GetName(),
	}

	response := dart.Response{
		Name:     outputType.GetName(),
		FileName: outputType.GetFile().GetName(),
	}

	return &dart.APIParam{
		HTTPMethod: endpoint.method,
		Path:       endpoint.path,
		Body:       jsonBody,
		APIName:    name,
		Request:    request,
		Response:   response,
	}, nil
}

type endpoint struct {
	method string
	path   string
}

func newEndpoint(rule *annotations.HttpRule) (*endpoint, error) {
	var e *endpoint
	switch opt := rule.GetPattern().(type) {
	case *annotations.HttpRule_Get:
		e = &endpoint{"GET", opt.Get}
	case *annotations.HttpRule_Put:
		e = &endpoint{"PUT", opt.Put}
	case *annotations.HttpRule_Post:
		e = &endpoint{"POST", opt.Post}
	case *annotations.HttpRule_Delete:
		e = &endpoint{"DELETE", opt.Delete}
	case *annotations.HttpRule_Patch:
		e = &endpoint{"PATCH", opt.Patch}
	default:
		return nil, xerrors.New("annotation http rule method dose not support type")
	}

	return e, nil
}
