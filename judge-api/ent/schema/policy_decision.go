package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Policy holds the schema definition for the Policy entity.
type PolicyDecision struct {
	ent.Schema
}

// Fields of the Policy.
func (PolicyDecision) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("subject_name").NotEmpty(),
		field.String("digest_id").NotEmpty(),
		field.Enum("decision").Values("allowed", "denied", "skipped").Default("denied"),
	}
}

// PolicyDecisionMixin is a mixin for adding common fields to the PolicyDecision.
type PolicyDecisionMixin struct {
	mixin.Schema
}

// Fields of the PolicyDecisionMixin.
func (PolicyDecisionMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("timestamp").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Policy.
func (PolicyDecision) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("policy_decisions"),
	}
}
