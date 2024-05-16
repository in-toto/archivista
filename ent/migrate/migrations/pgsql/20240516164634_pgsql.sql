-- Modify "vex_statements" table
ALTER TABLE "vex_statements" ADD COLUMN "vuln_id" character varying NOT NULL;
-- Create index "vexstatement_vuln_id" to table: "vex_statements"
CREATE INDEX "vexstatement_vuln_id" ON "vex_statements" ("vuln_id");
