SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `project_material` (
  `id_material` char(36) NOT NULL,
  `id_project`  char(36) NOT NULL,
  `name`        varchar(255) NOT NULL,
  `quantity`    decimal(10,2) NOT NULL DEFAULT 0,
  `unit`        varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id_material`),
  KEY `fk_material_project` (`id_project`),
  CONSTRAINT `fk_material_project` FOREIGN KEY (`id_project`) REFERENCES `project` (`id_project`)
);

CREATE TABLE `project_step` (
  `id_step`     char(36) NOT NULL,
  `id_project`  char(36) NOT NULL,
  `title`       varchar(255) NOT NULL,
  `description` text DEFAULT NULL,
  `status`      varchar(20) NOT NULL DEFAULT 'todo',
  `step_order`  int NOT NULL DEFAULT 0,
  PRIMARY KEY (`id_step`),
  KEY `fk_step_project` (`id_project`),
  CONSTRAINT `fk_step_project` FOREIGN KEY (`id_project`) REFERENCES `project` (`id_project`)
);

SET FOREIGN_KEY_CHECKS = 1;
