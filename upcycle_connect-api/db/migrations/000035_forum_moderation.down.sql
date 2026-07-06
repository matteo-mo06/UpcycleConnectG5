ALTER TABLE `topic`
  DROP FOREIGN KEY `fk_topic_moderated_by`,
  DROP KEY `fk_topic_moderated_by`,
  DROP COLUMN `moderated_by`,
  DROP COLUMN `moderation_message`;

ALTER TABLE `post`
  DROP FOREIGN KEY `fk_post_moderated_by`,
  DROP KEY `fk_post_moderated_by`,
  DROP COLUMN `moderated_by`,
  DROP COLUMN `moderation_message`;
