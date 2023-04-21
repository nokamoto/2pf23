package cligen

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func content(t *testing.T, content *v1.Package) *pluginpb.CodeGeneratorResponse {
	t.Helper()

	m := protojson.MarshalOptions{
		Multiline: true,
	}
	b, err := m.Marshal(content)
	if err != nil {
		t.Fatalf("failed to marshal content: %v", err)
	}

	return &pluginpb.CodeGeneratorResponse{
		File: []*pluginpb.CodeGeneratorResponse_File{
			{
				Name:    proto.String("test.json"),
				Content: proto.String(string(b) + "\n"),
			},
		},
	}
}

func TestPlugin_Run(t *testing.T) {
	testcases := []struct {
		name     string
		req      *pluginpb.CodeGeneratorRequest
		expected *pluginpb.CodeGeneratorResponse
	}{
		{
			name:     "empty",
			req:      &pluginpb.CodeGeneratorRequest{},
			expected: &pluginpb.CodeGeneratorResponse{},
		},
		{
			name: "create",
			req: &pluginpb.CodeGeneratorRequest{
				ProtoFile: []*descriptorpb.FileDescriptorProto{
					{
						Package: proto.String("api.foo.v1alpha"),
						Options: &descriptorpb.FileOptions{
							GoPackage: proto.String("github.com/foo;bar"),
						},
						Service: []*descriptorpb.ServiceDescriptorProto{
							{
								Method: []*descriptorpb.MethodDescriptorProto{
									{
										Name: proto.String("CreateFoo"),
									},
								},
							},
						},
					},
				},
			},
			expected: content(t, &v1.Package{
				Package: "foo",
				Use:     "foo",
				Short:   "foo is a CLI for mannaing the Foo.",
				Long:    "foo is a CLI for mannaing the Foo.",
				SubCommands: []*v1.Command{
					{
						Package:    "foo",
						Api:        "foo",
						ApiVersion: "v1alpha",
						ApiImportPath: &v1.ImportPath{
							Alias: "v1alpha",
							Path:  "github.com/foo",
						},
						Use:              "create",
						Short:            "create is a command to create a new Foo",
						Long:             "create is a command to create a new Foo",
						Method:           "Create",
						MethodType:       v1.MethodType_METHOD_TYPE_CREATE,
						CreateResourceId: "Foo",
						CreateResource: &v1.Resource{
							Type: "v1alpha.Foo",
						},
					},
				},
			}),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := proto.Marshal(tc.req)
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

			var actual pluginpb.CodeGeneratorResponse
			if err := proto.Unmarshal(out.Bytes(), &actual); err != nil {
				t.Fatalf("failed to unmarshal output: %v", err)
			}

			if diff := cmp.Diff(tc.expected, &actual, protocmp.Transform()); diff != "" {
				t.Errorf("unexpected output (-want +got):\n%s", diff)
			}
		})
	}
}
