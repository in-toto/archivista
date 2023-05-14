// Code generated by ent, DO NOT EDIT.

package runtime

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/project"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/schema"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/tenant"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/user"

	"entgo.io/ent"
	"entgo.io/ent/privacy"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	projectMixin := schema.Project{}.Mixin()
	project.Policy = privacy.NewPolicies(projectMixin[0], schema.Project{})
	project.Hooks[0] = func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if err := project.Policy.EvalMutation(ctx, m); err != nil {
				return nil, err
			}
			return next.Mutate(ctx, m)
		})
	}
	projectMixinFields1 := projectMixin[1].Fields()
	_ = projectMixinFields1
	projectFields := schema.Project{}.Fields()
	_ = projectFields
	// projectDescCreatedAt is the schema descriptor for created_at field.
	projectDescCreatedAt := projectMixinFields1[1].Descriptor()
	// project.DefaultCreatedAt holds the default value on creation for the created_at field.
	project.DefaultCreatedAt = projectDescCreatedAt.Default.(func() time.Time)
	// projectDescUpdatedAt is the schema descriptor for updated_at field.
	projectDescUpdatedAt := projectMixinFields1[2].Descriptor()
	// project.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	project.DefaultUpdatedAt = projectDescUpdatedAt.Default.(func() time.Time)
	// project.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	project.UpdateDefaultUpdatedAt = projectDescUpdatedAt.UpdateDefault.(func() time.Time)
	// projectDescRepoID is the schema descriptor for repo_id field.
	projectDescRepoID := projectFields[0].Descriptor()
	// project.RepoIDValidator is a validator for the "repo_id" field. It is called by the builders before save.
	project.RepoIDValidator = projectDescRepoID.Validators[0].(func(string) error)
	// projectDescName is the schema descriptor for name field.
	projectDescName := projectFields[1].Descriptor()
	// project.NameValidator is a validator for the "name" field. It is called by the builders before save.
	project.NameValidator = projectDescName.Validators[0].(func(string) error)
	// projectDescProjecturl is the schema descriptor for projecturl field.
	projectDescProjecturl := projectFields[2].Descriptor()
	// project.ProjecturlValidator is a validator for the "projecturl" field. It is called by the builders before save.
	project.ProjecturlValidator = projectDescProjecturl.Validators[0].(func(string) error)
	// projectDescID is the schema descriptor for id field.
	projectDescID := projectMixinFields1[0].Descriptor()
	// project.DefaultID holds the default value on creation for the id field.
	project.DefaultID = projectDescID.Default.(func() uuid.UUID)
	tenantMixin := schema.Tenant{}.Mixin()
	tenant.Policy = privacy.NewPolicies(schema.Tenant{})
	tenant.Hooks[0] = func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if err := tenant.Policy.EvalMutation(ctx, m); err != nil {
				return nil, err
			}
			return next.Mutate(ctx, m)
		})
	}
	tenantMixinFields0 := tenantMixin[0].Fields()
	_ = tenantMixinFields0
	tenantFields := schema.Tenant{}.Fields()
	_ = tenantFields
	// tenantDescCreatedAt is the schema descriptor for created_at field.
	tenantDescCreatedAt := tenantMixinFields0[1].Descriptor()
	// tenant.DefaultCreatedAt holds the default value on creation for the created_at field.
	tenant.DefaultCreatedAt = tenantDescCreatedAt.Default.(func() time.Time)
	// tenantDescUpdatedAt is the schema descriptor for updated_at field.
	tenantDescUpdatedAt := tenantMixinFields0[2].Descriptor()
	// tenant.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	tenant.DefaultUpdatedAt = tenantDescUpdatedAt.Default.(func() time.Time)
	// tenant.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	tenant.UpdateDefaultUpdatedAt = tenantDescUpdatedAt.UpdateDefault.(func() time.Time)
	// tenantDescID is the schema descriptor for id field.
	tenantDescID := tenantMixinFields0[0].Descriptor()
	// tenant.DefaultID holds the default value on creation for the id field.
	tenant.DefaultID = tenantDescID.Default.(func() uuid.UUID)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescIdentityID is the schema descriptor for identity_id field.
	userDescIdentityID := userFields[1].Descriptor()
	// user.IdentityIDValidator is a validator for the "identity_id" field. It is called by the builders before save.
	user.IdentityIDValidator = userDescIdentityID.Validators[0].(func(string) error)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}

const (
	Version = "v0.11.10"                                        // Version of ent codegen.
	Sum     = "h1:iqn32ybY5HRW3xSAyMNdNKpZhKgMf1Zunsej9yPKUI8=" // Sum of ent codegen.
)
