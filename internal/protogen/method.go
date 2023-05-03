package protogen

import (
	"strings"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/types/descriptorpb"
)

// MethodDescriptor describes a gRPC method from a proto file.
type MethodDescriptor struct {
	method *descriptorpb.MethodDescriptorProto
}

func NewMethodDescriptor(m *descriptorpb.MethodDescriptorProto) *MethodDescriptor {
	return &MethodDescriptor{
		method: m,
	}
}

// Type returns a method type. It determines the type of the method from its name.
func (m *MethodDescriptor) Type() v1.MethodType {
	if strings.HasPrefix(m.method.GetName(), "Create") {
		return v1.MethodType_METHOD_TYPE_CREATE
	}
	if strings.HasPrefix(m.method.GetName(), "Get") {
		return v1.MethodType_METHOD_TYPE_GET
	}
	return v1.MethodType_METHOD_TYPE_UNSPECIFIED
}

// ResourceNameAsCreateMethod returns a resource name from a standard create method name.
// For example, if the method name is `CreateFoo`, it returns `Foo`.
func (m *MethodDescriptor) ResourceNameAsCreateMethod() string {
	return strings.TrimPrefix(m.method.GetName(), "Create")
}
