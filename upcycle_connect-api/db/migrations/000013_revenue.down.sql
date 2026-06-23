DROP TABLE IF EXISTS `platform_settings`;

ALTER TABLE `payement`
  DROP COLUMN `announcement_id`,
  DROP COLUMN `buyer_id`,
  DROP COLUMN `commission_amount_cents`;
