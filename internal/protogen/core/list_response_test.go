package core

import (
	"errors"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestListResponseDescriptor_ListField(t *testing.T) {
	tt := []struct {
		name   string
		file   *descriptorpb.FileDescriptorProto
		method *descriptorpb.MethodDescriptorProto
		want   string
		err    error
	}{
		{
			name: "success",
			file: &descriptorpb.FileDescriptorProto{
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
			},
			method: &descriptorpb.MethodDescriptorProto{
				OutputType: proto.String(".test.ListClusterResponse"),
			},
			want: "Clusters",
		},
		{
			name: "not found",
			method: &descriptorpb.MethodDescriptorProto{
				OutputType: proto.String(".test.ListClusterResponse"),
			},
			err: errListFieldNotFound,
		},
		{
			name: "invalid message type",
			file: &descriptorpb.FileDescriptorProto{
				MessageType: []*descriptorpb.DescriptorProto{
					{
						Name: proto.String("ListClusterResponse"),
						Field: []*descriptorpb.FieldDescriptorProto{
							{
								Name:     proto.String("clusters"),
								JsonName: proto.String("clusters"),
							},
						},
					},
				},
			},
			method: &descriptorpb.MethodDescriptorProto{
				OutputType: proto.String(".test.ListClusterResponse"),
			},
			err: errInvalidNumberOfFields,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			l := NewListResponseDescriptor(tc.file, tc.method)
			got, err := l.ListField()
			if !errors.Is(err, tc.err) {
				t.Errorf("ListField() error = %v, want %v", err, tc.err)
			}
			if got != tc.want {
				t.Errorf("ListField() got = %v, want %v", got, tc.want)
			}
		})
	}
}
