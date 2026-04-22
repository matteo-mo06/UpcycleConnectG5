SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `user` (
  `id_user` char(36) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_user` varchar(255) NOT NULL,
  `first_name` varchar(150) DEFAULT NULL,
  `last_name` varchar(150) DEFAULT NULL,
  `upcycling_score` int NOT NULL DEFAULT '0',
  `premium` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` varchar(20) NOT NULL DEFAULT 'active',
  PRIMARY KEY (`id_user`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `advice` (
  `id_advice` char(36) NOT NULL,
  `title` varchar(255) NOT NULL,
  `tag` varchar(120) DEFAULT NULL,
  `date_advice` date DEFAULT NULL,
  PRIMARY KEY (`id_advice`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `category` (
  `id_category` char(36) NOT NULL,
  `name_category` varchar(100) NOT NULL,
  `description_category` text,
  PRIMARY KEY (`id_category`),
  UNIQUE KEY `name_category` (`name_category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `notification` (
  `id_notification` char(36) NOT NULL,
  `view` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_notification`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `payement` (
  `id_payement` char(36) NOT NULL,
  `amount_cents` int NOT NULL,
  `currency` char(3) NOT NULL,
  `provider_payement` varchar(60) DEFAULT NULL,
  `provider_ref` varchar(255) DEFAULT NULL,
  `status_payement` varchar(60) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `paid_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_payement`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `subscription` (
  `id_subscription` char(36) NOT NULL,
  `start_timestamp` datetime DEFAULT NULL,
  `end_timestamp` datetime DEFAULT NULL,
  `cancelled` tinyint(1) NOT NULL DEFAULT '0',
  `cancelled_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_subscription`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `subscription_plans` (
  `id_plan` char(36) NOT NULL,
  `price_cents` int NOT NULL,
  `interval_unit` varchar(30) NOT NULL,
  `interval_count` int NOT NULL DEFAULT '1',
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id_plan`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `ticket` (
  `id_ticket` char(36) NOT NULL,
  `timestamp_ticket` datetime DEFAULT NULL,
  PRIMARY KEY (`id_ticket`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `topic` (
  `id_topic` char(36) NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `post` (
  `id_post` char(36) NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_post`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `role` (
  `id_role` char(36) NOT NULL,
  `name_role` varchar(100) NOT NULL,
  PRIMARY KEY (`id_role`),
  UNIQUE KEY `name_role` (`name_role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `permission` (
  `id_permission` char(36) NOT NULL,
  `name_permission` varchar(100) NOT NULL,
  `domain` varchar(50) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_permission`),
  UNIQUE KEY `name_permission` (`name_permission`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `locker` (
  `id_locker` char(36) NOT NULL,
  `access_code` varchar(120) DEFAULT NULL,
  `locker_number` int DEFAULT NULL,
  PRIMARY KEY (`id_locker`),
  UNIQUE KEY `locker_number_unique` (`locker_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `project` (
  `id_project` char(36) NOT NULL,
  `start_date_project` date DEFAULT NULL,
  `end_date` date DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `id_creator` char(36) DEFAULT NULL,
  PRIMARY KEY (`id_project`),
  KEY `fk_project_creator` (`id_creator`),
  CONSTRAINT `fk_project_creator` FOREIGN KEY (`id_creator`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `formation` (
  `id_formation` char(36) NOT NULL,
  `id_creator` char(36) DEFAULT NULL,
  `id_formateur` char(36) DEFAULT NULL,
  PRIMARY KEY (`id_formation`),
  KEY `fk_formation_creator` (`id_creator`),
  KEY `fk_formation_formateur` (`id_formateur`),
  CONSTRAINT `fk_formation_creator` FOREIGN KEY (`id_creator`) REFERENCES `user` (`id_user`),
  CONSTRAINT `fk_formation_formateur` FOREIGN KEY (`id_formateur`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `event` (
  `id_event` char(36) NOT NULL,
  `title_event` varchar(255) NOT NULL DEFAULT '',
  `description_event` text,
  `date_event` datetime DEFAULT NULL,
  `location_event` varchar(255) DEFAULT NULL,
  `capacity` int DEFAULT NULL,
  `price_cents` int NOT NULL DEFAULT '0',
  `id_creator` char(36) DEFAULT NULL,
  PRIMARY KEY (`id_event`),
  KEY `fk_event_creator` (`id_creator`),
  CONSTRAINT `fk_event_creator` FOREIGN KEY (`id_creator`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `announcement` (
  `id_announcement` char(36) NOT NULL,
  `address_annoucement` varchar(255) DEFAULT NULL,
  `city` varchar(120) DEFAULT NULL,
  `postal` varchar(30) DEFAULT NULL,
  `description_annoucement` text,
  `availability_date` date NOT NULL,
  `price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `request` tinyint(1) NOT NULL DEFAULT '0',
  `state_annoucement` varchar(60) DEFAULT NULL,
  `id_category` char(36) DEFAULT NULL,
  `title_announcement` varchar(255) NOT NULL,
  `type_announcement` varchar(20) NOT NULL DEFAULT 'don',
  `condition_announcement` varchar(20) DEFAULT NULL,
  `id_buyer` char(36) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_announcement`),
  KEY `fk_announcement_category` (`id_category`),
  CONSTRAINT `fk_announcement_category` FOREIGN KEY (`id_category`) REFERENCES `category` (`id_category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `professional_request` (
  `id_request` char(36) NOT NULL,
  `id_user` char(36) NOT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `processed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_request`),
  KEY `fk_professional_request_user` (`id_user`),
  CONSTRAINT `fk_professional_request_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `document` (
  `id_document` char(36) NOT NULL,
  `id_user` char(36) NOT NULL,
  `category` varchar(120) DEFAULT NULL,
  `link` varchar(2048) DEFAULT NULL,
  PRIMARY KEY (`id_document`),
  KEY `fk_document_user` (`id_user`),
  CONSTRAINT `fk_document_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `report` (
  `id_report` char(36) NOT NULL,
  `id_user` char(36) DEFAULT NULL,
  `id_announcement` char(36) DEFAULT NULL,
  `id_topic` char(36) DEFAULT NULL,
  `id_post` char(36) DEFAULT NULL,
  `reason` text,
  `status` varchar(20) NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `id_reported_user` char(36) DEFAULT NULL,
  `resolved_by` char(36) DEFAULT NULL,
  `action_taken` varchar(50) DEFAULT NULL,
  `resolved_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id_report`),
  KEY `fk_report_user` (`id_user`),
  KEY `fk_report_announcement` (`id_announcement`),
  KEY `fk_report_topic` (`id_topic`),
  KEY `fk_report_post` (`id_post`),
  KEY `fk_report_reported_user` (`id_reported_user`),
  KEY `fk_report_resolved_by` (`resolved_by`),
  CONSTRAINT `fk_report_announcement` FOREIGN KEY (`id_announcement`) REFERENCES `announcement` (`id_announcement`),
  CONSTRAINT `fk_report_post` FOREIGN KEY (`id_post`) REFERENCES `post` (`id_post`),
  CONSTRAINT `fk_report_reported_user` FOREIGN KEY (`id_reported_user`) REFERENCES `user` (`id_user`),
  CONSTRAINT `fk_report_resolved_by` FOREIGN KEY (`resolved_by`) REFERENCES `user` (`id_user`),
  CONSTRAINT `fk_report_topic` FOREIGN KEY (`id_topic`) REFERENCES `topic` (`id_topic`),
  CONSTRAINT `fk_report_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_sanctions` (
  `id_sanction` char(36) NOT NULL,
  `id_user` char(36) NOT NULL,
  `id_admin` char(36) NOT NULL,
  `id_report` char(36) DEFAULT NULL,
  `type` varchar(20) NOT NULL,
  `reason` text,
  `duration_days` int DEFAULT NULL,
  `expires_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_sanction`),
  KEY `fk_sanctions_user` (`id_user`),
  KEY `fk_sanctions_admin` (`id_admin`),
  KEY `fk_sanctions_report` (`id_report`),
  CONSTRAINT `fk_sanctions_admin` FOREIGN KEY (`id_admin`) REFERENCES `user` (`id_user`),
  CONSTRAINT `fk_sanctions_report` FOREIGN KEY (`id_report`) REFERENCES `report` (`id_report`),
  CONSTRAINT `fk_sanctions_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `announcement_locker` (
  `id_announcement` char(36) NOT NULL,
  `id_locker` char(36) NOT NULL,
  PRIMARY KEY (`id_announcement`),
  UNIQUE KEY `id_locker` (`id_locker`),
  CONSTRAINT `fk_announcement_locker_announcement` FOREIGN KEY (`id_announcement`) REFERENCES `announcement` (`id_announcement`),
  CONSTRAINT `fk_announcement_locker_locker` FOREIGN KEY (`id_locker`) REFERENCES `locker` (`id_locker`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `role_permission` (
  `id` char(36) NOT NULL,
  `role_id` char(36) NOT NULL,
  `permission_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_role_permission` (`role_id`,`permission_id`),
  KEY `fk_rp_permission` (`permission_id`),
  CONSTRAINT `fk_rp_permission` FOREIGN KEY (`permission_id`) REFERENCES `permission` (`id_permission`),
  CONSTRAINT `fk_rp_role` FOREIGN KEY (`role_id`) REFERENCES `role` (`id_role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `topic_post` (
  `id_topic` char(36) NOT NULL,
  `id_post` char(36) NOT NULL,
  PRIMARY KEY (`id_topic`,`id_post`),
  KEY `fk_topic_post_post` (`id_post`),
  CONSTRAINT `fk_topic_post_post` FOREIGN KEY (`id_post`) REFERENCES `post` (`id_post`),
  CONSTRAINT `fk_topic_post_topic` FOREIGN KEY (`id_topic`) REFERENCES `topic` (`id_topic`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `sub_sub_plan` (
  `id_subscription` char(36) NOT NULL,
  `id_plan` char(36) NOT NULL,
  PRIMARY KEY (`id_subscription`,`id_plan`),
  KEY `fk_sub_sub_plan_plan` (`id_plan`),
  CONSTRAINT `fk_sub_sub_plan_plan` FOREIGN KEY (`id_plan`) REFERENCES `subscription_plans` (`id_plan`),
  CONSTRAINT `fk_sub_sub_plan_subscription` FOREIGN KEY (`id_subscription`) REFERENCES `subscription` (`id_subscription`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `payement_event` (
  `id_payement` char(36) NOT NULL,
  `id_event` char(36) NOT NULL,
  PRIMARY KEY (`id_payement`,`id_event`),
  KEY `fk_payement_event_event` (`id_event`),
  CONSTRAINT `fk_payement_event_event` FOREIGN KEY (`id_event`) REFERENCES `event` (`id_event`),
  CONSTRAINT `fk_payement_event_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `payement_formation` (
  `id_payement` char(36) NOT NULL,
  `id_formation` char(36) NOT NULL,
  PRIMARY KEY (`id_payement`,`id_formation`),
  KEY `fk_payement_formation_formation` (`id_formation`),
  CONSTRAINT `fk_payement_formation_formation` FOREIGN KEY (`id_formation`) REFERENCES `formation` (`id_formation`),
  CONSTRAINT `fk_payement_formation_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `payement_project` (
  `id_payement` char(36) NOT NULL,
  `id_project` char(36) NOT NULL,
  PRIMARY KEY (`id_payement`,`id_project`),
  KEY `fk_payement_project_project` (`id_project`),
  CONSTRAINT `fk_payement_project_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`),
  CONSTRAINT `fk_payement_project_project` FOREIGN KEY (`id_project`) REFERENCES `project` (`id_project`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_advice` (
  `id_user` char(36) NOT NULL,
  `id_advice` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_advice`),
  KEY `fk_user_advice_advice` (`id_advice`),
  CONSTRAINT `fk_user_advice_advice` FOREIGN KEY (`id_advice`) REFERENCES `advice` (`id_advice`),
  CONSTRAINT `fk_user_advice_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_announcement` (
  `id_user` char(36) NOT NULL,
  `id_announcement` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_announcement`),
  KEY `fk_user_announcement_announcement` (`id_announcement`),
  CONSTRAINT `fk_user_announcement_announcement` FOREIGN KEY (`id_announcement`) REFERENCES `announcement` (`id_announcement`),
  CONSTRAINT `fk_user_announcement_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_event_inscription` (
  `id_user` char(36) NOT NULL,
  `id_event` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_event`),
  KEY `fk_user_event_inscription_event` (`id_event`),
  CONSTRAINT `fk_user_event_inscription_event` FOREIGN KEY (`id_event`) REFERENCES `event` (`id_event`),
  CONSTRAINT `fk_user_event_inscription_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_formation_inscription` (
  `id_user` char(36) NOT NULL,
  `id_formation` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_formation`),
  KEY `fk_user_formation_inscription_formation` (`id_formation`),
  CONSTRAINT `fk_user_formation_inscription_formation` FOREIGN KEY (`id_formation`) REFERENCES `formation` (`id_formation`),
  CONSTRAINT `fk_user_formation_inscription_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_formation_inscription_payement` (
  `id_user` char(36) NOT NULL,
  `id_formation` char(36) NOT NULL,
  `id_payement` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_formation`,`id_payement`),
  KEY `fk_ufip_formation` (`id_formation`),
  KEY `fk_ufip_payement` (`id_payement`),
  CONSTRAINT `fk_ufip_formation` FOREIGN KEY (`id_formation`) REFERENCES `formation` (`id_formation`),
  CONSTRAINT `fk_ufip_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`),
  CONSTRAINT `fk_ufip_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_notification` (
  `id_user` char(36) NOT NULL,
  `id_notification` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_notification`),
  KEY `fk_user_notification_notification` (`id_notification`),
  CONSTRAINT `fk_user_notification_notification` FOREIGN KEY (`id_notification`) REFERENCES `notification` (`id_notification`),
  CONSTRAINT `fk_user_notification_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_payement` (
  `id_user` char(36) NOT NULL,
  `id_payement` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_payement`),
  KEY `fk_user_payement_payement` (`id_payement`),
  CONSTRAINT `fk_user_payement_payement` FOREIGN KEY (`id_payement`) REFERENCES `payement` (`id_payement`),
  CONSTRAINT `fk_user_payement_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_project_inscription` (
  `id_user` char(36) NOT NULL,
  `id_project` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_project`),
  KEY `fk_user_project_inscription_project` (`id_project`),
  CONSTRAINT `fk_user_project_inscription_project` FOREIGN KEY (`id_project`) REFERENCES `project` (`id_project`),
  CONSTRAINT `fk_user_project_inscription_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_project_updown` (
  `id_user` char(36) NOT NULL,
  `id_project` char(36) NOT NULL,
  `updown` smallint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_user`,`id_project`),
  KEY `fk_user_project_updown_project` (`id_project`),
  CONSTRAINT `fk_user_project_updown_project` FOREIGN KEY (`id_project`) REFERENCES `project` (`id_project`),
  CONSTRAINT `fk_user_project_updown_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_role` (
  `id_user` char(36) NOT NULL,
  `id_role` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_role`),
  KEY `fk_user_role_role` (`id_role`),
  CONSTRAINT `fk_user_role_role` FOREIGN KEY (`id_role`) REFERENCES `role` (`id_role`),
  CONSTRAINT `fk_user_role_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `user_subscription` (
  `id_user` char(36) NOT NULL,
  `id_subscription` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_subscription`),
  KEY `fk_user_subscription_subscription` (`id_subscription`),
  CONSTRAINT `fk_user_subscription_subscription` FOREIGN KEY (`id_subscription`) REFERENCES `subscription` (`id_subscription`),
  CONSTRAINT `fk_user_subscription_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Données de seed
INSERT INTO `user` VALUES ('7cd909a9-a27a-4849-a867-55e39b871c66','admin@upcycle.com','$2a$10$JHayhaAb0AYcARtBUWkyuuxf2fN64DydnGgbHSVPQfA60uvY5mo5K','Admin','Upcycle',0,0,'2026-04-22 21:47:08','active');

INSERT INTO `role` VALUES ('77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','admin'),('11111111-1111-1111-1111-000000000002','artisan'),('11111111-1111-1111-1111-000000000004','salarie'),('11111111-1111-1111-1111-000000000003','user');

INSERT INTO `user_role` VALUES ('7cd909a9-a27a-4849-a867-55e39b871c66','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab');

INSERT INTO `category` VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-000000000001','Couture & Textile','Retouches, créations et recyclage textile'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000002','Menuiserie & Bois','Réparation et création de meubles en bois'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000003','Électronique','Réparation d\'appareils électroniques'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000004','Maroquinerie','Réparation et création en cuir'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000005','Céramique & Poterie','Création et restauration céramique'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000006','Jardinage & Compostage','Compost, récupération déchets verts, jardinage écologique'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000007','Peinture & Décoration','Remise en peinture et relooking de mobilier usagé'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000008','Récupération Métaux','Travail des métaux et matériaux de récupération industrielle'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000009','Papeterie & Carton','Recyclage papier, carton et créations associées'),('aaaaaaaa-aaaa-aaaa-aaaa-000000000010','Mode & Couture Upcyclée','Transformation de vêtements et accessoires usagés');

INSERT INTO `permission` VALUES ('0b1c2d3e-4f5a-4b6c-9d7e-8f9a0b1c2d3e','moderate_forum','forum','2026-04-22 21:47:08'),('0e1f2a3b-4c5d-4e6f-8a7b-8c9d0e1f2a3b','manage_announcements','announcements','2026-04-22 21:47:08'),('1a2b3c4d-5e6f-4a7b-9c8d-0e1f2a3b3c4d','manage_payments','administration','2026-04-22 21:47:08'),('1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f','register_formation','formations','2026-04-22 21:47:08'),('1d2e3f4a-5b6c-4d7e-bf8a-9b0c1d2e3f4a','create_event','events','2026-04-22 21:47:08'),('2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d','request_locker_deposit','lockers','2026-04-22 21:47:08'),('2c3d4e5f-6a7b-4c8d-ae0f-1a2b3c4d5e6f','manage_categories','administration','2026-04-22 21:47:08'),('2d3e4f5a-6b7c-4d8e-9f0a-0b1c2d3e4f5a','manage_advice','conseils','2026-04-22 21:47:08'),('2e3f4a5b-6c7d-4e8f-8a9b-0c1d2e3f4a5b','manage_projects','projects','2026-04-22 21:47:08'),('3f2d8a1b-4c5e-4f6a-9b0c-1d2e3f4a5b6c','manage_users','administration','2026-04-22 21:47:08'),('3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c','register_event','events','2026-04-22 21:47:08'),('4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d','create_post','forum','2026-04-22 21:47:08'),('4c5d6e7f-8a9b-4c0d-8e1f-2a3b4c5d6e7f','manage_lockers','lockers','2026-04-22 21:47:08'),('4f5a6b7c-8d9e-4f0a-8b1c-2d3e4f5a6b7c','send_notifications','administration','2026-04-22 21:47:08'),('5b6c7d8e-9f0a-4b1c-8d2e-3f4a5b6c7d8e','manage_formations','formations','2026-04-22 21:47:08'),('5d6e7f8a-9b0c-4d1e-bf2a-3b4c5d6e7f8a','create_announcement','announcements','2026-04-22 21:47:08'),('6b7c8d9e-0f1a-4b2c-9d3e-4f5a6b7c8d9e','manage_subscriptions','administration','2026-04-22 21:47:08'),('6c7d8e9f-0a1b-4c2d-8e3f-4a5b6c7d8e9f','create_advice','conseils','2026-04-22 21:47:08'),('6d7e8f9a-0b1c-4d2e-9f3a-4b5c6d7e8f9a','create_project','projects','2026-04-22 21:47:08'),('73875cfa-b134-4b6d-9a46-43ea4fe9b500','validate_deposit','lockers','2026-04-22 21:47:08'),('7a8b9c0d-1e2f-4a3b-8c4d-5e6f7a8b9c0d','manage_roles','administration','2026-04-22 21:47:08'),('7e8f9a0b-1c2d-4e3f-9a4b-5c6d7e8f9a0b','manage_events','events','2026-04-22 21:47:08'),('7f8a9b0c-1d2e-4f3a-9b4c-5d6e7f8a9b0c','buy_announcement','announcements','2026-04-22 21:47:08'),('8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e','pickup_locker','lockers','2026-04-22 21:47:08'),('8e9f0a1b-2c3d-4e5f-9a6b-7c8d9e0f1a2b','manage_professional_requests','administration','2026-04-22 21:47:08'),('8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c','create_topic','forum','2026-04-22 21:47:08'),('9a0b1c2d-3e4f-4a5b-9c6d-7e8f9a0b1c2d','create_formation','formations','2026-04-22 21:47:08'),('9c0d1e2f-3a4b-4c5d-8e6f-7a8b9c0d1e2f','manage_documents','administration','2026-04-22 21:47:08');

INSERT INTO `role_permission` VALUES ('0bd46713-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','0e1f2a3b-4c5d-4e6f-8a7b-8c9d0e1f2a3b'),('0bd46ad6-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f'),('0bd4634c-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','1d2e3f4a-5b6c-4d7e-bf8a-9b0c1d2e3f4a'),('0bd46b6e-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d'),('0bd4667a-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','2d3e4f5a-6b7c-4d8e-9f0a-0b1c2d3e4f5a'),('0bd46908-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','2e3f4a5b-6c7d-4e8f-8a9b-0c1d2e3f4a5b'),('0bd46a3c-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c'),('0bd46494-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d'),('0bd46858-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','5b6c7d8e-9f0a-4b1c-8d2e-3f4a5b6c7d8e'),('0bd4627e-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','5d6e7f8a-9b0c-4d1e-bf2a-3b4c5d6e7f8a'),('0bd46074-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','6c7d8e9f-0a1b-4c2d-8e3f-4a5b6c7d8e9f'),('0bd46531-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','6d7e8f9a-0b1c-4d2e-9f3a-4b5c6d7e8f9a'),('0bd46c10-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','73875cfa-b134-4b6d-9a46-43ea4fe9b500'),('0bd467a6-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','7e8f9a0b-1c2d-4e3f-9a4b-5c6d7e8f9a0b'),('0bd469a1-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e'),('0bd465d9-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c'),('0bd463ee-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000002','9a0b1c2d-3e4f-4a5b-9c6d-7e8f9a0b1c2d'),('0bd4e56a-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f'),('0bd4e620-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d'),('0bd4e4b9-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c'),('0bd4e280-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d'),('0bd4dfce-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','7f8a9b0c-1d2e-4f3a-9b4c-5d6e7f8a9b0c'),('0bd4e406-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e'),('0bd4e34d-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000003','8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c'),('0bd4a82c-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f'),('0bd4a8d5-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d'),('0bd4a57c-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','2d3e4f5a-6b7c-4d8e-9f0a-0b1c2d3e4f5a'),('0bd4a625-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','2e3f4a5b-6c7d-4e8f-8a9b-0c1d2e3f4a5b'),('0bd4a778-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c'),('0bd4a377-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d'),('0bd4a2c8-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','5d6e7f8a-9b0c-4d1e-bf2a-3b4c5d6e7f8a'),('0bd4a1f6-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','6c7d8e9f-0a1b-4c2d-8e3f-4a5b6c7d8e9f'),('0bd4a422-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','6d7e8f9a-0b1c-4d2e-9f3a-4b5c6d7e8f9a'),('0bd49f8d-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','7f8a9b0c-1d2e-4f3a-9b4c-5d6e7f8a9b0c'),('0bd4a6cf-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e'),('0bd4a4ce-3e84-11f1-ab91-005056c00001','11111111-1111-1111-1111-000000000004','8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c'),('0bd42ccf-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','0b1c2d3e-4f5a-4b6c-9d7e-8f9a0b1c2d3e'),('0bd425d5-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','0e1f2a3b-4c5d-4e6f-8a7b-8c9d0e1f2a3b'),('0bd42937-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','1a2b3c4d-5e6f-4a7b-9c8d-0e1f2a3b3c4d'),('0bd42e67-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f'),('0bd4225e-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','1d2e3f4a-5b6c-4d7e-bf8a-9b0c1d2e3f4a'),('0bd42f02-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d'),('0bd42665-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','2c3d4e5f-6a7b-4c8d-ae0f-1a2b3c4d5e6f'),('0bd4254a-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','2d3e4f5a-6b7c-4d8e-9f0a-0b1c2d3e4f5a'),('0bd42a5e-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','2e3f4a5b-6c7d-4e8f-8a9b-0c1d2e3f4a5b'),('0bd42c49-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','3f2d8a1b-4c5e-4f6a-9b0c-1d2e3f4a5b6c'),('0bd42ddd-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c'),('0bd4238b-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d'),('0bd428a7-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','4c5d6e7f-8a9b-4c0d-8e1f-2a3b4c5d6e7f'),('0bd42f91-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','4f5a6b7c-8d9e-4f0a-8b1c-2d3e4f5a6b7c'),('0bd4281a-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','5b6c7d8e-9f0a-4b1c-8d2e-3f4a5b6c7d8e'),('0bd421d1-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','5d6e7f8a-9b0c-4d1e-bf2a-3b4c5d6e7f8a'),('0bd42bbb-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','6b7c8d9e-0f1a-4b2c-9d3e-4f5a6b7c8d9e'),('0bd420f3-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','6c7d8e9f-0a1b-4c2d-8e3f-4a5b6c7d8e9f'),('0bd4241c-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','6d7e8f9a-0b1c-4d2e-9f3a-4b5c6d7e8f9a'),('0bd4301c-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','73875cfa-b134-4b6d-9a46-43ea4fe9b500'),('0bd42ae6-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','7a8b9c0d-1e2f-4a3b-8c4d-5e6f7a8b9c0d'),('0bd42792-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','7e8f9a0b-1c2d-4e3f-9a4b-5c6d7e8f9a0b'),('0bd41f11-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','7f8a9b0c-1d2e-4f3a-9b4c-5d6e7f8a9b0c'),('0bd42d54-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e'),('0bd429c6-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','8e9f0a1b-2c3d-4e5f-9a6b-7c8d9e0f1a2b'),('0bd424b3-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c'),('0bd422ef-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','9a0b1c2d-3e4f-4a5b-9c6d-7e8f9a0b1c2d'),('0bd42700-3e84-11f1-ab91-005056c00001','77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab','9c0d1e2f-3a4b-4c5d-8e6f-7a8b9c0d1e2f');

INSERT INTO `locker` VALUES ('c46811ae-8fda-4b2c-9757-fc4d4e56ef08',NULL,1),('90c96690-e4ff-4711-b900-89be2fe9dcdb',NULL,2),('4fade1b3-d6b6-4635-b505-0820f31b41f2',NULL,3),('fbe9708e-3919-4b15-b59f-767821715847',NULL,4),('59adc1b0-5760-425f-81b7-3e61755cf3b0',NULL,5),('0c230337-3e84-11f1-ab91-005056c00001',NULL,6),('0c230546-3e84-11f1-ab91-005056c00001',NULL,7),('0c230599-3e84-11f1-ab91-005056c00001',NULL,8),('0c2305d8-3e84-11f1-ab91-005056c00001',NULL,9),('0c23060d-3e84-11f1-ab91-005056c00001',NULL,10),('0c230723-3e84-11f1-ab91-005056c00001',NULL,11),('0c23075c-3e84-11f1-ab91-005056c00001',NULL,12),('0c230794-3e84-11f1-ab91-005056c00001',NULL,13),('0c2307c1-3e84-11f1-ab91-005056c00001',NULL,14),('0c2307ee-3e84-11f1-ab91-005056c00001',NULL,15),('0c23081d-3e84-11f1-ab91-005056c00001',NULL,16),('0c23084e-3e84-11f1-ab91-005056c00001',NULL,17),('0c23087b-3e84-11f1-ab91-005056c00001',NULL,18),('0c2308aa-3e84-11f1-ab91-005056c00001',NULL,19),('0c2308d7-3e84-11f1-ab91-005056c00001',NULL,20);

SET FOREIGN_KEY_CHECKS = 1;
