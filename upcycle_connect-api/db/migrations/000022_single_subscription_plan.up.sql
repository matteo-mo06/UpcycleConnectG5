DELETE FROM `subscription_plans` WHERE `name` != 'Abonnement Mensuel';
UPDATE `subscription_plans` SET `name` = 'Abonnement Mensuel' WHERE `name` = 'Premium Mensuel';
