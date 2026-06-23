CREATE TABLE `platform_settings` (
  `key_setting` varchar(100) NOT NULL,
  `value_setting` varchar(255) NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`key_setting`)
);

INSERT INTO `platform_settings` (`key_setting`, `value_setting`) VALUES ('commission_rate', '5.00');

ALTER TABLE `payement`
  ADD COLUMN `announcement_id` char(36) DEFAULT NULL,
  ADD COLUMN `buyer_id` char(36) DEFAULT NULL,
  ADD COLUMN `commission_amount_cents` int NOT NULL DEFAULT 0;
