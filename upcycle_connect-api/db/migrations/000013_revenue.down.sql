DROP TABLE IF EXISTS `platform_settings`;

ALTER TABLE `payement`
  DROP COLUMN IF EXISTS `announcement_id`,
  DROP COLUMN IF EXISTS `buyer_id`,
  DROP COLUMN IF EXISTS `commission_amount_cents`;
