SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `advice`
    DROP FOREIGN KEY `fk_advice_creator`,
    DROP COLUMN `description`,
    DROP COLUMN `id_creator`,
    DROP COLUMN `created_at`;

SET FOREIGN_KEY_CHECKS = 1;
