ALTER TABLE `report`
  DROP FOREIGN KEY `fk_report_project`,
  DROP KEY `fk_report_project`,
  DROP COLUMN `id_project`;

DELETE rp FROM `role_permission` rp
  JOIN `permission` p ON p.id_permission = rp.permission_id
  WHERE p.name_permission = 'moderate_projects' AND p.domain = 'projects';

DELETE FROM `permission` WHERE name_permission = 'moderate_projects' AND domain = 'projects';
