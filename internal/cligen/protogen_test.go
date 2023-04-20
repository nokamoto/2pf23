package cligen

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestPlugin_Run(t *testing.T) {
	in := &pluginpb.CodeGeneratorRequest{}
	b, err := proto.Marshal(in)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}

	var out bytes.Buffer
	p := Plugin{
		in:    bytes.NewBuffer(b),
		out:   &out,
		debug: io.Discard,
	}
	if err := p.Run(); err != nil {
		t.Fatalf("failed to run plugin: %v", err)
	}

	expected := &pluginpb.CodeGeneratorResponse{}

	var actual pluginpb.CodeGeneratorResponse
	if err := proto.Unmarshal(out.Bytes(), &actual); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if diff := cmp.Diff(expected, &actual, protocmp.Transform()); diff != "" {
		t.Errorf("unexpected output (-want +got):\n%s", diff)
	}
}
