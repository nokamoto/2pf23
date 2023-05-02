package protogen

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestAPIDescriptor(t *testing.T) {
	sut := NewAPIDescriptor(&descriptorpb.FileDescriptorProto{
		Package: proto.String("api.ke.v1alpha"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("com.example.ke.v1alpha;kev1alpha"),
		},
	})

	if got, want := sut.ServiceName(), "ke"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if got, want := sut.APIVersion(), "v1alpha"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	want := &v1.ImportPath{
		Alias: "v1alpha",
		Path:  "com.example.ke.v1alpha",
	}
	if diff := cmp.Diff(sut.ImportPath(), want, protocmp.Transform()); diff != "" {
		t.Error(diff)
	}
}
