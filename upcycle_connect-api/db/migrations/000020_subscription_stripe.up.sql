ALTER TABLE `subscription_plans`
  ADD COLUMN `name` varchar(100) NOT NULL DEFAULT '' AFTER `id_plan`,
  ADD COLUMN `stripe_price_id` varchar(100) DEFAULT NULL;

ALTER TABLE `subscription`
  ADD COLUMN `stripe_subscription_id` varchar(100) DEFAULT NULL,
  ADD COLUMN `stripe_customer_id` varchar(100) DEFAULT NULL;

INSERT INTO `subscription_plans` (`id_plan`, `name`, `price_cents`, `interval_unit`, `interval_count`, `is_active`, `stripe_price_id`)
VALUES
  (UUID(), 'Abonnement Mensuel', 1500, 'month', 1, 1, NULL);
