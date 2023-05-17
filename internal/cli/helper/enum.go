package helper

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

// EnumFlag implements pflag.Value for proto enum T.
type EnumFlag[T ~int32] struct {
	value  T
	values map[string]int32
	names  map[int32]string
}

func NewEnumFlag[T ~int32](names map[int32]string, values map[string]int32) *EnumFlag[T] {
	return &EnumFlag[T]{
		values: values,
		names:  names,
	}
}

func (e *EnumFlag[T]) String() string {
	return e.names[int32(e.value)]
}

func (e *EnumFlag[T]) Set(s string) error {
	v, ok := e.values[s]
	if !ok {
		return fmt.Errorf("unknown value %v", s)
	}
	e.value = T(v)
	return nil
}

func (e *EnumFlag[T]) Type() string {
	return fmt.Sprintf("%T", e.value)
}

func (e *EnumFlag[T]) Value() T {
	return e.value
}

// Names returns all the names of the enum in order.
func (e *EnumFlag[T]) Names() []string {
	var i []int32
	for _, v := range e.values {
		i = append(i, v)
	}
	slices.Sort(i)
	var s []string
	for _, v := range i {
		s = append(s, e.names[v])
	}
	return s
}

// CompletionFunc returns a completion function for the flag.
// https://github.com/spf13/cobra/blob/main/shell_completions.md#completions-for-flags
func (e *EnumFlag[T]) CompletionFunc() func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return e.Names(), cobra.ShellCompDirectiveDefault
	}
}
