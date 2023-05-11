package protogen

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func content(t *testing.T, svc *v1.Service) *string {
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

func TestPlugin_Run(t *testing.T) {
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
			name: "create, get, delete, list",
			req: &pluginpb.CodeGeneratorRequest{
				Parameter: proto.String("multiline"),
				ProtoFile: []*descriptorpb.FileDescriptorProto{
					{
						Package: proto.String("api.ke.v1alpha"),
						Options: &descriptorpb.FileOptions{
							GoPackage: proto.String("com.example.ke.v1alpha;kev1alpha"),
						},
						MessageType: []*descriptorpb.DescriptorProto{
							{
								Name: proto.String("ListClusterResponse"),
								Field: []*descriptorpb.FieldDescriptorProto{
									{
										Name:     proto.String("clusters"),
										JsonName: proto.String("clusters"),
									},
									{
										Name: proto.String("next_page_token"),
									},
								},
							},
						},
						Service: []*descriptorpb.ServiceDescriptorProto{
							{
								Name: proto.String("KeService"),
								Method: []*descriptorpb.MethodDescriptorProto{
									{
										Name:       proto.String("CreateCluster"),
										InputType:  proto.String(".api.ke.v1alpha.CreateClusterRequest"),
										OutputType: proto.String(".api.ke.v1alpha.Cluster"),
									},
									{
										Name:       proto.String("GetCluster"),
										InputType:  proto.String(".api.ke.v1alpha.GetClusterRequest"),
										OutputType: proto.String(".api.ke.v1alpha.Cluster"),
									},
									{
										Name:       proto.String("DeleteCluster"),
										InputType:  proto.String(".api.ke.v1alpha.DeleteClusterRequest"),
										OutputType: proto.String(".google.protobuf.Empty"),
									},
									{
										Name:       proto.String("ListCluster"),
										InputType:  proto.String(".api.ke.v1alpha.ListClusterRequest"),
										OutputType: proto.String(".api.ke.v1alpha.ListClusterResponse"),
									},
								},
							},
						},
					},
				},
			},
			want: &pluginpb.CodeGeneratorResponse{
				File: []*pluginpb.CodeGeneratorResponse_File{
					{
						Name: proto.String("ke.v1alpha.json"),
						Content: content(t, &v1.Service{
							Name:       "ke",
							ApiVersion: "v1alpha",
							ApiImportPath: &v1.ImportPath{
								Alias: "v1alpha",
								Path:  "com.example.ke.v1alpha",
							},
							UnimplementedServer: "v1alpha.UnimplementedKeServiceServer",
							Calls: []*v1.Call{
								{
									Name:              "CreateCluster",
									MethodType:        v1.MethodType_METHOD_TYPE_CREATE,
									RequestType:       "v1alpha.CreateClusterRequest",
									ResponseType:      "v1alpha.Cluster",
									ResourceType:      "v1alpha.Cluster",
									GetResourceMethod: "GetCluster",
								},
								{
									Name:         "GetCluster",
									MethodType:   v1.MethodType_METHOD_TYPE_GET,
									RequestType:  "v1alpha.GetClusterRequest",
									ResponseType: "v1alpha.Cluster",
									ResourceType: "v1alpha.Cluster",
								},
								{
									Name:         "DeleteCluster",
									MethodType:   v1.MethodType_METHOD_TYPE_DELETE,
									RequestType:  "v1alpha.DeleteClusterRequest",
									ResponseType: "empty.Empty",
								},
								{
									Name:         "ListCluster",
									MethodType:   v1.MethodType_METHOD_TYPE_LIST,
									RequestType:  "v1alpha.ListClusterRequest",
									ResponseType: "v1alpha.ListClusterResponse",
									ResourceType: "v1alpha.Cluster",
									ListField:    "Clusters",
								},
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
