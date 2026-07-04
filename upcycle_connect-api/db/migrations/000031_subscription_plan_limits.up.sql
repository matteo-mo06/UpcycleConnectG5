ALTER TABLE subscription_plans
  ADD COLUMN max_announcements_per_month INT NULL,
  ADD COLUMN max_projects_per_month INT NULL,
  ADD COLUMN max_features_per_month INT NULL;

UPDATE subscription_plans
SET name = 'Medium',
    max_announcements_per_month = 10,
    max_projects_per_month = 5,
    max_features_per_month = 1
WHERE name = 'Abonnement Mensuel';
