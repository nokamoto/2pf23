package helper

import "github.com/bufbuild/connect-go"

// CodeOf returns connect.CodeOf(err) if err is not nil. Otherwise, returns 0 for OK.
// connect.CodeOf(err) returns codes.Unknown if err is nil, but it is not useful for testing.
func CodeOf(err error) connect.Code {
	if err == nil {
		return 0
	}
	return connect.CodeOf(err)
}

// GetMsg returns v.Msg if res is not nil. Otherwise, returns nil.
//
// Original grpc-go code implements nil check in generated code to access message field and able to get default value,
// but connect-go does not because [connect.Request].Msg panics if it is nil.
//
// ```
// var foo *T = nil
// foo.GetCluster() // ok
//
// var bar connect.Request[T] = nil
// bar.Msg.GetCluster() // panic
// ```
//
// To avoid this panic, this function is introduced.
//
// ```
// var bar connect.Request[T] = nil
// GetMsg(bar).GetCluster() // ok
// ```
func GetMsg[T any](v *connect.Request[T]) *T {
	if v == nil {
		return nil
	}
	return v.Msg
}

// GetResponseMsg returns v.Msg if res is not nil. Otherwise, returns nil.
// See GetMsg for details.
func GetResponseMsg[T any](v *connect.Response[T]) *T {
	if v == nil {
		return nil
	}
	return v.Msg
}
