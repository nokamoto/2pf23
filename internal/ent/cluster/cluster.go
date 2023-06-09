// Code generated by ent, DO NOT EDIT.

package cluster

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the cluster type in the database.
	Label = "cluster"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDisplayName holds the string denoting the display_name field in the database.
	FieldDisplayName = "display_name"
	// FieldNumNodes holds the string denoting the num_nodes field in the database.
	FieldNumNodes = "num_nodes"
	// FieldMachineType holds the string denoting the machine_type field in the database.
	FieldMachineType = "machine_type"
	// Table holds the table name of the cluster in the database.
	Table = "clusters"
)

// Columns holds all SQL columns for cluster fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDisplayName,
	FieldNumNodes,
	FieldMachineType,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// NumNodesValidator is a validator for the "num_nodes" field. It is called by the builders before save.
	NumNodesValidator func(int32) error
	// MachineTypeValidator is a validator for the "machine_type" field. It is called by the builders before save.
	MachineTypeValidator func(int32) error
)

// OrderOption defines the ordering options for the Cluster queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDisplayName orders the results by the display_name field.
func ByDisplayName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDisplayName, opts...).ToFunc()
}

// ByNumNodes orders the results by the num_nodes field.
func ByNumNodes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNumNodes, opts...).ToFunc()
}

// ByMachineType orders the results by the machine_type field.
func ByMachineType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMachineType, opts...).ToFunc()
}
