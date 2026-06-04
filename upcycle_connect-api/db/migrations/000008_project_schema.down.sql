DELETE rp FROM `role_permission` rp
  JOIN `permission` p ON p.id_permission = rp.permission_id
  WHERE p.name_permission = 'create_project' AND p.domain = 'projects'
    AND rp.role_id = '11111111-1111-1111-1111-000000000003';

DELETE rp FROM `role_permission` rp
  JOIN `permission` p ON p.id_permission = rp.permission_id
  WHERE p.name_permission = 'register_project' AND p.domain = 'projects';

DELETE FROM `permission` WHERE name_permission = 'register_project' AND domain = 'projects';

ALTER TABLE `project`
  DROP COLUMN `title_project`,
  DROP COLUMN `description_project`,
  DROP COLUMN `location_project`,
  DROP COLUMN `status`,
  DROP COLUMN `rejection_reason`,
  DROP COLUMN `created_at`;
