
CREATE TABLE IF NOT EXISTS `file` (
  `id` VARCHAR(128) NOT NULL,
  `name` VARCHAR(4096) NOT NULL,
  `path` TEXT NOT NULL,
  `mimetype` VARCHAR(256) NOT NULL,
  `extension` VARCHAR(128) NOT NULL,
  `size` BIGINT NOT NULL,
  `created_at` BIGINT NOT NULL,
  `updated_at` BIGINT NOT NULL,
  `deleted_at` BIGINT NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) 
DEFAULT CHARACTER SET utf8
COLLATE utf8_unicode_ci
ENGINE = InnoDB;
