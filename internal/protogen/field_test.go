package protogen

import "testing"

func TestGoTypeNameFromFullyQualified(t *testing.T) {
	if got, want := GoTypeNameFromFullyQualified(".com.example.v1.CreateRequest"), "v1.CreateRequest"; got != want {
		t.Errorf("GoTypeNameFromFullyQualified() = %v, want %v", got, want)
	}
}
