ALTER TABLE `topic`
  ADD COLUMN `title_topic` varchar(255) NOT NULL,
  ADD COLUMN `id_author` char(36) NOT NULL,
  ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ADD KEY `fk_topic_author` (`id_author`),
  ADD CONSTRAINT `fk_topic_author` FOREIGN KEY (`id_author`) REFERENCES `user` (`id_user`);

ALTER TABLE `post`
  ADD COLUMN `body_post` text NOT NULL,
  ADD COLUMN `id_author` char(36) NOT NULL,
  ADD COLUMN `id_parent_post` char(36) DEFAULT NULL,
  ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ADD KEY `fk_post_author` (`id_author`),
  ADD KEY `fk_post_parent` (`id_parent_post`),
  ADD CONSTRAINT `fk_post_author` FOREIGN KEY (`id_author`) REFERENCES `user` (`id_user`),
  ADD CONSTRAINT `fk_post_parent` FOREIGN KEY (`id_parent_post`) REFERENCES `post` (`id_post`);

CREATE TABLE `user_topic` (
  `id_user` char(36) NOT NULL,
  `id_topic` char(36) NOT NULL,
  PRIMARY KEY (`id_user`,`id_topic`),
  KEY `fk_user_topic_topic` (`id_topic`),
  CONSTRAINT `fk_user_topic_topic` FOREIGN KEY (`id_topic`) REFERENCES `topic` (`id_topic`),
  CONSTRAINT `fk_user_topic_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
