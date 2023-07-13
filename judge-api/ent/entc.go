//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithWhereInputs(true),
		entgql.WithRelaySpec(true),
		entgql.WithConfigPath("./gqlgen.yml"),
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("./ent.graphql"),
	)

	opts := []entc.Option{
		entc.FeatureNames("privacy"),
		entc.Extensions(ex),
		entc.Dependency(
			entc.DependencyName("Config"),
			entc.DependencyTypeInfo(&field.TypeInfo{
				Ident:   "configuration.Config",
				PkgPath: "github.com/testifysec/judge/judge-api/internal/configuration",
			}),
		),
	}

	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	if err := entc.Generate("./ent/schema", &gen.Config{}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
