DROP TABLE IF EXISTS featured_slot_log;
ALTER TABLE ANNOUNCEMENT
    DROP COLUMN is_featured,
    DROP COLUMN featured_until,
    DROP COLUMN featured_requested_at;
