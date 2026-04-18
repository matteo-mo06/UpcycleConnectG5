SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE REPORT;

CREATE TABLE REPORT_ANNOUCEMENT (
    id_report_annoucement CHAR(36) NOT NULL,
    id_user               CHAR(36) NULL,
    id_announcement       CHAR(36) NULL,
    reason                TEXT     NULL,
    PRIMARY KEY (id_report_annoucement),
    CONSTRAINT fk_report_announcement_user FOREIGN KEY (id_user)         REFERENCES USER(id_user),
    CONSTRAINT fk_report_announcement      FOREIGN KEY (id_announcement) REFERENCES ANNOUNCEMENT(id_announcement)
);

CREATE TABLE REPORT_TOPIC (
    id_report_topic CHAR(36) NOT NULL,
    id_user         CHAR(36) NULL,
    id_topic        CHAR(36) NULL,
    PRIMARY KEY (id_report_topic),
    CONSTRAINT fk_report_topic_user FOREIGN KEY (id_user)   REFERENCES USER(id_user),
    CONSTRAINT fk_report_topic      FOREIGN KEY (id_topic)  REFERENCES TOPIC(id_topic)
);

CREATE TABLE REPORT_POST (
    id_report_post CHAR(36) NOT NULL,
    id_user        CHAR(36) NULL,
    id_post        CHAR(36) NULL,
    PRIMARY KEY (id_report_post),
    CONSTRAINT fk_report_post_user FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_report_post      FOREIGN KEY (id_post) REFERENCES POST(id_post)
);

SET FOREIGN_KEY_CHECKS = 1;
