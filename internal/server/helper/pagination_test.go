package helper

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/2pf23/internal/app"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestPagination(t *testing.T) {
	token, err := PageToken(&v1.Pagination{
		Cursor: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		req  listRequest
		want *v1.Pagination
		err  error
	}{
		{
			name: "empty",
			req:  &kev1alpha.ListClusterRequest{},
		},
		{
			name: "ok",
			req: &kev1alpha.ListClusterRequest{
				PageToken: token,
			},
			want: &v1.Pagination{
				Cursor: 100,
			},
		},
		{
			name: "invalid page token not base64 encoded",
			req: &kev1alpha.ListClusterRequest{
				PageToken: "invalid",
			},
			err: app.ErrInvalidArgument,
		},
		{
			name: "invalid page token not proto",
			req: &kev1alpha.ListClusterRequest{
				PageToken: base64.StdEncoding.EncodeToString([]byte("invalid")),
			},
			err: app.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Pagination(tt.req)
			if !errors.Is(err, tt.err) {
				t.Errorf("Pagination() error = %v, wantErr %v", err, tt.err)
			}
			if diff := cmp.Diff(got, tt.want, protocmp.Transform()); diff != "" {
				t.Errorf("Pagination() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestPageToken(t *testing.T) {
	assert := func(p *v1.Pagination, want string) {
		t.Helper()
		token, err := PageToken(p)
		if err != nil {
			t.Fatal(err)
		}
		if token != want {
			t.Errorf("PageToken() = %v, want %v", token, want)
		}
	}
	assert(nil, "")
	assert(&v1.Pagination{Cursor: 100}, "CGQ=")
}
