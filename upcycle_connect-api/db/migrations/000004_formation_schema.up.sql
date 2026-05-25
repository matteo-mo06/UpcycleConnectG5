SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `formation`
    ADD COLUMN `title_formation` varchar(255) NOT NULL DEFAULT '',
    ADD COLUMN `description_formation` text,
    ADD COLUMN `date_formation` datetime DEFAULT NULL,
    ADD COLUMN `location_formation` varchar(255) DEFAULT NULL,
    ADD COLUMN `capacity` int DEFAULT NULL,
    ADD COLUMN `level` varchar(20) NOT NULL DEFAULT 'beginner',
    ADD COLUMN `duration_hours` int DEFAULT NULL,
    ADD COLUMN `status` varchar(20) NOT NULL DEFAULT 'pending',
    ADD COLUMN `rejection_reason` text DEFAULT NULL,
    ADD COLUMN `id_category` char(36) DEFAULT NULL,
    ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP;

INSERT IGNORE INTO `role_permission` (`id`, `role_id`, `permission_id`) VALUES
    ('f1e2d3c4-b5a6-4f7e-8d9c-0b1a2f3e4d5c', '11111111-1111-1111-1111-000000000004', '9a0b1c2d-3e4f-4a5b-9c6d-7e8f9a0b1c2d');

SET FOREIGN_KEY_CHECKS = 1;
