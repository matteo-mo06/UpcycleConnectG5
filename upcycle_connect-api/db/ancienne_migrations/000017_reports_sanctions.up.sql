SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE REPORT
    ADD COLUMN id_reported_user CHAR(36)    NULL,
    ADD COLUMN resolved_by      CHAR(36)    NULL,
    ADD COLUMN action_taken     VARCHAR(50) NULL,
    ADD COLUMN resolved_at      DATETIME    NULL;

ALTER TABLE REPORT
    ADD CONSTRAINT fk_report_reported_user FOREIGN KEY (id_reported_user) REFERENCES USER(id_user),
    ADD CONSTRAINT fk_report_resolved_by   FOREIGN KEY (resolved_by)      REFERENCES USER(id_user);

ALTER TABLE ANNOUNCEMENT ADD COLUMN deleted_at DATETIME NULL;
ALTER TABLE TOPIC        ADD COLUMN deleted_at DATETIME NULL;
ALTER TABLE POST         ADD COLUMN deleted_at DATETIME NULL;

CREATE TABLE USER_SANCTIONS (
    id_sanction   CHAR(36)    NOT NULL,
    id_user       CHAR(36)    NOT NULL,
    id_admin      CHAR(36)    NOT NULL,
    id_report     CHAR(36)    NULL,
    type          VARCHAR(20) NOT NULL,
    reason        TEXT        NULL,
    duration_days INT         NULL,
    expires_at    DATETIME    NULL,
    created_at    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_sanction),
    CONSTRAINT fk_sanctions_user   FOREIGN KEY (id_user)   REFERENCES USER(id_user),
    CONSTRAINT fk_sanctions_admin  FOREIGN KEY (id_admin)  REFERENCES USER(id_user),
    CONSTRAINT fk_sanctions_report FOREIGN KEY (id_report) REFERENCES REPORT(id_report)
);

SET FOREIGN_KEY_CHECKS = 1;
