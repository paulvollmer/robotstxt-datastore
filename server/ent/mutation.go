// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"time"

	"github.com/paulvollmer/robotstxt-datastore/server/ent/robotstxt"

	"github.com/facebookincubator/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeRobotstxt = "Robotstxt"
)

// RobotstxtMutation represents an operation that mutate the Robotstxts
// nodes in the graph.
type RobotstxtMutation struct {
	config
	op               Op
	typ              string
	id               *int
	created_at       *time.Time
	updated_at       *time.Time
	host             *string
	scheme           *string
	response_url     *string
	body             *[]byte
	statuscode       *int32
	addstatuscode    *int32
	response_time    *int64
	addresponse_time *int64
	clearedFields    map[string]struct{}
}

var _ ent.Mutation = (*RobotstxtMutation)(nil)

// newRobotstxtMutation creates new mutation for $n.Name.
func newRobotstxtMutation(c config, op Op) *RobotstxtMutation {
	return &RobotstxtMutation{
		config:        c,
		op:            op,
		typ:           TypeRobotstxt,
		clearedFields: make(map[string]struct{}),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m RobotstxtMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m RobotstxtMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *RobotstxtMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreatedAt sets the created_at field.
func (m *RobotstxtMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the created_at value in the mutation.
func (m *RobotstxtMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreatedAt reset all changes of the created_at field.
func (m *RobotstxtMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the updated_at field.
func (m *RobotstxtMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the updated_at value in the mutation.
func (m *RobotstxtMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdatedAt reset all changes of the updated_at field.
func (m *RobotstxtMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetHost sets the host field.
func (m *RobotstxtMutation) SetHost(s string) {
	m.host = &s
}

// Host returns the host value in the mutation.
func (m *RobotstxtMutation) Host() (r string, exists bool) {
	v := m.host
	if v == nil {
		return
	}
	return *v, true
}

// ResetHost reset all changes of the host field.
func (m *RobotstxtMutation) ResetHost() {
	m.host = nil
}

// SetScheme sets the scheme field.
func (m *RobotstxtMutation) SetScheme(s string) {
	m.scheme = &s
}

// Scheme returns the scheme value in the mutation.
func (m *RobotstxtMutation) Scheme() (r string, exists bool) {
	v := m.scheme
	if v == nil {
		return
	}
	return *v, true
}

// ResetScheme reset all changes of the scheme field.
func (m *RobotstxtMutation) ResetScheme() {
	m.scheme = nil
}

// SetResponseURL sets the response_url field.
func (m *RobotstxtMutation) SetResponseURL(s string) {
	m.response_url = &s
}

// ResponseURL returns the response_url value in the mutation.
func (m *RobotstxtMutation) ResponseURL() (r string, exists bool) {
	v := m.response_url
	if v == nil {
		return
	}
	return *v, true
}

// ResetResponseURL reset all changes of the response_url field.
func (m *RobotstxtMutation) ResetResponseURL() {
	m.response_url = nil
}

// SetBody sets the body field.
func (m *RobotstxtMutation) SetBody(b []byte) {
	m.body = &b
}

// Body returns the body value in the mutation.
func (m *RobotstxtMutation) Body() (r []byte, exists bool) {
	v := m.body
	if v == nil {
		return
	}
	return *v, true
}

// ResetBody reset all changes of the body field.
func (m *RobotstxtMutation) ResetBody() {
	m.body = nil
}

// SetStatuscode sets the statuscode field.
func (m *RobotstxtMutation) SetStatuscode(i int32) {
	m.statuscode = &i
	m.addstatuscode = nil
}

// Statuscode returns the statuscode value in the mutation.
func (m *RobotstxtMutation) Statuscode() (r int32, exists bool) {
	v := m.statuscode
	if v == nil {
		return
	}
	return *v, true
}

// AddStatuscode adds i to statuscode.
func (m *RobotstxtMutation) AddStatuscode(i int32) {
	if m.addstatuscode != nil {
		*m.addstatuscode += i
	} else {
		m.addstatuscode = &i
	}
}

// AddedStatuscode returns the value that was added to the statuscode field in this mutation.
func (m *RobotstxtMutation) AddedStatuscode() (r int32, exists bool) {
	v := m.addstatuscode
	if v == nil {
		return
	}
	return *v, true
}

// ResetStatuscode reset all changes of the statuscode field.
func (m *RobotstxtMutation) ResetStatuscode() {
	m.statuscode = nil
	m.addstatuscode = nil
}

// SetResponseTime sets the response_time field.
func (m *RobotstxtMutation) SetResponseTime(i int64) {
	m.response_time = &i
	m.addresponse_time = nil
}

// ResponseTime returns the response_time value in the mutation.
func (m *RobotstxtMutation) ResponseTime() (r int64, exists bool) {
	v := m.response_time
	if v == nil {
		return
	}
	return *v, true
}

// AddResponseTime adds i to response_time.
func (m *RobotstxtMutation) AddResponseTime(i int64) {
	if m.addresponse_time != nil {
		*m.addresponse_time += i
	} else {
		m.addresponse_time = &i
	}
}

// AddedResponseTime returns the value that was added to the response_time field in this mutation.
func (m *RobotstxtMutation) AddedResponseTime() (r int64, exists bool) {
	v := m.addresponse_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetResponseTime reset all changes of the response_time field.
func (m *RobotstxtMutation) ResetResponseTime() {
	m.response_time = nil
	m.addresponse_time = nil
}

// Op returns the operation name.
func (m *RobotstxtMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Robotstxt).
func (m *RobotstxtMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *RobotstxtMutation) Fields() []string {
	fields := make([]string, 0, 8)
	if m.created_at != nil {
		fields = append(fields, robotstxt.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, robotstxt.FieldUpdatedAt)
	}
	if m.host != nil {
		fields = append(fields, robotstxt.FieldHost)
	}
	if m.scheme != nil {
		fields = append(fields, robotstxt.FieldScheme)
	}
	if m.response_url != nil {
		fields = append(fields, robotstxt.FieldResponseURL)
	}
	if m.body != nil {
		fields = append(fields, robotstxt.FieldBody)
	}
	if m.statuscode != nil {
		fields = append(fields, robotstxt.FieldStatuscode)
	}
	if m.response_time != nil {
		fields = append(fields, robotstxt.FieldResponseTime)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *RobotstxtMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case robotstxt.FieldCreatedAt:
		return m.CreatedAt()
	case robotstxt.FieldUpdatedAt:
		return m.UpdatedAt()
	case robotstxt.FieldHost:
		return m.Host()
	case robotstxt.FieldScheme:
		return m.Scheme()
	case robotstxt.FieldResponseURL:
		return m.ResponseURL()
	case robotstxt.FieldBody:
		return m.Body()
	case robotstxt.FieldStatuscode:
		return m.Statuscode()
	case robotstxt.FieldResponseTime:
		return m.ResponseTime()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *RobotstxtMutation) SetField(name string, value ent.Value) error {
	switch name {
	case robotstxt.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case robotstxt.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case robotstxt.FieldHost:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetHost(v)
		return nil
	case robotstxt.FieldScheme:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetScheme(v)
		return nil
	case robotstxt.FieldResponseURL:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetResponseURL(v)
		return nil
	case robotstxt.FieldBody:
		v, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBody(v)
		return nil
	case robotstxt.FieldStatuscode:
		v, ok := value.(int32)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatuscode(v)
		return nil
	case robotstxt.FieldResponseTime:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetResponseTime(v)
		return nil
	}
	return fmt.Errorf("unknown Robotstxt field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *RobotstxtMutation) AddedFields() []string {
	var fields []string
	if m.addstatuscode != nil {
		fields = append(fields, robotstxt.FieldStatuscode)
	}
	if m.addresponse_time != nil {
		fields = append(fields, robotstxt.FieldResponseTime)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *RobotstxtMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case robotstxt.FieldStatuscode:
		return m.AddedStatuscode()
	case robotstxt.FieldResponseTime:
		return m.AddedResponseTime()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *RobotstxtMutation) AddField(name string, value ent.Value) error {
	switch name {
	case robotstxt.FieldStatuscode:
		v, ok := value.(int32)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddStatuscode(v)
		return nil
	case robotstxt.FieldResponseTime:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddResponseTime(v)
		return nil
	}
	return fmt.Errorf("unknown Robotstxt numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *RobotstxtMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *RobotstxtMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *RobotstxtMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Robotstxt nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *RobotstxtMutation) ResetField(name string) error {
	switch name {
	case robotstxt.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case robotstxt.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case robotstxt.FieldHost:
		m.ResetHost()
		return nil
	case robotstxt.FieldScheme:
		m.ResetScheme()
		return nil
	case robotstxt.FieldResponseURL:
		m.ResetResponseURL()
		return nil
	case robotstxt.FieldBody:
		m.ResetBody()
		return nil
	case robotstxt.FieldStatuscode:
		m.ResetStatuscode()
		return nil
	case robotstxt.FieldResponseTime:
		m.ResetResponseTime()
		return nil
	}
	return fmt.Errorf("unknown Robotstxt field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *RobotstxtMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *RobotstxtMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *RobotstxtMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *RobotstxtMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *RobotstxtMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *RobotstxtMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *RobotstxtMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Robotstxt unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *RobotstxtMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Robotstxt edge %s", name)
}
