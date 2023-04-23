package protogen

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nokamoto/2pf23/internal/cligen/protogen/api"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type Plugin struct {
	in        io.Reader
	out       io.Writer
	debug     io.Writer
	multiline bool
}

// NewPlugin returns a new Plugin with stdin and stdout.
func NewPlugin() *Plugin {
	return &Plugin{
		in:        os.Stdin,
		out:       os.Stdout,
		debug:     io.Discard,
		multiline: true,
	}
}

// Run reads CodeGeneratorRequest from stdin, writes CodeGeneratorResponse to stdout.
//
// if the parameter is "debug", it writes debug messages to stderr.
func (p *Plugin) Run() error {
	bytes, err := io.ReadAll(p.in)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	var req pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(bytes, &req); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	p.setParam(&req)

	pkg, err := p.codeGeneratorRequest(&req)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	resp, err := p.codeGeneratorResponse(pkg)
	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}
	bytes, err = proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}

	if _, err := p.out.Write(bytes); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	return nil
}

func (p *Plugin) setParam(req *pluginpb.CodeGeneratorRequest) {
	if req.GetParameter() == "debug" {
		p.debug = os.Stderr
	}
}

func (p *Plugin) debugf(format string, args ...any) {
	fmt.Fprintf(p.debug, "debug: "+format+"\n", args...)
}

func (p *Plugin) codeGeneratorRequest(req *pluginpb.CodeGeneratorRequest) (*v1.Package, error) {
	var resp *v1.Package
	for _, file := range req.GetProtoFile() {
		apidesc := api.NewAPI(file)
		// discard noisy unused information
		file.SourceCodeInfo = nil

		if len(file.GetService()) == 0 {
			p.debugf("skipped: no services: %s", file.GetName())
			continue
		}

		debug, _ := protojson.Marshal(file)
		p.debugf("FileDescriptorProto: %s", debug)

		f, err := p.fileDescriptorProto(req, file, apidesc)
		if err != nil {
			return nil, fmt.Errorf("failed to generate file: %w", err)
		}

		if f != nil {
			resp = f
		}
	}
	return resp, nil
}

func (p *Plugin) fileDescriptorProto(req *pluginpb.CodeGeneratorRequest, file *descriptorpb.FileDescriptorProto, apidesc *api.API) (*v1.Package, error) {
	resp := &v1.Package{}

	for _, service := range file.GetService() {
		serviceName := apidesc.ServiceName()
		short := fmt.Sprintf("%s is a CLI for mannaing the %s.", serviceName, serviceName)

		pkg := &v1.Package{
			Package: serviceName,
			Use:     serviceName,
			Short:   short,
			Long:    short,
		}
		apiPackage, err := p.serviceDescriptorProto(service, file, apidesc)
		if err != nil {
			return nil, fmt.Errorf("failed to generate service: %w", err)
		}
		pkg.SubPackages = append(pkg.SubPackages, apiPackage)
		resp.SubPackages = append(resp.SubPackages, pkg)
	}
	return resp, nil
}

