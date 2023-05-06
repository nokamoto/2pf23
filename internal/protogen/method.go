package protogen

import (
	"strings"

	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	createPrefix = "Create"
	getPrefix    = "Get"
	deletePrefix = "Delete"
	listPrefix   = "List"
)

var prefixes = map[string]v1.MethodType{
	createPrefix: v1.MethodType_METHOD_TYPE_CREATE,
	getPrefix:    v1.MethodType_METHOD_TYPE_GET,
	deletePrefix: v1.MethodType_METHOD_TYPE_DELETE,
	listPrefix:   v1.MethodType_METHOD_TYPE_LIST,
}

// MethodDescriptor describes a gRPC method from a proto file.
type MethodDescriptor struct {
	*descriptorpb.MethodDescriptorProto
}

func NewMethodDescriptor(m *descriptorpb.MethodDescriptorProto) *MethodDescriptor {
	return &MethodDescriptor{
		MethodDescriptorProto: m,
	}
}

// Type returns a method type. It determines the type of the method from its name.
func (m *MethodDescriptor) Type() v1.MethodType {
	for prefix, methodType := range prefixes {
		if strings.HasPrefix(m.GetName(), prefix) {
			return methodType
		}
	}
	return v1.MethodType_METHOD_TYPE_UNSPECIFIED
}

// ResourceName returns a resource name from a method type.
// For example, if the method type is `Create` and the method name is `CreateFoo`, it returns `Foo`.
func (m *MethodDescriptor) ResourceName() string {
	typ := m.Type()
	for prefix, methodType := range prefixes {
		if methodType == typ {
			return strings.TrimPrefix(m.GetName(), prefix)
		}
	}
	return ""
}
