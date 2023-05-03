package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Cluster holds the schema definition for the Cluster entity.
type Cluster struct {
	ent.Schema
}

// Fields of the Cluster.
func (Cluster) Fields() []ent.Field {
	return []ent.Field{
		field.Text("name").NotEmpty().Unique(),
		field.Text("display_name"),
		field.Int("num_nodes").Positive(),
	}
}

// Edges of the Cluster.
func (Cluster) Edges() []ent.Edge {
	return nil
}