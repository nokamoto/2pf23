package protogen

import (
	"fmt"

	"github.com/nokamoto/2pf23/internal/protogen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin is a protoc plugin to generate server code.
type Plugin struct {
	*protogen.Plugin
}

// NewPlugin returns a new Plugin with stdin and stdout.
func NewPlugin() *Plugin {
	var p *Plugin
	p = &Plugin{
		protogen.NewPlugin(func(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
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
			resp = append(resp, p.service(svc, file))
		}
	}
	return resp, nil
}

func (p *Plugin) service(svc *descriptorpb.ServiceDescriptorProto, file *descriptorpb.FileDescriptorProto) *v1.Service {
	api := protogen.NewAPIDescriptor(file)
	resp := &v1.Service{
		Name:                svc.GetName(),
		ApiVersion:          api.APIVersion(),
		ApiImportPath:       api.ImportPath(),
		UnimplementedServer: fmt.Sprintf("Unimplemented%sServer", svc.GetName()),
	}
	return resp
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
