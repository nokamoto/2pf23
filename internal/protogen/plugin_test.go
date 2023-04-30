package protogen

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func in(v *pluginpb.CodeGeneratorRequest) io.Reader {
	b, _ := proto.Marshal(v)
	return bytes.NewReader(b)
}

func TestPlugin_Debugf(t *testing.T) {
	p := NewPlugin(func(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
		return &pluginpb.CodeGeneratorResponse{}, nil
	})
	p.in = in(&pluginpb.CodeGeneratorRequest{
		Parameter: proto.String("debug"),
	})
	p.out = io.Discard

	// debug is io.Discard by default
	if p.debug != io.Discard {
		t.Error("debug writer is not io.Discard")
	}

	// debug is os.Stderr if the parameter is "debug
	if err := p.Run(); err != nil {
		t.Fatal(err)
	}
	if p.debug != os.Stderr {
		t.Error("debug writer is not stderr")
	}

	// debugf writes debug messages to debug writer
	var buf bytes.Buffer
	p.debug = &buf
	p.Debugf("hello %s", "world")
	if got, want := buf.String(), "debug: hello world\n"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPlugin_MarshalJsonProto(t *testing.T) {
	p := NewPlugin(func(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
		return &pluginpb.CodeGeneratorResponse{}, nil
	})
	p.in = in(&pluginpb.CodeGeneratorRequest{
		Parameter: proto.String("multiline"),
	})
	p.out = io.Discard

	// multiline is false by default
	if p.multiline != false {
		t.Error("multiline is not false")
	}
	got, err := p.MarshalJsonProto(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{"foo.proto"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `{"fileToGenerate":["foo.proto"]}`
	if diff := cmp.Diff([]byte(want), got); diff != "" {
		t.Errorf("got %q, want %q", got, want)
	}

	// multiline is true if the parameter is "multiline"
	if err := p.Run(); err != nil {
		t.Fatal(err)
	}
	if p.multiline != true {
		t.Error("multiline is not true")
	}
}
