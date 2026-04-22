SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE REPORT_ANNOUCEMENT;
DROP TABLE REPORT_TOPIC;
DROP TABLE REPORT_POST;

CREATE TABLE REPORT (
    id_report       CHAR(36)    NOT NULL,
    id_user         CHAR(36)    NULL,
    id_announcement CHAR(36)    NULL,
    id_topic        CHAR(36)    NULL,
    id_post         CHAR(36)    NULL,
    reason          TEXT        NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at      DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id_report),
    CONSTRAINT fk_report_user         FOREIGN KEY (id_user)         REFERENCES USER(id_user),
    CONSTRAINT fk_report_announcement FOREIGN KEY (id_announcement) REFERENCES ANNOUNCEMENT(id_announcement),
    CONSTRAINT fk_report_topic        FOREIGN KEY (id_topic)        REFERENCES TOPIC(id_topic),
    CONSTRAINT fk_report_post         FOREIGN KEY (id_post)         REFERENCES POST(id_post)
);

SET FOREIGN_KEY_CHECKS = 1;
