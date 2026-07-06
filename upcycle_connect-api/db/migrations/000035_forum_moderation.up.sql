ALTER TABLE `topic`
  ADD COLUMN `moderated_by` CHAR(36) NULL DEFAULT NULL,
  ADD COLUMN `moderation_message` TEXT NULL,
  ADD KEY `fk_topic_moderated_by` (`moderated_by`),
  ADD CONSTRAINT `fk_topic_moderated_by` FOREIGN KEY (`moderated_by`) REFERENCES `user` (`id_user`) ON DELETE SET NULL;

ALTER TABLE `post`
  ADD COLUMN `moderated_by` CHAR(36) NULL DEFAULT NULL,
  ADD COLUMN `moderation_message` TEXT NULL,
  ADD KEY `fk_post_moderated_by` (`moderated_by`),
  ADD CONSTRAINT `fk_post_moderated_by` FOREIGN KEY (`moderated_by`) REFERENCES `user` (`id_user`) ON DELETE SET NULL;
