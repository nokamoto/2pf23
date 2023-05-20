package ent

import (
	"bytes"
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

func content(t *testing.T, svc *v1.Ent) *string {
	t.Helper()
	m := protojson.MarshalOptions{
		Multiline: true,
	}
	b, err := m.Marshal(svc)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	return proto.String(string(b))
}

func TestNewPlugin(t *testing.T) {
	withEntQuery := func(v *descriptorpb.DescriptorProto) *descriptorpb.DescriptorProto {
		v.Options = &descriptorpb.MessageOptions{}
		proto.SetExtension(v.Options, optionv1.E_Resource_EntQuery, true)
		return v
	}

	testcases := []struct {
		name string
		req  *pluginpb.CodeGeneratorRequest
		want *pluginpb.CodeGeneratorResponse
	}{
		{
			name: "empty",
			req:  &pluginpb.CodeGeneratorRequest{},
			want: &pluginpb.CodeGeneratorResponse{},
		},
		{
			name: "shelf",
			req: &pluginpb.CodeGeneratorRequest{
				Parameter: proto.String("multiline"),
				ProtoFile: []*descriptorpb.FileDescriptorProto{
					{
						Package: proto.String("api.library.v1alpha"),
						Options: &descriptorpb.FileOptions{
							GoPackage: proto.String("github.com/library;libraryv1alpha"),
						},
						MessageType: []*descriptorpb.DescriptorProto{
							withEntQuery(&descriptorpb.DescriptorProto{
								Name: proto.String("Shelf"),
								Field: []*descriptorpb.FieldDescriptorProto{
									{
										JsonName: proto.String("name"),
										Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
									},
									{
										JsonName: proto.String("category"),
										Type:     descriptorpb.FieldDescriptorProto_TYPE_ENUM.Enum(),
									},
									{
										JsonName: proto.String("numBooks"),
										Type:     descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
									},
								},
							}),
						},
					},
				},
			},
			want: &pluginpb.CodeGeneratorResponse{
				File: []*pluginpb.CodeGeneratorResponse_File{
					{
						Name: proto.String("shelf.json"),
						Content: content(t, &v1.Ent{
							Name: "Shelf",
							ImportPath: &v1.ImportPath{
								Alias: "v1alpha",
								Path:  "github.com/library",
							},
							Fields: []string{
								"Name",
								"NumBooks",
							},
							EnumFields: []string{
								"Category",
							},
						}),
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := proto.Marshal(tc.req)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}
			var buf bytes.Buffer
			sut := NewPlugin()
			sut.SetInput(bytes.NewReader(b))
			sut.SetOutput(&buf)
			err = sut.Run()
			if err != nil {
				t.Errorf("got %v, want nil", err)
			}

			var got pluginpb.CodeGeneratorResponse
			if err := proto.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			if diff := cmp.Diff(&got, tc.want, protocmp.Transform()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
