package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Target represents a TUF target metadata.
type Target struct {
	ent.Schema
}

func (Target) Fields() []ent.Field {
	return []ent.Field{
		field.String("gitoid_sha256").NotEmpty().Unique(),
		field.Int("length"),
		field.String("version"),
		field.Enum("target_type").Values("policy"),
	}
}

func (Target) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("signatures", Signature.Type),
		edge.To("payload_digests", PayloadDigest.Type),
		edge.To("dsse", Dsse.Type),
		edge.To("subjects", Subject.Type),
	}
}

// Key represents a key in delegations.
type Key struct {
	ent.Schema
}

// public key data
func (Key) Fields() []ent.Field {
	return []ent.Field{
		//actual key data in the tuf document
		field.String("keyid").NotEmpty(),
		field.String("key").NotEmpty(),
		field.Enum("keytype").Values("ed25519", "ecdsa", "rsa"),
		field.Bool("intermediate"),
		field.Bool("certificate_authority"),
		field.String("scheme").NotEmpty(),

		// Certificate-specific fields
		field.String("issuer").NotEmpty(),              // Issuer DN
		field.String("subject").NotEmpty(),             // Subject DN
		field.Time("not_before"),                       // Certificate validity not before
		field.Time("not_after"),                        // Certificate validity not after
		field.String("serial_number").NotEmpty(),       // Serial number
		field.String("signature_algorithm").NotEmpty(), // Signature algorithm
		field.String("common_name").NotEmpty(),         // Common name
		field.String("organization").NotEmpty(),        // Organization
		field.String("organizational_unit").NotEmpty(), // Organizational unit
		field.Strings("dns_names"),
		field.Strings("email_addresses"),
		field.Strings("uris"),

		// ux
		field.String("name"),
		field.String("description"),
		field.String("comment"),
		field.String("usage"),
	}
}

func (Key) Edges() []ent.Edge {
	return []ent.Edge{
		//digest of the pem encoded key
		edge.To("payload_digests", PayloadDigest.Type),
	}
}

// CertConstraint represents certificate constraints.
type CertConstraint struct {
	ent.Schema
}

func (CertConstraint) Fields() []ent.Field {
	return []ent.Field{
		field.String("commonname").NotEmpty(),
		field.Strings("dnsnames"),
		field.Strings("emails"),
		field.Strings("organizations"),
		field.Strings("uris"),
	}
}

// RootMetadata represents the root TUF metadata.
type RootMetadata struct {
	ent.Schema
}

func (RootMetadata) Fields() []ent.Field {
	return []ent.Field{
		field.String("gitoid_sha256").NotEmpty().Unique(),
		field.String("_type").Default("targets"),
		field.String("spec_version"),
		field.Int("version"),
		field.Time("expires"),
	}
}

func (RootMetadata) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("payload_digests", PayloadDigest.Type),
		edge.To("targets", Target.Type),
		edge.To("keys", Key.Type),
		edge.To("certConstraints", CertConstraint.Type),
		edge.To("signatures", Signature.Type),
		edge.To("dsse", Dsse.Type),
	}
}
