package ent

import (
	"fmt"
	"strings"

	"github.com/nokamoto/2pf23/internal/protogen"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	optionv1 "github.com/nokamoto/2pf23/pkg/api/option/v1"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin is a protoc plugin to generate ent code.
type Plugin struct {
	protogen.Plugin
}

// NewPlugin returns a new Plugin with stdin and stdout.
func NewPlugin() *Plugin {
	var p *Plugin
	p = &Plugin{
		*protogen.NewPlugin(func(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
			m := map[string]*v1.Ent{}

			for _, file := range req.GetProtoFile() {
				for _, msg := range file.GetMessageType() {
					if !proto.HasExtension(msg.GetOptions(), optionv1.E_Resource_EntQuery) {
						continue
					}
					entQuery := proto.GetExtension(msg.GetOptions(), optionv1.E_Resource_EntQuery).(bool)
					if !entQuery {
						continue
					}

					filename := fmt.Sprintf("%s.json", strings.ToLower(msg.GetName()))
					ent := &v1.Ent{
						Name:       msg.GetName(),
						ImportPath: protogen.NewAPIDescriptor(file).ImportPath(),
					}

					for _, field := range msg.GetField() {
						goName := cases.Title(language.English, cases.NoLower).String(field.GetJsonName())

						switch field.GetType() {
						case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
							ent.EnumFields = append(ent.EnumFields, goName)

						case descriptorpb.FieldDescriptorProto_TYPE_STRING, descriptorpb.FieldDescriptorProto_TYPE_INT32:
							ent.Fields = append(ent.Fields, goName)
						}
					}

					m[filename] = ent
				}
			}

			var resp pluginpb.CodeGeneratorResponse
			for file, ent := range m {
				content, err := p.MarshalJsonProto(ent)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal json: %w", err)
				}
				resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
					Name:    proto.String(file),
					Content: proto.String(string(content)),
				})
			}

			return &resp, nil
		}),
	}
	return p
}
