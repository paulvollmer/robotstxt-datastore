// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/paulvollmer/robotstxt-datastore/server/ent/predicate"
	"github.com/paulvollmer/robotstxt-datastore/server/ent/robotstxt"
)

// RobotstxtDelete is the builder for deleting a Robotstxt entity.
type RobotstxtDelete struct {
	config
	hooks      []Hook
	mutation   *RobotstxtMutation
	predicates []predicate.Robotstxt
}

// Where adds a new predicate to the delete builder.
func (rd *RobotstxtDelete) Where(ps ...predicate.Robotstxt) *RobotstxtDelete {
	rd.predicates = append(rd.predicates, ps...)
	return rd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (rd *RobotstxtDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(rd.hooks) == 0 {
		affected, err = rd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RobotstxtMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			rd.mutation = mutation
			affected, err = rd.sqlExec(ctx)
			return affected, err
		})
		for i := len(rd.hooks) - 1; i >= 0; i-- {
			mut = rd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (rd *RobotstxtDelete) ExecX(ctx context.Context) int {
	n, err := rd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (rd *RobotstxtDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: robotstxt.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: robotstxt.FieldID,
			},
		},
	}
	if ps := rd.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, rd.driver, _spec)
}

// RobotstxtDeleteOne is the builder for deleting a single Robotstxt entity.
type RobotstxtDeleteOne struct {
	rd *RobotstxtDelete
}

// Exec executes the deletion query.
func (rdo *RobotstxtDeleteOne) Exec(ctx context.Context) error {
	n, err := rdo.rd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{robotstxt.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (rdo *RobotstxtDeleteOne) ExecX(ctx context.Context) {
	rdo.rd.ExecX(ctx)
}
