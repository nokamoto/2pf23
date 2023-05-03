// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/nokamoto/2pf23/internal/ent/cluster"
	"github.com/nokamoto/2pf23/internal/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	clusterFields := schema.Cluster{}.Fields()
	_ = clusterFields
	// clusterDescName is the schema descriptor for name field.
	clusterDescName := clusterFields[0].Descriptor()
	// cluster.NameValidator is a validator for the "name" field. It is called by the builders before save.
	cluster.NameValidator = clusterDescName.Validators[0].(func(string) error)
	// clusterDescNumNodes is the schema descriptor for num_nodes field.
	clusterDescNumNodes := clusterFields[2].Descriptor()
	// cluster.NumNodesValidator is a validator for the "num_nodes" field. It is called by the builders before save.
	cluster.NumNodesValidator = clusterDescNumNodes.Validators[0].(func(int) error)
}