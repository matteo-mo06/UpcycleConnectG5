SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `advice`
    DROP FOREIGN KEY `fk_advice_category`,
    DROP COLUMN `id_category`,
    ADD COLUMN `tag` varchar(120) DEFAULT NULL;

SET FOREIGN_KEY_CHECKS = 1;
