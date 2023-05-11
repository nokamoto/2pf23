// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ClustersColumns holds the columns for the "clusters" table.
	ClustersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true, Size: 2147483647},
		{Name: "display_name", Type: field.TypeString, Size: 2147483647},
		{Name: "num_nodes", Type: field.TypeInt32},
	}
	// ClustersTable holds the schema information for the "clusters" table.
	ClustersTable = &schema.Table{
		Name:       "clusters",
		Columns:    ClustersColumns,
		PrimaryKey: []*schema.Column{ClustersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ClustersTable,
	}
)

func init() {
}
