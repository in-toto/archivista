// Code generated by ent, DO NOT EDIT.

package mapping

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the mapping type in the database.
	Label = "mapping"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldPath holds the string denoting the path field in the database.
	FieldPath = "path"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldSha1 holds the string denoting the sha1 field in the database.
	FieldSha1 = "sha1"
	// FieldSha256 holds the string denoting the sha256 field in the database.
	FieldSha256 = "sha256"
	// FieldGitoidSha1 holds the string denoting the gitoidsha1 field in the database.
	FieldGitoidSha1 = "gitoid_sha1"
	// FieldGitoidSha256 holds the string denoting the gitoidsha256 field in the database.
	FieldGitoidSha256 = "gitoid_sha256"
	// EdgePosix holds the string denoting the posix edge name in mutations.
	EdgePosix = "posix"
	// EdgeOmnitrail holds the string denoting the omnitrail edge name in mutations.
	EdgeOmnitrail = "omnitrail"
	// Table holds the table name of the mapping in the database.
	Table = "mappings"
	// PosixTable is the table that holds the posix relation/edge.
	PosixTable = "posixes"
	// PosixInverseTable is the table name for the Posix entity.
	// It exists in this package in order to avoid circular dependency with the "posix" package.
	PosixInverseTable = "posixes"
	// PosixColumn is the table column denoting the posix relation/edge.
	PosixColumn = "mapping_posix"
	// OmnitrailTable is the table that holds the omnitrail relation/edge.
	OmnitrailTable = "mappings"
	// OmnitrailInverseTable is the table name for the Omnitrail entity.
	// It exists in this package in order to avoid circular dependency with the "omnitrail" package.
	OmnitrailInverseTable = "omnitrails"
	// OmnitrailColumn is the table column denoting the omnitrail relation/edge.
	OmnitrailColumn = "omnitrail_mappings"
)

// Columns holds all SQL columns for mapping fields.
var Columns = []string{
	FieldID,
	FieldPath,
	FieldType,
	FieldSha1,
	FieldSha256,
	FieldGitoidSha1,
	FieldGitoidSha256,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "mappings"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"omnitrail_mappings",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// PathValidator is a validator for the "path" field. It is called by the builders before save.
	PathValidator func(string) error
	// TypeValidator is a validator for the "type" field. It is called by the builders before save.
	TypeValidator func(string) error
	// GitoidSha1Validator is a validator for the "gitoidSha1" field. It is called by the builders before save.
	GitoidSha1Validator func(string) error
	// GitoidSha256Validator is a validator for the "gitoidSha256" field. It is called by the builders before save.
	GitoidSha256Validator func(string) error
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the Mapping queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByPath orders the results by the path field.
func ByPath(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPath, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// BySha1 orders the results by the sha1 field.
func BySha1(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSha1, opts...).ToFunc()
}

// BySha256 orders the results by the sha256 field.
func BySha256(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSha256, opts...).ToFunc()
}

// ByGitoidSha1 orders the results by the gitoidSha1 field.
func ByGitoidSha1(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGitoidSha1, opts...).ToFunc()
}

// ByGitoidSha256 orders the results by the gitoidSha256 field.
func ByGitoidSha256(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGitoidSha256, opts...).ToFunc()
}

// ByPosixCount orders the results by posix count.
func ByPosixCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newPosixStep(), opts...)
	}
}

// ByPosix orders the results by posix terms.
func ByPosix(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPosixStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByOmnitrailField orders the results by omnitrail field.
func ByOmnitrailField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOmnitrailStep(), sql.OrderByField(field, opts...))
	}
}
func newPosixStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PosixInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, PosixTable, PosixColumn),
	)
}
func newOmnitrailStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OmnitrailInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, OmnitrailTable, OmnitrailColumn),
	)
}