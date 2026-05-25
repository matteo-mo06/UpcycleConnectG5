SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `formation`
    DROP COLUMN `title_formation`,
    DROP COLUMN `description_formation`,
    DROP COLUMN `date_formation`,
    DROP COLUMN `location_formation`,
    DROP COLUMN `capacity`,
    DROP COLUMN `level`,
    DROP COLUMN `duration_hours`,
    DROP COLUMN `status`,
    DROP COLUMN `rejection_reason`,
    DROP COLUMN `id_category`,
    DROP COLUMN `created_at`;

DELETE FROM `role_permission` WHERE `id` = 'f1e2d3c4-b5a6-4f7e-8d9c-0b1a2f3e4d5c';

SET FOREIGN_KEY_CHECKS = 1;
