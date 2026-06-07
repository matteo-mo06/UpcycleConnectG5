ALTER TABLE `project`
  ADD COLUMN `title_project`       VARCHAR(255) NOT NULL DEFAULT '',
  ADD COLUMN `description_project` TEXT         DEFAULT NULL,
  ADD COLUMN `location_project`    VARCHAR(255) DEFAULT NULL,
  ADD COLUMN `status`              VARCHAR(60)  NOT NULL DEFAULT 'pending',
  ADD COLUMN `rejection_reason`    TEXT         DEFAULT NULL,
  ADD COLUMN `created_at`          DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP;

INSERT INTO `permission` (`id_permission`, `name_permission`, `domain`, `created_at`) VALUES
  (UUID(), 'register_project', 'projects', NOW());

INSERT IGNORE INTO `role_permission` (`id`, `role_id`, `permission_id`)
  SELECT UUID(), '11111111-1111-1111-1111-000000000003', id_permission FROM `permission` WHERE name_permission = 'create_project' AND domain = 'projects';
INSERT IGNORE INTO `role_permission` (`id`, `role_id`, `permission_id`)
  SELECT UUID(), '11111111-1111-1111-1111-000000000003', id_permission FROM `permission` WHERE name_permission = 'register_project' AND domain = 'projects';
INSERT IGNORE INTO `role_permission` (`id`, `role_id`, `permission_id`)
  SELECT UUID(), '11111111-1111-1111-1111-000000000002', id_permission FROM `permission` WHERE name_permission = 'register_project' AND domain = 'projects';
