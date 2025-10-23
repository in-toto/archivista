-- Modify "signatures" table
ALTER TABLE `signatures` ADD COLUMN `certificate` blob NULL AFTER `signature`, ADD COLUMN `intermediates` json NULL AFTER `certificate`;
-- Modify "timestamps" table
ALTER TABLE `timestamps` ADD COLUMN `data` blob NULL AFTER `timestamp`;
-- Create "sigstore_bundles" table
CREATE TABLE `sigstore_bundles` (
  `id` char(36) NOT NULL,
  `gitoid_sha256` varchar(255) NOT NULL,
  `media_type` varchar(255) NOT NULL,
  `version` varchar(255) NULL,
  `created_at` timestamp NOT NULL,
  `dsse_bundle` char(36) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `dsse_bundle` (`dsse_bundle`),
  UNIQUE INDEX `gitoid_sha256` (`gitoid_sha256`),
  INDEX `sigstorebundle_gitoid_sha256` (`gitoid_sha256`),
  CONSTRAINT `sigstore_bundles_dsses_bundle` FOREIGN KEY (`dsse_bundle`) REFERENCES `dsses` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
