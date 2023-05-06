package helper

import (
	"encoding/base64"
	"fmt"

	"github.com/nokamoto/2pf23/internal/app"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/proto"
)

type listRequest interface {
	GetPageToken() string
}

// Pagination returns a pagination object from the list request.
// If the page token is empty, it returns nil. It means that the first page is requested.
// If the page token is invalid, it returns app.ErrInvalidArgument.
func Pagination(req listRequest) (*v1.Pagination, error) {
	if req.GetPageToken() == "" {
		return nil, nil
	}
	var res v1.Pagination
	b, err := base64.StdEncoding.DecodeString(req.GetPageToken())
	if err != nil {
		return nil, fmt.Errorf("%w: invalid page token: %s", app.ErrInvalidArgument, err)
	}
	if err := proto.Unmarshal(b, &res); err != nil {
		return nil, fmt.Errorf("%w: invalid page token: %s", app.ErrInvalidArgument, err)
	}
	return &res, nil
}

// PageToken returns a page token from the pagination object.
// If the pagination object is nil, it returns an empty string. It means that the last page is returned.
func PageToken(p *v1.Pagination) (string, error) {
	if p == nil {
		return "", nil
	}
	b, err := proto.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("failed to marshal pagination: %s", err)
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
