CREATE TABLE `score_action` (
  `action_type` varchar(60) NOT NULL,
  `points`      int         NOT NULL,
  `description` varchar(255) NOT NULL,
  PRIMARY KEY (`action_type`)
);

CREATE TABLE `score_log` (
  `id`           char(36)    NOT NULL,
  `id_user`      char(36)    NOT NULL,
  `action_type`  varchar(60) NOT NULL,
  `id_reference` char(36)    NOT NULL,
  `points`       int         NOT NULL,
  `created_at`   datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_score_log` (`id_user`, `action_type`, `id_reference`),
  KEY `fk_score_log_user` (`id_user`),
  KEY `fk_score_log_action` (`action_type`),
  CONSTRAINT `fk_score_log_user`   FOREIGN KEY (`id_user`)     REFERENCES `user`         (`id_user`),
  CONSTRAINT `fk_score_log_action` FOREIGN KEY (`action_type`) REFERENCES `score_action` (`action_type`)
);

INSERT INTO `score_action` (`action_type`, `points`, `description`) VALUES
  ('project_created', 50, 'Création d''un projet collectif'),
  ('formation_registration', 30, 'Inscription à une formation'),
  ('project_registration', 25, 'Inscription à un projet collectif'),
  ('event_registration', 15, 'Inscription à un événement'),
  ('deposit_validated', 20, 'Dépôt d''objet physiquement validé'),
  ('announcement_sold', 15, 'Objet récupéré — côté vendeur'),
  ('announcement_bought', 10, 'Objet récupéré — côté acheteur'),
  ('sanction_warning', -15, 'Avertissement reçu'),
  ('sanction_suspension', -40, 'Suspension reçue');
