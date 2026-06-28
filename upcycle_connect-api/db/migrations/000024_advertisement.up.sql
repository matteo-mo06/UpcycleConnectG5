CREATE TABLE `advertisement` (
  `id_advertisement` char(36) NOT NULL,
  `id_user` char(36) NOT NULL,
  `title` varchar(200) NOT NULL,
  `image_url` varchar(500) NOT NULL,
  `link_url` varchar(500) DEFAULT NULL,
  `state` enum('pending','approved','rejected','active','expired') NOT NULL DEFAULT 'pending',
  `price_cents` int NOT NULL,
  `rejection_reason` varchar(500) DEFAULT NULL,
  `stripe_checkout_session_id` varchar(200) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `approved_at` timestamp NULL DEFAULT NULL,
  `paid_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id_advertisement`),
  KEY `idx_ad_user` (`id_user`),
  KEY `idx_ad_state` (`state`),
  CONSTRAINT `fk_ad_user` FOREIGN KEY (`id_user`) REFERENCES `user` (`id_user`)
);

INSERT INTO `platform_settings` (`key_setting`, `value_setting`) VALUES ('advertisement_price_cents', '50000');
