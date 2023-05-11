package protogen

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/descriptorpb"
)

const nextPageTokenField = "nextPageToken"

var (
	errInvalidNumberOfFields = fmt.Errorf("the response message must have two fields")
	errListFieldNotFound     = fmt.Errorf("the list field not found")
)

type ListResponseDescriptor struct {
	file   *descriptorpb.FileDescriptorProto
	method *descriptorpb.MethodDescriptorProto
}

func NewListResponseDescriptor(f *descriptorpb.FileDescriptorProto, m *descriptorpb.MethodDescriptorProto) *ListResponseDescriptor {
	return &ListResponseDescriptor{f, m}
}

// ListField returns the field name of the resource list.
//
// The response message of the list method must have two fields.
// One is the resource list and the other is the next page token. If the response message has more than two fields, it returns an error.
//
// For example, the following response message has two fields:
// ```
//
//	message ListClusterResponse {
//	  repeated Cluster clusters = 1;
//	  string next_page_token = 2;
//	}
//
// ```
// In this case, it returns `Clusters`.
func (l *ListResponseDescriptor) ListField() (string, error) {
	var res string
	for _, typ := range l.file.GetMessageType() {
		if strings.HasSuffix(l.method.GetOutputType(), fmt.Sprintf(".%s", typ.GetName())) {
			if size := len(typ.GetField()); size != 2 {
				return "", fmt.Errorf("%w: %d field(s) in %s", errInvalidNumberOfFields, size, l.method.GetOutputType())
			}
			for _, field := range typ.GetField() {
				if name := field.GetJsonName(); name != nextPageTokenField {
					res = name
					break
				}
			}
			break
		}
	}
	if res == "" {
		return "", fmt.Errorf("%w: %s", errListFieldNotFound, l.method.GetOutputType())
	}
	return cases.Title(language.English, cases.NoLower).String(res), nil
}
