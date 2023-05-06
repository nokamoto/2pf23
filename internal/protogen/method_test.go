package protogen

import (
	"reflect"
	"testing"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestMethodDescriptor_Type(t *testing.T) {
	tests := []struct {
		name   string
		method *descriptorpb.MethodDescriptorProto
		want   v1.MethodType
	}{
		{
			name: "create",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("CreateFoo"),
			},
			want: v1.MethodType_METHOD_TYPE_CREATE,
		},
		{
			name: "get",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("GetFoo"),
			},
			want: v1.MethodType_METHOD_TYPE_GET,
		},
		{
			name: "delete",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("DeleteFoo"),
			},
			want: v1.MethodType_METHOD_TYPE_DELETE,
		},
		{
			name: "list",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("ListFoo"),
			},
			want: v1.MethodType_METHOD_TYPE_LIST,
		},
		{
			name: "unspecified",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("SearchFoo"),
			},
			want: v1.MethodType_METHOD_TYPE_UNSPECIFIED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMethodDescriptor(tt.method)
			if got := m.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MethodDescriptor.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMethodDescriptor_ResourceName(t *testing.T) {
	tests := []struct {
		name   string
		method *descriptorpb.MethodDescriptorProto
		want   string
	}{
		{
			name: "create",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("CreateFoo"),
			},
			want: "Foo",
		},
		{
			name: "get",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("GetFoo"),
			},
			want: "Foo",
		},
		{
			name: "delete",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("DeleteFoo"),
			},
			want: "Foo",
		},
		{
			name: "list",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("ListFoo"),
			},
			want: "Foo",
		},
		{
			name: "unspecified",
			method: &descriptorpb.MethodDescriptorProto{
				Name: proto.String("SearchFoo"),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMethodDescriptor(tt.method)
			if got := m.ResourceName(); got != tt.want {
				t.Errorf("MethodDescriptor.ResourceName = %v, want %v", got, tt.want)
			}
		})
	}
}
