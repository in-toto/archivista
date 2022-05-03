// Code generated by entc, DO NOT EDIT.

package signature

const (
	// Label holds the string label denoting the signature type in the database.
	Label = "signature"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldKeyID holds the string denoting the key_id field in the database.
	FieldKeyID = "key_id"
	// FieldSignature holds the string denoting the signature field in the database.
	FieldSignature = "signature"
	// EdgeDsse holds the string denoting the dsse edge name in mutations.
	EdgeDsse = "dsse"
	// Table holds the table name of the signature in the database.
	Table = "signatures"
	// DsseTable is the table that holds the dsse relation/edge.
	DsseTable = "signatures"
	// DsseInverseTable is the table name for the Dsse entity.
	// It exists in this package in order to avoid circular dependency with the "dsse" package.
	DsseInverseTable = "dsses"
	// DsseColumn is the table column denoting the dsse relation/edge.
	DsseColumn = "dsse_signatures"
)

// Columns holds all SQL columns for signature fields.
var Columns = []string{
	FieldID,
	FieldKeyID,
	FieldSignature,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "signatures"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"dsse_signatures",
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
	// KeyIDValidator is a validator for the "key_id" field. It is called by the builders before save.
	KeyIDValidator func(string) error
	// SignatureValidator is a validator for the "signature" field. It is called by the builders before save.
	SignatureValidator func(string) error
)
