INSERT INTO `permission` (`id_permission`, `name_permission`, `domain`, `created_at`) VALUES
  (UUID(), 'moderate_projects', 'projects', NOW());

INSERT IGNORE INTO `role_permission` (`id`, `role_id`, `permission_id`)
  SELECT UUID(), '11111111-1111-1111-1111-000000000004', id_permission FROM `permission` WHERE name_permission = 'moderate_projects' AND domain = 'projects';

ALTER TABLE `report`
  ADD COLUMN `id_project` CHAR(36) NULL DEFAULT NULL,
  ADD KEY `fk_report_project` (`id_project`),
  ADD CONSTRAINT `fk_report_project` FOREIGN KEY (`id_project`) REFERENCES `project` (`id_project`) ON DELETE SET NULL;
