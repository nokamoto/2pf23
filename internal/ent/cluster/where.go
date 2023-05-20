// Code generated by ent, DO NOT EDIT.

package cluster

import (
	"entgo.io/ent/dialect/sql"
	"github.com/nokamoto/2pf23/internal/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int64) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int64) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int64) predicate.Cluster {
	return predicate.Cluster(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int64) predicate.Cluster {
	return predicate.Cluster(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int64) predicate.Cluster {
	return predicate.Cluster(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int64) predicate.Cluster {
	return predicate.Cluster(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int64) predicate.Cluster {
	return predicate.Cluster(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int64) predicate.Cluster {
	return predicate.Cluster(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int64) predicate.Cluster {
	return predicate.Cluster(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldName, v))
}

// DisplayName applies equality check predicate on the "display_name" field. It's identical to DisplayNameEQ.
func DisplayName(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldDisplayName, v))
}

// NumNodes applies equality check predicate on the "num_nodes" field. It's identical to NumNodesEQ.
func NumNodes(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldNumNodes, v))
}

// MachineType applies equality check predicate on the "machine_type" field. It's identical to MachineTypeEQ.
func MachineType(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldMachineType, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Cluster {
	return predicate.Cluster(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Cluster {
	return predicate.Cluster(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldContainsFold(FieldName, v))
}

// DisplayNameEQ applies the EQ predicate on the "display_name" field.
func DisplayNameEQ(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldDisplayName, v))
}

// DisplayNameNEQ applies the NEQ predicate on the "display_name" field.
func DisplayNameNEQ(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldNEQ(FieldDisplayName, v))
}

// DisplayNameIn applies the In predicate on the "display_name" field.
func DisplayNameIn(vs ...string) predicate.Cluster {
	return predicate.Cluster(sql.FieldIn(FieldDisplayName, vs...))
}

// DisplayNameNotIn applies the NotIn predicate on the "display_name" field.
func DisplayNameNotIn(vs ...string) predicate.Cluster {
	return predicate.Cluster(sql.FieldNotIn(FieldDisplayName, vs...))
}

// DisplayNameGT applies the GT predicate on the "display_name" field.
func DisplayNameGT(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldGT(FieldDisplayName, v))
}

// DisplayNameGTE applies the GTE predicate on the "display_name" field.
func DisplayNameGTE(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldGTE(FieldDisplayName, v))
}

// DisplayNameLT applies the LT predicate on the "display_name" field.
func DisplayNameLT(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldLT(FieldDisplayName, v))
}

// DisplayNameLTE applies the LTE predicate on the "display_name" field.
func DisplayNameLTE(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldLTE(FieldDisplayName, v))
}

// DisplayNameContains applies the Contains predicate on the "display_name" field.
func DisplayNameContains(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldContains(FieldDisplayName, v))
}

// DisplayNameHasPrefix applies the HasPrefix predicate on the "display_name" field.
func DisplayNameHasPrefix(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldHasPrefix(FieldDisplayName, v))
}

// DisplayNameHasSuffix applies the HasSuffix predicate on the "display_name" field.
func DisplayNameHasSuffix(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldHasSuffix(FieldDisplayName, v))
}

// DisplayNameEqualFold applies the EqualFold predicate on the "display_name" field.
func DisplayNameEqualFold(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldEqualFold(FieldDisplayName, v))
}

// DisplayNameContainsFold applies the ContainsFold predicate on the "display_name" field.
func DisplayNameContainsFold(v string) predicate.Cluster {
	return predicate.Cluster(sql.FieldContainsFold(FieldDisplayName, v))
}

// NumNodesEQ applies the EQ predicate on the "num_nodes" field.
func NumNodesEQ(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldNumNodes, v))
}

// NumNodesNEQ applies the NEQ predicate on the "num_nodes" field.
func NumNodesNEQ(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldNEQ(FieldNumNodes, v))
}

// NumNodesIn applies the In predicate on the "num_nodes" field.
func NumNodesIn(vs ...int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldIn(FieldNumNodes, vs...))
}

// NumNodesNotIn applies the NotIn predicate on the "num_nodes" field.
func NumNodesNotIn(vs ...int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldNotIn(FieldNumNodes, vs...))
}

// NumNodesGT applies the GT predicate on the "num_nodes" field.
func NumNodesGT(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldGT(FieldNumNodes, v))
}

// NumNodesGTE applies the GTE predicate on the "num_nodes" field.
func NumNodesGTE(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldGTE(FieldNumNodes, v))
}

// NumNodesLT applies the LT predicate on the "num_nodes" field.
func NumNodesLT(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldLT(FieldNumNodes, v))
}

// NumNodesLTE applies the LTE predicate on the "num_nodes" field.
func NumNodesLTE(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldLTE(FieldNumNodes, v))
}

// MachineTypeEQ applies the EQ predicate on the "machine_type" field.
func MachineTypeEQ(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldEQ(FieldMachineType, v))
}

// MachineTypeNEQ applies the NEQ predicate on the "machine_type" field.
func MachineTypeNEQ(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldNEQ(FieldMachineType, v))
}

// MachineTypeIn applies the In predicate on the "machine_type" field.
func MachineTypeIn(vs ...int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldIn(FieldMachineType, vs...))
}

// MachineTypeNotIn applies the NotIn predicate on the "machine_type" field.
func MachineTypeNotIn(vs ...int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldNotIn(FieldMachineType, vs...))
}

// MachineTypeGT applies the GT predicate on the "machine_type" field.
func MachineTypeGT(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldGT(FieldMachineType, v))
}

// MachineTypeGTE applies the GTE predicate on the "machine_type" field.
func MachineTypeGTE(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldGTE(FieldMachineType, v))
}

// MachineTypeLT applies the LT predicate on the "machine_type" field.
func MachineTypeLT(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldLT(FieldMachineType, v))
}

// MachineTypeLTE applies the LTE predicate on the "machine_type" field.
func MachineTypeLTE(v int32) predicate.Cluster {
	return predicate.Cluster(sql.FieldLTE(FieldMachineType, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Cluster) predicate.Cluster {
	return predicate.Cluster(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Cluster) predicate.Cluster {
	return predicate.Cluster(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Cluster) predicate.Cluster {
	return predicate.Cluster(func(s *sql.Selector) {
		p(s.Not())
	})
}
