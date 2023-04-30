package protogen

import (
	"fmt"

	"github.com/nokamoto/2pf23/internal/protogen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
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

func (p *Plugin) codeGeneratorRequest(req *pluginpb.CodeGeneratorRequest) (*v1.Service, error) {
	var resp *v1.Service
	return resp, nil
}

func (p *Plugin) codeGeneratorResponse(svc *v1.Service) (*pluginpb.CodeGeneratorResponse, error) {
	var resp *pluginpb.CodeGeneratorResponse
	return resp, nil
}
