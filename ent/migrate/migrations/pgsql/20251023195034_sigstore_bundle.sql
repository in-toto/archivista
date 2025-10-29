-- Modify "signatures" table
ALTER TABLE "signatures" ADD COLUMN "certificate" bytea NULL, ADD COLUMN "intermediates" jsonb NULL;
-- Modify "timestamps" table
ALTER TABLE "timestamps" ADD COLUMN "data" bytea NULL;
-- Create "sigstore_bundles" table
CREATE TABLE "sigstore_bundles" (
  "id" uuid NOT NULL,
  "gitoid_sha256" character varying NOT NULL,
  "media_type" character varying NOT NULL,
  "version" character varying NULL,
  "created_at" timestamptz NOT NULL,
  "dsse_bundle" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "sigstore_bundles_dsses_bundle" FOREIGN KEY ("dsse_bundle") REFERENCES "dsses" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "sigstore_bundles_dsse_bundle_key" to table: "sigstore_bundles"
CREATE UNIQUE INDEX "sigstore_bundles_dsse_bundle_key" ON "sigstore_bundles" ("dsse_bundle");
-- Create index "sigstore_bundles_gitoid_sha256_key" to table: "sigstore_bundles"
CREATE UNIQUE INDEX "sigstore_bundles_gitoid_sha256_key" ON "sigstore_bundles" ("gitoid_sha256");
-- Create index "sigstorebundle_gitoid_sha256" to table: "sigstore_bundles"
CREATE INDEX "sigstorebundle_gitoid_sha256" ON "sigstore_bundles" ("gitoid_sha256");
