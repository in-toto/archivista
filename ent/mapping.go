// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/in-toto/archivista/ent/mapping"
	"github.com/in-toto/archivista/ent/omnitrail"
)

// Mapping is the model entity for the Mapping schema.
type Mapping struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Path holds the value of the "path" field.
	Path string `json:"path,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Sha1 holds the value of the "sha1" field.
	Sha1 string `json:"sha1,omitempty"`
	// Sha256 holds the value of the "sha256" field.
	Sha256 string `json:"sha256,omitempty"`
	// GitoidSha1 holds the value of the "gitoidSha1" field.
	GitoidSha1 string `json:"gitoidSha1,omitempty"`
	// GitoidSha256 holds the value of the "gitoidSha256" field.
	GitoidSha256 string `json:"gitoidSha256,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MappingQuery when eager-loading is set.
	Edges              MappingEdges `json:"edges"`
	omnitrail_mappings *uuid.UUID
	selectValues       sql.SelectValues
}

// MappingEdges holds the relations/edges for other nodes in the graph.
type MappingEdges struct {
	// Posix holds the value of the posix edge.
	Posix []*Posix `json:"posix,omitempty"`
	// Omnitrail holds the value of the omnitrail edge.
	Omnitrail *Omnitrail `json:"omnitrail,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
	// totalCount holds the count of the edges above.
	totalCount [2]map[string]int

	namedPosix map[string][]*Posix
}

// PosixOrErr returns the Posix value or an error if the edge
// was not loaded in eager-loading.
func (e MappingEdges) PosixOrErr() ([]*Posix, error) {
	if e.loadedTypes[0] {
		return e.Posix, nil
	}
	return nil, &NotLoadedError{edge: "posix"}
}

// OmnitrailOrErr returns the Omnitrail value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MappingEdges) OmnitrailOrErr() (*Omnitrail, error) {
	if e.Omnitrail != nil {
		return e.Omnitrail, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: omnitrail.Label}
	}
	return nil, &NotLoadedError{edge: "omnitrail"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Mapping) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case mapping.FieldPath, mapping.FieldType, mapping.FieldSha1, mapping.FieldSha256, mapping.FieldGitoidSha1, mapping.FieldGitoidSha256:
			values[i] = new(sql.NullString)
		case mapping.FieldID:
			values[i] = new(uuid.UUID)
		case mapping.ForeignKeys[0]: // omnitrail_mappings
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Mapping fields.
func (m *Mapping) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case mapping.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				m.ID = *value
			}
		case mapping.FieldPath:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field path", values[i])
			} else if value.Valid {
				m.Path = value.String
			}
		case mapping.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				m.Type = value.String
			}
		case mapping.FieldSha1:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sha1", values[i])
			} else if value.Valid {
				m.Sha1 = value.String
			}
		case mapping.FieldSha256:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sha256", values[i])
			} else if value.Valid {
				m.Sha256 = value.String
			}
		case mapping.FieldGitoidSha1:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field gitoidSha1", values[i])
			} else if value.Valid {
				m.GitoidSha1 = value.String
			}
		case mapping.FieldGitoidSha256:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field gitoidSha256", values[i])
			} else if value.Valid {
				m.GitoidSha256 = value.String
			}
		case mapping.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field omnitrail_mappings", values[i])
			} else if value.Valid {
				m.omnitrail_mappings = new(uuid.UUID)
				*m.omnitrail_mappings = *value.S.(*uuid.UUID)
			}
		default:
			m.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Mapping.
// This includes values selected through modifiers, order, etc.
func (m *Mapping) Value(name string) (ent.Value, error) {
	return m.selectValues.Get(name)
}

// QueryPosix queries the "posix" edge of the Mapping entity.
func (m *Mapping) QueryPosix() *PosixQuery {
	return NewMappingClient(m.config).QueryPosix(m)
}

// QueryOmnitrail queries the "omnitrail" edge of the Mapping entity.
func (m *Mapping) QueryOmnitrail() *OmnitrailQuery {
	return NewMappingClient(m.config).QueryOmnitrail(m)
}

// Update returns a builder for updating this Mapping.
// Note that you need to call Mapping.Unwrap() before calling this method if this Mapping
// was returned from a transaction, and the transaction was committed or rolled back.
func (m *Mapping) Update() *MappingUpdateOne {
	return NewMappingClient(m.config).UpdateOne(m)
}

// Unwrap unwraps the Mapping entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (m *Mapping) Unwrap() *Mapping {
	_tx, ok := m.config.driver.(*txDriver)
	if !ok {
		panic("ent: Mapping is not a transactional entity")
	}
	m.config.driver = _tx.drv
	return m
}

// String implements the fmt.Stringer.
func (m *Mapping) String() string {
	var builder strings.Builder
	builder.WriteString("Mapping(")
	builder.WriteString(fmt.Sprintf("id=%v, ", m.ID))
	builder.WriteString("path=")
	builder.WriteString(m.Path)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(m.Type)
	builder.WriteString(", ")
	builder.WriteString("sha1=")
	builder.WriteString(m.Sha1)
	builder.WriteString(", ")
	builder.WriteString("sha256=")
	builder.WriteString(m.Sha256)
	builder.WriteString(", ")
	builder.WriteString("gitoidSha1=")
	builder.WriteString(m.GitoidSha1)
	builder.WriteString(", ")
	builder.WriteString("gitoidSha256=")
	builder.WriteString(m.GitoidSha256)
	builder.WriteByte(')')
	return builder.String()
}

// NamedPosix returns the Posix named value or an error if the edge was not
// loaded in eager-loading with this name.
func (m *Mapping) NamedPosix(name string) ([]*Posix, error) {
	if m.Edges.namedPosix == nil {
		return nil, &NotLoadedError{edge: name}
	}
	nodes, ok := m.Edges.namedPosix[name]
	if !ok {
		return nil, &NotLoadedError{edge: name}
	}
	return nodes, nil
}

func (m *Mapping) appendNamedPosix(name string, edges ...*Posix) {
	if m.Edges.namedPosix == nil {
		m.Edges.namedPosix = make(map[string][]*Posix)
	}
	if len(edges) == 0 {
		m.Edges.namedPosix[name] = []*Posix{}
	} else {
		m.Edges.namedPosix[name] = append(m.Edges.namedPosix[name], edges...)
	}
}

// Mappings is a parsable slice of Mapping.
type Mappings []*Mapping