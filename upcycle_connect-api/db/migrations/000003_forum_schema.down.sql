DROP TABLE IF EXISTS `user_topic`;

ALTER TABLE `post`
  DROP FOREIGN KEY `fk_post_parent`,
  DROP FOREIGN KEY `fk_post_author`,
  DROP KEY `fk_post_parent`,
  DROP KEY `fk_post_author`,
  DROP COLUMN `created_at`,
  DROP COLUMN `id_parent_post`,
  DROP COLUMN `id_author`,
  DROP COLUMN `body_post`;

ALTER TABLE `topic`
  DROP FOREIGN KEY `fk_topic_author`,
  DROP KEY `fk_topic_author`,
  DROP COLUMN `created_at`,
  DROP COLUMN `id_author`,
  DROP COLUMN `title_topic`;
