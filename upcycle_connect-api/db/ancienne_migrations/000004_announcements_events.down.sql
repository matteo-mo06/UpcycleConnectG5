SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE ANNOUNCEMENT
    DROP FOREIGN KEY fk_announcement_category,
    DROP COLUMN id_category,
    DROP COLUMN title_announcement;

ALTER TABLE EVENT
    DROP COLUMN title_event,
    DROP COLUMN description_event,
    DROP COLUMN date_event,
    DROP COLUMN location_event,
    DROP COLUMN capacity,
    DROP COLUMN price_cents;

DROP TABLE IF EXISTS CATEGORY;

SET FOREIGN_KEY_CHECKS = 1;
