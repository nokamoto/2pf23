package protogen

import "strings"

// GoTypeNameFromFullyQualified returns a go type name from a fully qualified type name.
// For example, if typ is `.com.example.v1.CreateRequest`, it returns `v1.CreateRequest`.
// The fully qualified type name is typically from FieldDescriptorProto.TypeName.
func GoTypeNameFromFullyQualified(typ string) string {
	v := strings.Split(typ, ".")
	return strings.Join(v[len(v)-2:], ".")
}
