// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nokamoto/2pf23/internal/ent/cluster"
	"github.com/nokamoto/2pf23/internal/ent/predicate"
)

// ClusterUpdate is the builder for updating Cluster entities.
type ClusterUpdate struct {
	config
	hooks    []Hook
	mutation *ClusterMutation
}

// Where appends a list predicates to the ClusterUpdate builder.
func (cu *ClusterUpdate) Where(ps ...predicate.Cluster) *ClusterUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetName sets the "name" field.
func (cu *ClusterUpdate) SetName(s string) *ClusterUpdate {
	cu.mutation.SetName(s)
	return cu
}

// SetDisplayName sets the "display_name" field.
func (cu *ClusterUpdate) SetDisplayName(s string) *ClusterUpdate {
	cu.mutation.SetDisplayName(s)
	return cu
}

// SetNumNodes sets the "num_nodes" field.
func (cu *ClusterUpdate) SetNumNodes(i int32) *ClusterUpdate {
	cu.mutation.ResetNumNodes()
	cu.mutation.SetNumNodes(i)
	return cu
}

// AddNumNodes adds i to the "num_nodes" field.
func (cu *ClusterUpdate) AddNumNodes(i int32) *ClusterUpdate {
	cu.mutation.AddNumNodes(i)
	return cu
}

// SetMachineType sets the "machine_type" field.
func (cu *ClusterUpdate) SetMachineType(i int32) *ClusterUpdate {
	cu.mutation.ResetMachineType()
	cu.mutation.SetMachineType(i)
	return cu
}

// AddMachineType adds i to the "machine_type" field.
func (cu *ClusterUpdate) AddMachineType(i int32) *ClusterUpdate {
	cu.mutation.AddMachineType(i)
	return cu
}

// Mutation returns the ClusterMutation object of the builder.
func (cu *ClusterUpdate) Mutation() *ClusterMutation {
	return cu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *ClusterUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *ClusterUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *ClusterUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *ClusterUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cu *ClusterUpdate) check() error {
	if v, ok := cu.mutation.Name(); ok {
		if err := cluster.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Cluster.name": %w`, err)}
		}
	}
	if v, ok := cu.mutation.NumNodes(); ok {
		if err := cluster.NumNodesValidator(v); err != nil {
			return &ValidationError{Name: "num_nodes", err: fmt.Errorf(`ent: validator failed for field "Cluster.num_nodes": %w`, err)}
		}
	}
	if v, ok := cu.mutation.MachineType(); ok {
		if err := cluster.MachineTypeValidator(v); err != nil {
			return &ValidationError{Name: "machine_type", err: fmt.Errorf(`ent: validator failed for field "Cluster.machine_type": %w`, err)}
		}
	}
	return nil
}

