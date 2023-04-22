package cligen

import (
	"fmt"
	"io"
	"os"
	"strings"

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
		// discard noisy unused information
		file.SourceCodeInfo = nil

		if len(file.GetService()) == 0 {
			p.debugf("skipped: no services: %s", file.GetName())
			continue
		}

		debug, _ := protojson.Marshal(file)
		p.debugf("FileDescriptorProto: %s", debug)

		f, err := p.fileDescriptorProto(req, file)
		if err != nil {
			return nil, fmt.Errorf("failed to generate file: %w", err)
		}

		if f != nil {
			resp = f
		}
	}
	return resp, nil
}

func (p *Plugin) fileDescriptorProto(req *pluginpb.CodeGeneratorRequest, file *descriptorpb.FileDescriptorProto) (*v1.Package, error) {
	var resp *v1.Package
	for i, service := range file.GetService() {
		p.debugf("[%d] ServiceDescriptorProto: %s", i, service.GetName())

		for j, method := range service.GetMethod() {
			p.debugf("[%d/%d] MethodDescriptorProto: %s", i, j, method.GetName())

			if strings.HasPrefix(method.GetName(), "Create") {
				resource := strings.TrimPrefix(method.GetName(), "Create")
				lowercaseResource := strings.ToLower(resource)

				short := fmt.Sprintf("%s is a CLI for mannaing the %s.", lowercaseResource, resource)

				sub, err := p.createCommand(file, method)
				if err != nil {
					return nil, fmt.Errorf("failed to create command: %w", err)
				}

				resp = &v1.Package{
					Package:     lowercaseResource,
					Use:         lowercaseResource,
					Short:       short,
					Long:        short,
					SubCommands: []*v1.Command{sub},
					SubPackages: []*v1.Package{},
				}
				continue
			}

			return nil, fmt.Errorf("unsupported method: %s", method.GetName())
		}
	}
	return resp, nil
}

func serviceFromPackage(file *descriptorpb.FileDescriptorProto) string {
	v := strings.Split(file.GetPackage(), ".")
	return v[len(v)-2]
}

func (p *Plugin) children(typ string, file *descriptorpb.FileDescriptorProto) (*v1.Resource, error) {
	var fields []*v1.ResourceField
	for _, message := range file.GetMessageType() {
		if strings.HasSuffix(typ, message.GetName()) {
			for _, field := range message.GetField() {
				if field.GetName() == "name" {
					// `name` is output only field
					continue
				}
				fields = append(fields, &v1.ResourceField{
					Id:       cases.Title(language.English, cases.NoLower).String(*field.JsonName),
					FlagName: *field.JsonName,
				})
			}
			return &v1.Resource{
				Type:   typ,
				Fields: fields,
			}, nil
		}
	}
	return nil, fmt.Errorf("not found: %s in %s", typ, file.GetName())
}

func (p *Plugin) createCommand(file *descriptorpb.FileDescriptorProto, method *descriptorpb.MethodDescriptorProto) (*v1.Command, error) {
	first := func(s []string) string {
		return s[0]
	}
	last := func(s []string) string {
		return s[len(s)-1]
	}

	pkg := last(strings.Split(file.GetPackage(), "."))
	resource := strings.TrimPrefix(method.GetName(), "Create")
	short := fmt.Sprintf("create is a command to create a new %s", resource)

	child, err := p.children(fmt.Sprintf("%s.%s", pkg, resource), file)
	if err != nil {
		return nil, fmt.Errorf("failed to generate child: %w", err)
	}

	return &v1.Command{
		Api:        serviceFromPackage(file),
		ApiVersion: pkg,
		ApiImportPath: &v1.ImportPath{
			Alias: pkg,
			Path:  first(strings.Split(file.GetOptions().GetGoPackage(), ";")),
		},
		Package:          strings.ToLower(resource),
		Use:              "create",
		Short:            short,
		Long:             short,
		Method:           "Create",
		MethodType:       v1.MethodType_METHOD_TYPE_CREATE,
		CreateResourceId: resource,
		CreateResource: &v1.Resource{
			Type:     fmt.Sprintf("%s.%s", pkg, resource),
			Fields:   []*v1.ResourceField{},
			Children: []*v1.Resource{child},
		},
		StringFlags: []*v1.Flag{},
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
