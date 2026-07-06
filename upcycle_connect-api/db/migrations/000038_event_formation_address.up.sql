ALTER TABLE `event`
    ADD COLUMN `address_event` varchar(255) DEFAULT NULL,
    ADD COLUMN `city_event`    varchar(120) DEFAULT NULL,
    ADD COLUMN `postal_event`  varchar(30)  DEFAULT NULL;

UPDATE `event` SET `address_event` = `location_event` WHERE `location_event` IS NOT NULL;

ALTER TABLE `event` DROP COLUMN `location_event`;

ALTER TABLE `formation`
    ADD COLUMN `address_formation` varchar(255) DEFAULT NULL,
    ADD COLUMN `city_formation`    varchar(120) DEFAULT NULL,
    ADD COLUMN `postal_formation`  varchar(30)  DEFAULT NULL;

UPDATE `formation` SET `address_formation` = `location_formation` WHERE `location_formation` IS NOT NULL;

ALTER TABLE `formation` DROP COLUMN `location_formation`;