func (cu *ClusterUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := cu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(cluster.Table, cluster.Columns, sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeInt64))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.Name(); ok {
		_spec.SetField(cluster.FieldName, field.TypeString, value)
	}
	if value, ok := cu.mutation.DisplayName(); ok {
		_spec.SetField(cluster.FieldDisplayName, field.TypeString, value)
	}
	if value, ok := cu.mutation.NumNodes(); ok {
		_spec.SetField(cluster.FieldNumNodes, field.TypeInt32, value)
	}
	if value, ok := cu.mutation.AddedNumNodes(); ok {
		_spec.AddField(cluster.FieldNumNodes, field.TypeInt32, value)
	}
	if value, ok := cu.mutation.MachineType(); ok {
		_spec.SetField(cluster.FieldMachineType, field.TypeInt32, value)
	}
	if value, ok := cu.mutation.AddedMachineType(); ok {
		_spec.AddField(cluster.FieldMachineType, field.TypeInt32, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cluster.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// ClusterUpdateOne is the builder for updating a single Cluster entity.
type ClusterUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ClusterMutation
}

// SetName sets the "name" field.
func (cuo *ClusterUpdateOne) SetName(s string) *ClusterUpdateOne {
	cuo.mutation.SetName(s)
	return cuo
}

// SetDisplayName sets the "display_name" field.
func (cuo *ClusterUpdateOne) SetDisplayName(s string) *ClusterUpdateOne {
	cuo.mutation.SetDisplayName(s)
	return cuo
}

// SetNumNodes sets the "num_nodes" field.
func (cuo *ClusterUpdateOne) SetNumNodes(i int32) *ClusterUpdateOne {
	cuo.mutation.ResetNumNodes()
	cuo.mutation.SetNumNodes(i)
	return cuo
}

// AddNumNodes adds i to the "num_nodes" field.
func (cuo *ClusterUpdateOne) AddNumNodes(i int32) *ClusterUpdateOne {
	cuo.mutation.AddNumNodes(i)
	return cuo
}

// SetMachineType sets the "machine_type" field.
func (cuo *ClusterUpdateOne) SetMachineType(i int32) *ClusterUpdateOne {
	cuo.mutation.ResetMachineType()
	cuo.mutation.SetMachineType(i)
	return cuo
}

// AddMachineType adds i to the "machine_type" field.
func (cuo *ClusterUpdateOne) AddMachineType(i int32) *ClusterUpdateOne {
	cuo.mutation.AddMachineType(i)
	return cuo
}

// Mutation returns the ClusterMutation object of the builder.
func (cuo *ClusterUpdateOne) Mutation() *ClusterMutation {
	return cuo.mutation
}

// Where appends a list predicates to the ClusterUpdate builder.
func (cuo *ClusterUpdateOne) Where(ps ...predicate.Cluster) *ClusterUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *ClusterUpdateOne) Select(field string, fields ...string) *ClusterUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Cluster entity.
func (cuo *ClusterUpdateOne) Save(ctx context.Context) (*Cluster, error) {
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *ClusterUpdateOne) SaveX(ctx context.Context) *Cluster {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *ClusterUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *ClusterUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cuo *ClusterUpdateOne) check() error {
	if v, ok := cuo.mutation.Name(); ok {
		if err := cluster.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Cluster.name": %w`, err)}
		}
	}
	if v, ok := cuo.mutation.NumNodes(); ok {
		if err := cluster.NumNodesValidator(v); err != nil {
			return &ValidationError{Name: "num_nodes", err: fmt.Errorf(`ent: validator failed for field "Cluster.num_nodes": %w`, err)}
		}
	}
	if v, ok := cuo.mutation.MachineType(); ok {
		if err := cluster.MachineTypeValidator(v); err != nil {
			return &ValidationError{Name: "machine_type", err: fmt.Errorf(`ent: validator failed for field "Cluster.machine_type": %w`, err)}
		}
	}
	return nil
}

func (cuo *ClusterUpdateOne) sqlSave(ctx context.Context) (_node *Cluster, err error) {
	if err := cuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(cluster.Table, cluster.Columns, sqlgraph.NewFieldSpec(cluster.FieldID, field.TypeInt64))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Cluster.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, cluster.FieldID)
		for _, f := range fields {
			if !cluster.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != cluster.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.Name(); ok {
		_spec.SetField(cluster.FieldName, field.TypeString, value)
	}
	if value, ok := cuo.mutation.DisplayName(); ok {
		_spec.SetField(cluster.FieldDisplayName, field.TypeString, value)
	}
	if value, ok := cuo.mutation.NumNodes(); ok {
		_spec.SetField(cluster.FieldNumNodes, field.TypeInt32, value)
	}
	if value, ok := cuo.mutation.AddedNumNodes(); ok {
		_spec.AddField(cluster.FieldNumNodes, field.TypeInt32, value)
	}
	if value, ok := cuo.mutation.MachineType(); ok {
		_spec.SetField(cluster.FieldMachineType, field.TypeInt32, value)
	}
	if value, ok := cuo.mutation.AddedMachineType(); ok {
		_spec.AddField(cluster.FieldMachineType, field.TypeInt32, value)
	}
	_node = &Cluster{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{cluster.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}
