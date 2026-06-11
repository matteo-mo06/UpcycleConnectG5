SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `advice`
    ADD COLUMN `description` text DEFAULT NULL,
    ADD COLUMN `id_creator` char(36) DEFAULT NULL,
    ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD CONSTRAINT `fk_advice_creator` FOREIGN KEY (`id_creator`) REFERENCES `user` (`id_user`);

SET FOREIGN_KEY_CHECKS = 1;
