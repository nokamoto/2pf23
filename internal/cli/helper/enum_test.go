package helper

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
)

func TestEnumFlag_Set(t *testing.T) {
	sut := NewEnumFlag[kev1alpha.MachineType](kev1alpha.MachineType_name, kev1alpha.MachineType_value)

	assert := func(want kev1alpha.MachineType) {
		if s := sut.String(); s != want.String() {
			t.Errorf("got %v, want %v", s, want.String())
		}
		if v := sut.Value(); v != want {
			t.Errorf("got %v, want %v", v, want)
		}
	}

	// default
	assert(kev1alpha.MachineType_MACHINE_TYPE_UNSPECIFIED)

	// set
	if err := sut.Set(kev1alpha.MachineType_MACHINE_TYPE_STANDARD.String()); err != nil {
		t.Errorf("got %v, want nil", err)
	}
	assert(kev1alpha.MachineType_MACHINE_TYPE_STANDARD)

	// unchanged
	if err := sut.Set("unknown"); err == nil {
		t.Errorf("got nil, want error")
	}
	assert(kev1alpha.MachineType_MACHINE_TYPE_STANDARD)
}

func TestEnumFlag_Type(t *testing.T) {
	sut := NewEnumFlag[kev1alpha.MachineType](kev1alpha.MachineType_name, kev1alpha.MachineType_value)

	if s := sut.Type(); s != "kev1alpha.MachineType" {
		t.Errorf("got %v, want %v", s, "kev1alpha.MachineType")
	}
}

func TestEnumFlag_CompletionFunc(t *testing.T) {
	sut := NewEnumFlag[kev1alpha.MachineType](kev1alpha.MachineType_name, kev1alpha.MachineType_value)

	got, _ := sut.CompletionFunc()(nil, nil, "")
	want := []string{
		kev1alpha.MachineType_MACHINE_TYPE_UNSPECIFIED.String(),
		kev1alpha.MachineType_MACHINE_TYPE_STANDARD.String(),
		kev1alpha.MachineType_MACHINE_TYPE_HIGHMEM.String(),
		kev1alpha.MachineType_MACHINE_TYPE_HIGHCPU.String(),
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}
