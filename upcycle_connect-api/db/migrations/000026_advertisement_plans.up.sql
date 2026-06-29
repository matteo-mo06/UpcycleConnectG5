CREATE TABLE advertisement_plans (
    id INT AUTO_INCREMENT PRIMARY KEY,
    weeks INT NOT NULL,
    price_cents INT NOT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1
);

INSERT INTO advertisement_plans (weeks, price_cents) VALUES (2, 20000), (4, 40000), (6, 60000);

ALTER TABLE advertisement ADD COLUMN plan_id INT NULL;
ALTER TABLE advertisement ADD CONSTRAINT fk_ad_plan FOREIGN KEY (plan_id) REFERENCES advertisement_plans(id);