func (p *Plugin) serviceDescriptorProto(service *descriptorpb.ServiceDescriptorProto, file *descriptorpb.FileDescriptorProto, apidesc *api.API) (*v1.Package, error) {
	apiVersion := apidesc.APIVersion()
	serviceName := apidesc.ServiceName()
	short := fmt.Sprintf("%s.%s is a CLI for mannaing the %s.", serviceName, apiVersion, serviceName)
	resp := &v1.Package{
		Package: apiVersion,
		Use:     apiVersion,
		Short:   short,
		Long:    short,
	}
	resources := map[string][]*v1.Command{}
	for _, method := range service.GetMethod() {
		resource, cmd, err := p.methodDescriptorProto(method, file, apidesc)
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

func (p *Plugin) methodDescriptorProto(method *descriptorpb.MethodDescriptorProto, file *descriptorpb.FileDescriptorProto, apidesc *api.API) (string, *v1.Command, error) {
	if strings.HasPrefix(method.GetName(), "Create") {
		resource := strings.ToLower(strings.TrimPrefix(method.GetName(), "Create"))
		cmd, err := p.createCommand(file, method, apidesc)
		if err != nil {
			return "", nil, fmt.Errorf("failed to create command: %w", err)
		}
		return resource, cmd, nil
	}

	return "", nil, fmt.Errorf("unsupported method: %s", method.GetName())
}

func (p *Plugin) requestMessage(typ string, name string, file *descriptorpb.FileDescriptorProto) (*v1.RequestMessage, []*v1.Flag, error) {
	goType := strings.Split(typ, ".")
	goType = goType[len(goType)-2:]
	resp := &v1.RequestMessage{
		Type: strings.Join(goType, "."),
		Name: name,
	}

	var flags []*v1.Flag

	found := false
	for _, message := range file.GetMessageType() {
		absoluteName := fmt.Sprintf(".%s.%s", file.GetPackage(), message.GetName())
		if found = typ == absoluteName; found {
			for _, field := range message.GetField() {
				if field.GetName() == "name" {
					// `name` is output only field
					continue
				}
				switch field.GetType() {
				case descriptorpb.FieldDescriptorProto_TYPE_STRING:
					resp.Fields = append(resp.Fields, &v1.RequestMessageField{
						Name:  cases.Title(language.English, cases.NoLower).String(*field.JsonName),
						Value: *field.JsonName,
					})
					flags = append(flags, &v1.Flag{
						Name:        *field.JsonName,
						DisplayName: strings.ReplaceAll(field.GetName(), "_", "-"),
						Value:       "",
						Usage:       "todo",
					})

				case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
					sub, fs, err := p.requestMessage(field.GetTypeName(), cases.Title(language.English, cases.NoLower).String(field.GetJsonName()), file)
					if err != nil {
						return nil, nil, fmt.Errorf("failed to create field request message: %w", err)
					}
					resp.Children = append(resp.Children, sub)
					flags = append(flags, fs...)

				default:
					return nil, nil, fmt.Errorf("unsupported field type: %s", field.GetType())
				}
			}
			break
		}
	}

	if !found {
		return nil, nil, fmt.Errorf("failed to find message: %s", typ)
	}

	return resp, flags, nil
}

func (p *Plugin) createCommand(file *descriptorpb.FileDescriptorProto, method *descriptorpb.MethodDescriptorProto, apidesc *api.API) (*v1.Command, error) {
	apiVersion := apidesc.APIVersion()
	resource := strings.TrimPrefix(method.GetName(), "Create")
	short := fmt.Sprintf("create is a command to create a new %s", resource)
	req, flags, err := p.requestMessage(*method.InputType, "", file)
	if err != nil {
		return nil, fmt.Errorf("failed to create request message: %w", err)
	}

	return &v1.Command{
		Api:           apidesc.ServiceName(),
		ApiVersion:    apiVersion,
		ApiImportPath: apidesc.ImportPath(),
		Package:       strings.ToLower(resource),
		Use:           "create",
		Short:         short,
		Long:          short,
		Method:        method.GetName(),
		MethodType:    v1.MethodType_METHOD_TYPE_CREATE,
		Request:       req,
		StringFlags:   flags,
	}, nil
}

func (p *Plugin) codeGeneratorResponse(pkg *v1.Package) (*pluginpb.CodeGeneratorResponse, error) {
	var resp pluginpb.CodeGeneratorResponse
	if pkg == nil {
		return &resp, nil
	}

	m := protojson.MarshalOptions{
		Multiline: p.multiline,
	}
	bytes, err := m.Marshal(pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal package: %w", err)
	}

	resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String("test.json"),
		Content: proto.String(string(bytes) + "\n"),
	})
	return &resp, nil
}
