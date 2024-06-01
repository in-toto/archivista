-- Create "statements" table
CREATE TABLE "statements" ("id" uuid NOT NULL, "predicate" character varying NOT NULL, PRIMARY KEY ("id"));
-- Create index "statement_predicate" to table: "statements"
CREATE INDEX "statement_predicate" ON "statements" ("predicate");
-- Create "attestation_collections" table
CREATE TABLE "attestation_collections" ("id" uuid NOT NULL, "name" character varying NOT NULL, "statement_attestation_collections" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "attestation_collections_statements_attestation_collections" FOREIGN KEY ("statement_attestation_collections") REFERENCES "statements" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "attestation_collections_statement_attestation_collections_key" to table: "attestation_collections"
CREATE UNIQUE INDEX "attestation_collections_statement_attestation_collections_key" ON "attestation_collections" ("statement_attestation_collections");
-- Create index "attestationcollection_name" to table: "attestation_collections"
CREATE INDEX "attestationcollection_name" ON "attestation_collections" ("name");
-- Create "attestation_policies" table
CREATE TABLE "attestation_policies" ("id" uuid NOT NULL, "name" character varying NOT NULL, "statement_policy" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "attestation_policies_statements_policy" FOREIGN KEY ("statement_policy") REFERENCES "statements" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "attestation_policies_statement_policy_key" to table: "attestation_policies"
CREATE UNIQUE INDEX "attestation_policies_statement_policy_key" ON "attestation_policies" ("statement_policy");
-- Create index "attestationpolicy_name" to table: "attestation_policies"
CREATE INDEX "attestationpolicy_name" ON "attestation_policies" ("name");
-- Create "attestations" table
CREATE TABLE "attestations" ("id" uuid NOT NULL, "type" character varying NOT NULL, "attestation_collection_attestations" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "attestations_attestation_collections_attestations" FOREIGN KEY ("attestation_collection_attestations") REFERENCES "attestation_collections" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "attestation_type" to table: "attestations"
CREATE INDEX "attestation_type" ON "attestations" ("type");
-- Create "dsses" table
CREATE TABLE "dsses" ("id" uuid NOT NULL, "gitoid_sha256" character varying NOT NULL, "payload_type" character varying NOT NULL, "dsse_statement" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "dsses_statements_statement" FOREIGN KEY ("dsse_statement") REFERENCES "statements" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "dsses_gitoid_sha256_key" to table: "dsses"
CREATE UNIQUE INDEX "dsses_gitoid_sha256_key" ON "dsses" ("gitoid_sha256");
-- Create "payload_digests" table
CREATE TABLE "payload_digests" ("id" uuid NOT NULL, "algorithm" character varying NOT NULL, "value" character varying NOT NULL, "dsse_payload_digests" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "payload_digests_dsses_payload_digests" FOREIGN KEY ("dsse_payload_digests") REFERENCES "dsses" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "payloaddigest_value" to table: "payload_digests"
CREATE INDEX "payloaddigest_value" ON "payload_digests" ("value");
-- Create "signatures" table
CREATE TABLE "signatures" ("id" uuid NOT NULL, "key_id" character varying NOT NULL, "signature" character varying NOT NULL, "dsse_signatures" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "signatures_dsses_signatures" FOREIGN KEY ("dsse_signatures") REFERENCES "dsses" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "signature_key_id" to table: "signatures"
CREATE INDEX "signature_key_id" ON "signatures" ("key_id");
-- Create "subjects" table
CREATE TABLE "subjects" ("id" uuid NOT NULL, "name" character varying NOT NULL, "statement_subjects" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "subjects_statements_subjects" FOREIGN KEY ("statement_subjects") REFERENCES "statements" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "subject_name" to table: "subjects"
CREATE INDEX "subject_name" ON "subjects" ("name");
-- Create "subject_digests" table
CREATE TABLE "subject_digests" ("id" uuid NOT NULL, "algorithm" character varying NOT NULL, "value" character varying NOT NULL, "subject_subject_digests" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "subject_digests_subjects_subject_digests" FOREIGN KEY ("subject_subject_digests") REFERENCES "subjects" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "subjectdigest_value" to table: "subject_digests"
CREATE INDEX "subjectdigest_value" ON "subject_digests" ("value");
-- Create "timestamps" table
CREATE TABLE "timestamps" ("id" uuid NOT NULL, "type" character varying NOT NULL, "timestamp" timestamptz NOT NULL, "signature_timestamps" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "timestamps_signatures_timestamps" FOREIGN KEY ("signature_timestamps") REFERENCES "signatures" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
