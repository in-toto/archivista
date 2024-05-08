-- Create "metadata" table
CREATE TABLE `metadata` (`id` bigint NOT NULL AUTO_INCREMENT, `key` varchar(255) NOT NULL, `value` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE INDEX `metadata_key_value` (`key`, `value`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "dsse_metadata" table
CREATE TABLE `dsse_metadata` (`dsse_id` bigint NOT NULL, `metadata_id` bigint NOT NULL, PRIMARY KEY (`dsse_id`, `metadata_id`), INDEX `dsse_metadata_metadata_id` (`metadata_id`), CONSTRAINT `dsse_metadata_dsse_id` FOREIGN KEY (`dsse_id`) REFERENCES `dsses` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT `dsse_metadata_metadata_id` FOREIGN KEY (`metadata_id`) REFERENCES `metadata` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE) CHARSET utf8mb4 COLLATE utf8mb4_bin;
