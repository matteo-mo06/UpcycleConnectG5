SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS USER_SANCTIONS;

ALTER TABLE ANNOUNCEMENT DROP COLUMN deleted_at;
ALTER TABLE TOPIC        DROP COLUMN deleted_at;
ALTER TABLE POST         DROP COLUMN deleted_at;

ALTER TABLE REPORT
    DROP FOREIGN KEY fk_report_reported_user,
    DROP FOREIGN KEY fk_report_resolved_by,
    DROP COLUMN id_reported_user,
    DROP COLUMN resolved_by,
    DROP COLUMN action_taken,
    DROP COLUMN resolved_at;

SET FOREIGN_KEY_CHECKS = 1;
