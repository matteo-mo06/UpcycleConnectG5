CREATE TABLE CATEGORY (
    id_category CHAR(36) NOT NULL,
    name_category VARCHAR(100) NOT NULL UNIQUE,
    description_category TEXT NULL,
    PRIMARY KEY (id_category)
);

INSERT INTO CATEGORY (id_category, name_category, description_category) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-000000000001', 'Couture & Textile', 'Retouches, créations et recyclage textile'),
('aaaaaaaa-aaaa-aaaa-aaaa-000000000002', 'Menuiserie & Bois', 'Réparation et création de meubles en bois'),
('aaaaaaaa-aaaa-aaaa-aaaa-000000000003', 'Électronique', 'Réparation d''appareils électroniques'),
('aaaaaaaa-aaaa-aaaa-aaaa-000000000004', 'Maroquinerie', 'Réparation et création en cuir'),
('aaaaaaaa-aaaa-aaaa-aaaa-000000000005', 'Céramique & Poterie', 'Création et restauration céramique');

ALTER TABLE ANNOUNCEMENT
    ADD COLUMN id_category CHAR(36) NULL,
    ADD COLUMN title_announcement VARCHAR(255) NULL,
    ADD CONSTRAINT fk_announcement_category
        FOREIGN KEY (id_category) REFERENCES CATEGORY(id_category);

UPDATE ANNOUNCEMENT
SET id_category = 'aaaaaaaa-aaaa-aaaa-aaaa-000000000001', title_announcement = 'Test Announcement One'
WHERE id_announcement = 'TEST-ANN-0001';

UPDATE ANNOUNCEMENT
SET id_category = 'aaaaaaaa-aaaa-aaaa-aaaa-000000000002', title_announcement = 'Test Announcement Two'
WHERE id_announcement = 'TEST-ANN-0002';

ALTER TABLE ANNOUNCEMENT
    MODIFY COLUMN id_category CHAR(36) NOT NULL,
    MODIFY COLUMN title_announcement VARCHAR(255) NOT NULL;

ALTER TABLE EVENT
    ADD COLUMN title_event VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN description_event TEXT NULL,
    ADD COLUMN date_event DATETIME NULL,
    ADD COLUMN location_event VARCHAR(255) NULL,
    ADD COLUMN capacity INT NULL,
    ADD COLUMN price_cents INT NOT NULL DEFAULT 0;
