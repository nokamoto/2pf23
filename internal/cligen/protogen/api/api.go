package api

import (
	"strings"

	"github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/types/descriptorpb"
)

// API describes a gRPC API from a proto file.
type API struct {
	file *descriptorpb.FileDescriptorProto
}

func NewAPI(f *descriptorpb.FileDescriptorProto) *API {
	return &API{f}
}

func (a *API) lastOfPackage(i int) string {
	v := strings.Split(a.file.GetPackage(), ".")
	return v[len(v)-i]
}

// ServiceName returns the service name from the package name.
// For example, if the package name is "com.example.v1", it returns "example".
func (a *API) ServiceName() string {
	return a.lastOfPackage(2)
}

// APIVersion returns the API version from the package name.
// For example, if the package name is "com.example.v1", it returns "v1".
func (a *API) APIVersion() string {
	return a.lastOfPackage(1)
}

// ImportPath returns the import path from the go_package option.
// For example, if the go_package option is "com.example.v1;example", it returns "com.example.v1".
func (a *API) ImportPath() *v1.ImportPath {
	return &v1.ImportPath{
		Alias: a.APIVersion(),
		Path:  strings.Split(a.file.GetOptions().GetGoPackage(), ";")[0],
	}
}
