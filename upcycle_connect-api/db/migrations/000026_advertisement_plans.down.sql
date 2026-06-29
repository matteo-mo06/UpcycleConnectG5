ALTER TABLE advertisement DROP FOREIGN KEY fk_ad_plan;
ALTER TABLE advertisement DROP COLUMN plan_id;
DROP TABLE IF EXISTS advertisement_plans;
