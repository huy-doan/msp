-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `role_id` int NOT NULL,
  `enabled_mfa` tinyint(1) NOT NULL DEFAULT '1',
  `enabled_mfa_type_id` int DEFAULT NULL,
  `last_name` varchar(100) NOT NULL,
  `first_name` varchar(100) NOT NULL,
  `last_name_kana` varchar(100) NOT NULL,
  `first_name_kana` varchar(100) NOT NULL,
  `avatar_url` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `deleted_at` (`deleted_at`),
  KEY `fk_users_role` (`role_id`),
  KEY `fk_users_mfa_type` (`enabled_mfa_type_id`),
  CONSTRAINT `fk_users_mfa_type` FOREIGN KEY (`enabled_mfa_type_id`) REFERENCES `master_mfa_types` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE RESTRICT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
