package protogen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nokamoto/2pf23/internal/protogen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin is a protoc plugin to generate cli code.
type Plugin struct {
	protogen.Plugin
}

// NewPlugin returns a new Plugin with stdin and stdout.
func NewPlugin() *Plugin {
	var p *Plugin
	p = &Plugin{
		*protogen.NewPlugin(func(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
			pkg, err := p.codeGeneratorRequest(req)
			if err != nil {
				return nil, fmt.Errorf("failed to generate code: %w", err)
			}

			return p.codeGeneratorResponse(pkg)
		}),
	}
	return p
}

func (p *Plugin) codeGeneratorRequest(req *pluginpb.CodeGeneratorRequest) (*v1.Package, error) {
	var resp *v1.Package
	for _, file := range req.GetProtoFile() {
		api := protogen.NewAPIDescriptor(file)
		// discard noisy unused information
		file.SourceCodeInfo = nil

		if len(file.GetService()) == 0 {
			p.Debugf("skipped: no services: %s", file.GetName())
			continue
		}

		debug, _ := protojson.Marshal(file)
		p.Debugf("FileDescriptorProto: %s", debug)

		f, err := p.fileDescriptorProto(req, file, api)
		if err != nil {
			return nil, fmt.Errorf("failed to generate file: %w", err)
		}

		if f != nil {
			resp = f
		}
	}
	return resp, nil
}

func (p *Plugin) fileDescriptorProto(req *pluginpb.CodeGeneratorRequest, file *descriptorpb.FileDescriptorProto, api *protogen.APIDescriptor) (*v1.Package, error) {
	resp := &v1.Package{}

	for _, service := range file.GetService() {
		serviceName := api.ServiceName()
		short := fmt.Sprintf("%s is a CLI for mannaing the %s.", serviceName, serviceName)

		pkg := &v1.Package{
			Package: serviceName,
			Use:     serviceName,
			Short:   short,
			Long:    short,
		}
		apiPackage, err := p.serviceDescriptorProto(service, file, api)
		if err != nil {
			return nil, fmt.Errorf("failed to generate service: %w", err)
		}
		pkg.SubPackages = append(pkg.SubPackages, apiPackage)
		resp.SubPackages = append(resp.SubPackages, pkg)
	}
	return resp, nil
}

var errUnimplemented = fmt.Errorf("todo: implement later")

func (p *Plugin) serviceDescriptorProto(service *descriptorpb.ServiceDescriptorProto, file *descriptorpb.FileDescriptorProto, api *protogen.APIDescriptor) (*v1.Package, error) {
	apiVersion := api.APIVersion()
	serviceName := api.ServiceName()
	short := fmt.Sprintf("%s.%s is a CLI for mannaing the %s.", serviceName, apiVersion, serviceName)
	resp := &v1.Package{
		Package: apiVersion,
		Use:     apiVersion,
		Short:   short,
		Long:    short,
	}
	resources := map[string][]*v1.Command{}
	for _, method := range service.GetMethod() {
		resource, cmd, err := p.methodDescriptorProto(method, file, api)
		if errors.Is(err, errUnimplemented) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to generate method: %w", err)
		}
		if _, ok := resources[resource]; !ok {
			resources[resource] = []*v1.Command{}
		}
		resources[resource] = append(resources[resource], cmd)
	}
	for resource, commands := range resources {
		resp.SubPackages = append(resp.SubPackages, &v1.Package{
			Package:     resource,
			Use:         resource,
			Short:       fmt.Sprintf("%s is a CLI for mannaing the %s.", resource, resource),
			Long:        fmt.Sprintf("%s is a CLI for mannaing the %s.", resource, resource),
			SubCommands: commands,
		})
	}
	return resp, nil
}

func (p *Plugin) methodDescriptorProto(method *descriptorpb.MethodDescriptorProto, file *descriptorpb.FileDescriptorProto, api *protogen.APIDescriptor) (string, *v1.Command, error) {
	m := protogen.NewMethodDescriptor(method)
	switch m.Type() {
	case v1.MethodType_METHOD_TYPE_CREATE:
		resource := strings.ToLower(m.ResourceNameAsCreateMethod())
		cmd, err := p.createCommand(file, method, api)
		if err != nil {
			return "", nil, fmt.Errorf("failed to create command: %w", err)
		}
		return resource, cmd, nil

	case v1.MethodType_METHOD_TYPE_GET:
		return "", nil, errUnimplemented
	}

	return "", nil, fmt.Errorf("unsupported method: %s", method.GetName())
}

func (p *Plugin) createCommand(file *descriptorpb.FileDescriptorProto, method *descriptorpb.MethodDescriptorProto, api *protogen.APIDescriptor) (*v1.Command, error) {
	resource := strings.TrimPrefix(method.GetName(), "Create")
	short := fmt.Sprintf("create is a command to create a new %s", resource)

	req, err := NewRequestMessageDescriptor(file).RequestMessage(*method.InputType)
	if err != nil {
		return nil, fmt.Errorf("failed to create request message: %w", err)
	}

	return &v1.Command{
		Api:           api.ServiceName(),
		ApiVersion:    api.APIVersion(),
		ApiImportPath: api.ImportPath(),
		Package:       strings.ToLower(resource),
		Use:           "create",
		Short:         short,
		Long:          short,
		Method:        method.GetName(),
		MethodType:    v1.MethodType_METHOD_TYPE_CREATE,
		Request:       req.Message,
		StringFlags:   req.StringFlags,
		Int32Flags:    req.Int32Flags,
	}, nil
}

func (p *Plugin) codeGeneratorResponse(pkg *v1.Package) (*pluginpb.CodeGeneratorResponse, error) {
	var resp pluginpb.CodeGeneratorResponse
	if pkg == nil {
		return &resp, nil
	}

	bytes, err := p.MarshalJsonProto(pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal package: %w", err)
	}

	resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String("test.json"),
		Content: proto.String(string(bytes) + "\n"),
	})
	return &resp, nil
}
