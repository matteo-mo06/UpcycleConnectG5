DELETE FROM `subscription_plans` WHERE `stripe_price_id` IS NULL AND `name` IN ('Premium Mensuel', 'Premium Annuel');

ALTER TABLE `subscription`
  DROP COLUMN `stripe_customer_id`,
  DROP COLUMN `stripe_subscription_id`;

ALTER TABLE `subscription_plans`
  DROP COLUMN `stripe_price_id`,
  DROP COLUMN `name`;
