package protogen

import (
	"fmt"

	"github.com/nokamoto/2pf23/internal/protogen/core"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin is a protoc plugin to generate server code.
type Plugin struct {
	*core.Plugin
}

// NewPlugin returns a new Plugin with stdin and stdout.
func NewPlugin() *Plugin {
	var p *Plugin
	p = &Plugin{
		core.NewPlugin(func(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
			svc, err := p.codeGeneratorRequest(req)
			if err != nil {
				return nil, fmt.Errorf("failed to generate code: %w", err)
			}
			resp, err := p.codeGeneratorResponse(svc)
			if err != nil {
				return nil, fmt.Errorf("failed to generate code: %w", err)
			}
			return resp, nil
		}),
	}
	return p
}

func (p *Plugin) codeGeneratorRequest(req *pluginpb.CodeGeneratorRequest) ([]*v1.Service, error) {
	var resp []*v1.Service
	for _, file := range req.GetProtoFile() {
		for _, svc := range file.GetService() {
			s, err := p.service(svc, file)
			if err != nil {
				return nil, fmt.Errorf("failed to generate service: %w", err)
			}
			resp = append(resp, s)
		}
	}
	return resp, nil
}

func (p *Plugin) setCall(call *v1.Call, m *core.MethodDescriptor) {
	call.Name = m.GetName()
	call.RequestType = core.GoTypeNameFromFullyQualified(m.GetInputType())
	call.ResponseType = core.GoTypeNameFromFullyQualified(m.GetOutputType())
	call.ResourceType = core.GoTypeNameFromFullyQualified(m.GetOutputType())
}

func (p *Plugin) createCall(m *core.MethodDescriptor) *v1.Call {
	accessor := fmt.Sprintf("Get%s", m.ResourceName())
	resp := &v1.Call{
		MethodType:        v1.MethodType_METHOD_TYPE_CREATE,
		GetResourceMethod: accessor,
	}
	p.setCall(resp, m)
	return resp
}

func (p *Plugin) getCall(m *core.MethodDescriptor) *v1.Call {
	resp := &v1.Call{
		MethodType: v1.MethodType_METHOD_TYPE_GET,
	}
	p.setCall(resp, m)
	return resp
}

func (p *Plugin) deleteCall(m *core.MethodDescriptor) *v1.Call {
	resp := &v1.Call{
		Name:         m.GetName(),
		MethodType:   v1.MethodType_METHOD_TYPE_DELETE,
		RequestType:  core.GoTypeNameFromFullyQualified(m.GetInputType()),
		ResponseType: "empty.Empty",
	}
	return resp
}

func (p *Plugin) listCall(file *descriptorpb.FileDescriptorProto, m *core.MethodDescriptor, api *core.APIDescriptor) (*v1.Call, error) {
	resp := &v1.Call{
		Name:         m.GetName(),
		MethodType:   v1.MethodType_METHOD_TYPE_LIST,
		RequestType:  core.GoTypeNameFromFullyQualified(m.GetInputType()),
		ResponseType: core.GoTypeNameFromFullyQualified(m.GetOutputType()),
		ResourceType: fmt.Sprintf("%s.%s", api.APIVersion(), m.ResourceName()),
	}
	listResponse := core.NewListResponseDescriptor(file, m.MethodDescriptorProto)
	listField, err := listResponse.ListField()
	if err != nil {
		return nil, fmt.Errorf("failed to get list field: %w", err)
	}
	resp.ListField = listField
	return resp, nil
}

func (p *Plugin) updateCall(m *core.MethodDescriptor) *v1.Call {
	accessor := fmt.Sprintf("Get%s", m.ResourceName())
	resp := &v1.Call{
		MethodType:        v1.MethodType_METHOD_TYPE_UPDATE,
		GetResourceMethod: accessor,
	}
	p.setCall(resp, m)
	return resp
}

func (p *Plugin) service(svc *descriptorpb.ServiceDescriptorProto, file *descriptorpb.FileDescriptorProto) (*v1.Service, error) {
	api := core.NewAPIDescriptor(file)
	resp := &v1.Service{
		Name:          api.ServiceName(),
		ApiVersion:    api.APIVersion(),
		ApiImportPath: api.ImportPath(),
	}
	for _, method := range svc.GetMethod() {
		m := core.NewMethodDescriptor(method)
		switch m.Type() {
		case v1.MethodType_METHOD_TYPE_CREATE:
			resp.Calls = append(resp.Calls, p.createCall(m))

		case v1.MethodType_METHOD_TYPE_GET:
			resp.Calls = append(resp.Calls, p.getCall(m))

		case v1.MethodType_METHOD_TYPE_DELETE:
			resp.Calls = append(resp.Calls, p.deleteCall(m))

		case v1.MethodType_METHOD_TYPE_LIST:
			call, err := p.listCall(file, m, api)
			if err != nil {
				return nil, fmt.Errorf("failed to generate list call: %w", err)
			}
			resp.Calls = append(resp.Calls, call)

		case v1.MethodType_METHOD_TYPE_UPDATE:
			resp.Calls = append(resp.Calls, p.updateCall(m))
		}
	}
	return resp, nil
}

func (p *Plugin) codeGeneratorResponse(services []*v1.Service) (*pluginpb.CodeGeneratorResponse, error) {
	var resp pluginpb.CodeGeneratorResponse
	for _, svc := range services {
		content, err := p.MarshalJsonProto(svc)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}

		filename := fmt.Sprintf("%s.%s.json", svc.GetName(), svc.GetApiVersion())
		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(filename),
			Content: proto.String(string(content)),
		})
	}
	return &resp, nil
}
