-- Modify "dsses" table
ALTER TABLE "dsses" ADD COLUMN "created_at" timestamptz NULL;
-- Modify "subjects" table
ALTER TABLE "subjects" ADD COLUMN "created_at" timestamptz NULL;
