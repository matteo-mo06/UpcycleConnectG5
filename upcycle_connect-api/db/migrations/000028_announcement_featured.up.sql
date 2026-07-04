ALTER TABLE ANNOUNCEMENT
    ADD COLUMN is_featured TINYINT(1) NOT NULL DEFAULT 0,
    ADD COLUMN featured_until DATETIME NULL,
    ADD COLUMN featured_requested_at DATETIME NULL;

CREATE TABLE featured_slot_log (
    id_slot CHAR(36) NOT NULL,
    id_announcement CHAR(36) NOT NULL,
    started_at DATETIME NOT NULL,
    ends_at DATETIME NOT NULL,
    PRIMARY KEY (id_slot),
    CONSTRAINT fk_fsl_announcement FOREIGN KEY (id_announcement)
        REFERENCES ANNOUNCEMENT(id_announcement) ON DELETE CASCADE
);
