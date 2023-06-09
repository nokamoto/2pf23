// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/nokamoto/2pf23/internal/ent/cluster"
)

// Cluster is the model entity for the Cluster schema.
type Cluster struct {
	config `json:"-"`
	// ID of the ent.
	ID int64 `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// DisplayName holds the value of the "display_name" field.
	DisplayName string `json:"display_name,omitempty"`
	// NumNodes holds the value of the "num_nodes" field.
	NumNodes int32 `json:"num_nodes,omitempty"`
	// MachineType holds the value of the "machine_type" field.
	MachineType  int32 `json:"machine_type,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Cluster) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case cluster.FieldID, cluster.FieldNumNodes, cluster.FieldMachineType:
			values[i] = new(sql.NullInt64)
		case cluster.FieldName, cluster.FieldDisplayName:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Cluster fields.
func (c *Cluster) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case cluster.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int64(value.Int64)
		case cluster.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case cluster.FieldDisplayName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field display_name", values[i])
			} else if value.Valid {
				c.DisplayName = value.String
			}
		case cluster.FieldNumNodes:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field num_nodes", values[i])
			} else if value.Valid {
				c.NumNodes = int32(value.Int64)
			}
		case cluster.FieldMachineType:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field machine_type", values[i])
			} else if value.Valid {
				c.MachineType = int32(value.Int64)
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Cluster.
// This includes values selected through modifiers, order, etc.
func (c *Cluster) Value(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// Update returns a builder for updating this Cluster.
// Note that you need to call Cluster.Unwrap() before calling this method if this Cluster
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Cluster) Update() *ClusterUpdateOne {
	return NewClusterClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Cluster entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Cluster) Unwrap() *Cluster {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Cluster is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Cluster) String() string {
	var builder strings.Builder
	builder.WriteString("Cluster(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("name=")
	builder.WriteString(c.Name)
	builder.WriteString(", ")
	builder.WriteString("display_name=")
	builder.WriteString(c.DisplayName)
	builder.WriteString(", ")
	builder.WriteString("num_nodes=")
	builder.WriteString(fmt.Sprintf("%v", c.NumNodes))
	builder.WriteString(", ")
	builder.WriteString("machine_type=")
	builder.WriteString(fmt.Sprintf("%v", c.MachineType))
	builder.WriteByte(')')
	return builder.String()
}

// Clusters is a parsable slice of Cluster.
type Clusters []*Cluster
