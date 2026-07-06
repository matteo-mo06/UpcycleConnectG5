SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `formation_step` (
  `id_step`      char(36) NOT NULL,
  `id_formation` char(36) NOT NULL,
  `title`        varchar(255) NOT NULL,
  `description`  text DEFAULT NULL,
  `status`       varchar(20) NOT NULL DEFAULT 'todo',
  `step_order`   int NOT NULL DEFAULT 0,
  PRIMARY KEY (`id_step`),
  KEY `fk_step_formation` (`id_formation`),
  CONSTRAINT `fk_step_formation` FOREIGN KEY (`id_formation`) REFERENCES `formation` (`id_formation`)
);

SET FOREIGN_KEY_CHECKS = 1;
