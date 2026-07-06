ALTER TABLE `event`
    ADD COLUMN `location_event` varchar(255) DEFAULT NULL;

UPDATE `event` SET `location_event` = `address_event`;

ALTER TABLE `event`
    DROP COLUMN `address_event`,
    DROP COLUMN `city_event`,
    DROP COLUMN `postal_event`;

ALTER TABLE `formation`
    ADD COLUMN `location_formation` varchar(255) DEFAULT NULL;

UPDATE `formation` SET `location_formation` = `address_formation`;

ALTER TABLE `formation`
    DROP COLUMN `address_formation`,
    DROP COLUMN `city_formation`,
    DROP COLUMN `postal_formation`;
