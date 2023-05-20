package helper

import (
	"context"

	"github.com/bufbuild/connect-go"
	"github.com/nokamoto/2pf23/internal/util/helper"
)

// ListAll is a helper function to list all items from a paginated API.
//
// list is the function to call to list items. (e.g. kev1alpha.KeService.ListCluster)
// f1 is the function to create the request object. It takes a page token as input and returns the request message (e.g. func(string) *kev1alpha.ListClusterRequest)
// f2 is the function to process the response object. It takes the response message as input. (e.g. func(*kev1alpha.ListClusterResponse))
// f3 is the function to retrieve the next page token from the response object. It takes the response message as input and returns the next page token. (e.g. func(*kev1alpha.ListClusterResponse) string)
//
// The page token is used to retrieve the next page of results. If it is empty, the first page will be retrieved.
// The page token is retrieved from the response object using the GetNextPageToken() method. If the token is empty, the last page has been reached and the function will return.
func ListAll[T1 any, T2 any](ctx context.Context, list func(context.Context, *connect.Request[T1]) (*connect.Response[T2], error), f1 func(string) *T1, f2 func(*T2), f3 func(*T2) string) error {
	token := ""
	for {
		v, err := list(ctx, connect.NewRequest(f1(token)))
		if err != nil {
			return err
		}
		msg := helper.GetResponseMsg(v)
		f2(msg)
		if token = f3(msg); token == "" {
			break
		}
	}
	return nil
}
