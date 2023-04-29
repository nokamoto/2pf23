package protogen

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	optionv1 "github.com/nokamoto/2pf23/pkg/api/option/v1"
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

func withUsage(field *descriptorpb.FieldDescriptorProto, usage string) *descriptorpb.FieldDescriptorProto {
	field.Options = &descriptorpb.FieldOptions{}
	proto.SetExtension(field.Options, optionv1.E_Resource_Usage, usage)
	return field
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
						Package: proto.String("api.library.v1alpha"),
						Options: &descriptorpb.FileOptions{
							GoPackage: proto.String("github.com/library;libraryv1alpha"),
						},
						MessageType: []*descriptorpb.DescriptorProto{
							{
								Name: proto.String("Shelf"),
								Field: []*descriptorpb.FieldDescriptorProto{
									withUsage(
										&descriptorpb.FieldDescriptorProto{
											Name:     proto.String("display_name"),
											JsonName: proto.String("displayName"),
											Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
										},
										"display name usage",
									),
								},
							},
							{
								Name: proto.String("CreateShelfRequest"),
								Field: []*descriptorpb.FieldDescriptorProto{
									{
										Name:     proto.String("shelf"),
										JsonName: proto.String("shelf"),
										Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
										TypeName: proto.String(".api.library.v1alpha.Shelf"),
									},
								},
							},
						},
						Service: []*descriptorpb.ServiceDescriptorProto{
							{
								Method: []*descriptorpb.MethodDescriptorProto{
									{
										Name:      proto.String("CreateShelf"),
										InputType: proto.String(".api.library.v1alpha.CreateShelfRequest"),
									},
								},
							},
						},
					},
				},
			},
			expected: content(t, &v1.Package{
				SubPackages: []*v1.Package{
					{
						Package: "library",
						Use:     "library",
						Short:   "library is a CLI for mannaing the library.",
						Long:    "library is a CLI for mannaing the library.",
						SubPackages: []*v1.Package{
							{
								Package: "v1alpha",
								Use:     "v1alpha",
								Short:   "library.v1alpha is a CLI for mannaing the library.",
								Long:    "library.v1alpha is a CLI for mannaing the library.",
								SubPackages: []*v1.Package{
									{
										Package: "shelf",
										Use:     "shelf",
										Short:   "shelf is a CLI for mannaing the shelf.",
										Long:    "shelf is a CLI for mannaing the shelf.",
										SubCommands: []*v1.Command{
											{
												Api:        "library",
												ApiVersion: "v1alpha",
												ApiImportPath: &v1.ImportPath{
													Alias: "v1alpha",
													Path:  "github.com/library",
												},
												Package:    "shelf",
												Use:        "create",
												Short:      "create is a command to create a new Shelf",
												Long:       "create is a command to create a new Shelf",
												Method:     "CreateShelf",
												MethodType: v1.MethodType_METHOD_TYPE_CREATE,
												Request: &v1.RequestMessage{
													Type: "v1alpha.CreateShelfRequest",
													Children: []*v1.RequestMessage{
														{
															Name: "Shelf",
															Type: "v1alpha.Shelf",
															Fields: []*v1.RequestMessageField{
																{
																	Name:  "DisplayName",
																	Value: "displayName",
																},
															},
														},
													},
												},
												StringFlags: []*v1.Flag{
													{
														Name:        "displayName",
														DisplayName: "display-name",
														Usage:       "display name usage",
													},
												},
											},
										},
									},
								},
							},
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
				in:        bytes.NewBuffer(b),
				out:       &out,
				debug:     io.Discard,
				multiline: true,
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
