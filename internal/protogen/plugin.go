package protogen

import (
	"fmt"
	"io"
	"os"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin is a plugin for protoc.
type Plugin struct {
	in        io.Reader
	out       io.Writer
	debug     io.Writer
	multiline bool
	handler   func(*pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error)
}

// NewPlugin returns a new Plugin with stdin and stdout.
func NewPlugin(handler func(*pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error)) *Plugin {
	return &Plugin{
		in:        os.Stdin,
		out:       os.Stdout,
		debug:     io.Discard,
		multiline: false,
		handler:   handler,
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

	resp, err := p.handler(&req)
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
	for _, param := range strings.Split(req.GetParameter(), ",") {
		switch param {
		case "debug":
			p.debug = os.Stderr
		case "multiline":
			p.multiline = true
		}
	}
}

// Debugf writes debug messages to stderr if the parameter is "debug".
func (p *Plugin) Debugf(format string, args ...any) {
	fmt.Fprintf(p.debug, "debug: "+format+"\n", args...)
}

// MarshalJSONProto marshals a proto message to JSON in protojson format.
//
// If the parameter is "multiline", the output is formatted with newlines and indentation.
func (p *Plugin) MarshalJsonProto(v proto.Message) ([]byte, error) {
	m := protojson.MarshalOptions{
		Multiline: p.multiline,
	}
	return m.Marshal(v)
}

// SetInput sets the input reader. It is used for testing.
func (p *Plugin) SetInput(in io.Reader) {
	p.in = in
}

// SetOutput sets the output writer. It is used for testing.
func (p *Plugin) SetOutput(out io.Writer) {
	p.out = out
}
