-- Modify "vex_statements" table
ALTER TABLE `vex_statements` ADD COLUMN `vuln_id` varchar(255) NOT NULL, ADD INDEX `vexstatement_vuln_id` (`vuln_id`);
