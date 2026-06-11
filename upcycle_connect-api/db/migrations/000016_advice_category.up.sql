SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `advice`
    DROP COLUMN `tag`,
    ADD COLUMN `id_category` char(36) DEFAULT NULL,
    ADD CONSTRAINT `fk_advice_category` FOREIGN KEY (`id_category`) REFERENCES `category` (`id_category`);

SET FOREIGN_KEY_CHECKS = 1;
