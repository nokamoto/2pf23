package cli

import (
	"fmt"
	"strings"

	"github.com/nokamoto/2pf23/internal/protogen/core"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	optionv1 "github.com/nokamoto/2pf23/pkg/api/option/v1"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// RequestMessageDescriptor describes a construct of a request message from a proto file.
// It describes relationships between flags and a request message.
type RequestMessageDescriptor struct {
	file *descriptorpb.FileDescriptorProto

	// StringFlags is a list of string flags.
	StringFlags []*v1.Flag
	// Int32Flags is a list of number flags.
	Int32Flags []*v1.Flag
	// EnumFlags is a list of enum flags.
	EnumFlags []*v1.EnumFlag
	// Message is a construct of a request message.
	Message *v1.RequestMessage
}

func NewRequestMessageDescriptor(file *descriptorpb.FileDescriptorProto) *RequestMessageDescriptor {
	return &RequestMessageDescriptor{file: file}
}

func (r *RequestMessageDescriptor) requestMessage(typ string, name string, withName bool, withPath bool) (*v1.RequestMessage, error) {
	resp := &v1.RequestMessage{
		Type: core.GoTypeNameFromFullyQualified(typ),
		Name: name,
	}

	found := false
	for _, message := range r.file.GetMessageType() {
		goType := fmt.Sprintf(".%s.%s", r.file.GetPackage(), message.GetName())
		if found = typ == goType; found {
			for _, field := range message.GetField() {
				if field.GetName() == "name" {
					if withName {
						resp.Fields = append(resp.Fields, &v1.RequestMessageField{
							Name:  "Name",
							Value: "args[0]",
						})
					}
					continue
				}
				if field.GetName() == "update_mask" {
					resp.Fields = append(resp.Fields, &v1.RequestMessageField{
						Name:  "UpdateMask",
						Value: "mask",
					})
					continue
				}

				flag := &v1.Flag{
					Name:        field.GetJsonName(),
					DisplayName: strings.ReplaceAll(field.GetName(), "_", "-"),
					Value:       "",
				}
				if proto.HasExtension(field.GetOptions(), optionv1.E_Resource_Usage) {
					flag.Usage = proto.GetExtension(field.GetOptions(), optionv1.E_Resource_Usage).(string)
				}
				if withPath {
					flag.Path = field.GetName()
				}

				goFieldName := cases.Title(language.English, cases.NoLower).String(field.GetJsonName())

				goField := &v1.RequestMessageField{
					Name:  goFieldName,
					Value: flag.Name,
				}

				switch field.GetType() {
				case descriptorpb.FieldDescriptorProto_TYPE_INT32:
					flag.Value = "0"
					resp.Fields = append(resp.Fields, goField)
					r.Int32Flags = append(r.Int32Flags, flag)

				case descriptorpb.FieldDescriptorProto_TYPE_STRING:
					flag.Value = ""
					resp.Fields = append(resp.Fields, goField)
					r.StringFlags = append(r.StringFlags, flag)

				case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
					resp.Fields = append(resp.Fields, &v1.RequestMessageField{
						Name:  goFieldName,
						Value: fmt.Sprintf("%s.Value()", flag.Name),
					})
					r.EnumFlags = append(r.EnumFlags, &v1.EnumFlag{
						Name:        flag.GetName(),
						DisplayName: flag.GetDisplayName(),
						Type:        core.GoTypeNameFromFullyQualified(field.GetTypeName()),
						Usage:       flag.GetUsage(),
					})

				case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
					sub, err := r.requestMessage(field.GetTypeName(), goFieldName, withName, withPath)
					if err != nil {
						return nil, fmt.Errorf("failed to create field request message: %w", err)
					}
					resp.Children = append(resp.Children, sub)

				default:
					return nil, fmt.Errorf("unsupported field type: %s", field.GetType())
				}
			}
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("failed to find message: %s", typ)
	}

	return resp, nil
}

// RequestMessage returns a request message descriptor. It sets and returns itself.
//
// typ is a type of the request message. It is resolved in the same way as FieldDescriptorProto.type_name, but must refer to a message type.
// (e.g. ".com.example.FooRequest")
func (r *RequestMessageDescriptor) RequestMessage(typ string, withName bool, withPath bool) (*RequestMessageDescriptor, error) {
	res, err := r.requestMessage(typ, "", withName, withPath)
	if err != nil {
		r.Int32Flags = nil
		r.StringFlags = nil
		return nil, err
	}
	r.Message = res
	return r, nil
}
