package cligen

import (
	"fmt"
	"io"
	"os"
	"strings"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type Plugin struct {
	in    io.Reader
	out   io.Writer
	debug io.Writer
}

// NewPlugin returns a new Plugin with stdin and stdout.
func NewPlugin() *Plugin {
	return &Plugin{
		in:    os.Stdin,
		out:   os.Stdout,
		debug: io.Discard,
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

	packages, err := p.codeGeneratorRequest(&req)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	bytes, err = proto.Marshal(p.codeGeneratorResponse(packages))
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

func (p *Plugin) codeGeneratorRequest(req *pluginpb.CodeGeneratorRequest) ([]*v1.Package, error) {
	var resp []*v1.Package
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
			resp = append(resp, f...)
		}
	}

	return resp, nil
}

func (p *Plugin) fileDescriptorProto(req *pluginpb.CodeGeneratorRequest, file *descriptorpb.FileDescriptorProto) ([]*v1.Package, error) {
	var resp []*v1.Package
	for i, service := range file.GetService() {
		p.debugf("[%d] ServiceDescriptorProto: %s", i, service.GetName())

		for j, method := range service.GetMethod() {
			p.debugf("[%d/%d] MethodDescriptorProto: %s", i, j, method.GetName())

			if strings.HasPrefix(method.GetName(), "Create") {
				resp = append(resp, &v1.Package{
					Package:     "test1",
					ImportPath:  "test2",
					Use:         "test3",
					Short:       "test4",
					Long:        "test5",
					SubCommands: []*v1.Command{},
					SubPackages: []*v1.Package{},
				})
				continue
			}

			panic("not implemented")
		}
	}
	return resp, nil
}

func (p *Plugin) codeGeneratorResponse(packages []*v1.Package) *pluginpb.CodeGeneratorResponse {
	var resp pluginpb.CodeGeneratorResponse
	for _, pkg := range packages {
		bytes, err := protojson.Marshal(pkg)
		if err != nil {
			panic(err)
		}

		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(pkg.GetPackage() + ".json"),
			Content: proto.String(string(bytes)),
		})
	}
	return &resp
}
