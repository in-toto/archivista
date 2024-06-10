-- Create "sarifs" table
CREATE TABLE "sarifs" ("id" uuid NOT NULL, "report_file_name" character varying NOT NULL, "statement_sarif" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "sarifs_statements_sarif" FOREIGN KEY ("statement_sarif") REFERENCES "statements" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "sarif_report_file_name" to table: "sarifs"
CREATE INDEX "sarif_report_file_name" ON "sarifs" ("report_file_name");
-- Create index "sarifs_statement_sarif_key" to table: "sarifs"
CREATE UNIQUE INDEX "sarifs_statement_sarif_key" ON "sarifs" ("statement_sarif");
-- Create "sarif_rules" table
CREATE TABLE "sarif_rules" ("id" uuid NOT NULL, "rule_id" character varying NOT NULL, "rule_name" character varying NOT NULL, "short_description" character varying NOT NULL, "sarif_sarif_rules" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "sarif_rules_sarifs_sarif_rules" FOREIGN KEY ("sarif_sarif_rules") REFERENCES "sarifs" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "sarif_rules_sarif_sarif_rules_key" to table: "sarif_rules"
CREATE UNIQUE INDEX "sarif_rules_sarif_sarif_rules_key" ON "sarif_rules" ("sarif_sarif_rules");
-- Create index "sarifrule_rule_id" to table: "sarif_rules"
CREATE INDEX "sarifrule_rule_id" ON "sarif_rules" ("rule_id");
-- Create index "sarifrule_rule_name" to table: "sarif_rules"
CREATE INDEX "sarifrule_rule_name" ON "sarif_rules" ("rule_name");
-- Create index "sarifrule_short_description" to table: "sarif_rules"
CREATE INDEX "sarifrule_short_description" ON "sarif_rules" ("short_description");
