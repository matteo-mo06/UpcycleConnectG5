SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE USER_FORMATION_FORMATER (
    id_user      CHAR(36) NOT NULL,
    id_formation CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_formation),
    CONSTRAINT fk_user_formation_formater_user
        FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_formation_formater_formation
        FOREIGN KEY (id_formation) REFERENCES FORMATION(id_formation)
);

CREATE TABLE USER_FORMATION_CREATE (
    id_user      CHAR(36) NOT NULL,
    id_formation CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_formation),
    CONSTRAINT fk_user_formation_create_user
        FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_formation_create_formation
        FOREIGN KEY (id_formation) REFERENCES FORMATION(id_formation)
);

ALTER TABLE FORMATION
    DROP FOREIGN KEY fk_formation_formateur,
    DROP FOREIGN KEY fk_formation_creator,
    DROP COLUMN id_formateur,
    DROP COLUMN id_creator;

CREATE TABLE USER_PROJECT_CREATE (
    id_user    CHAR(36) NOT NULL,
    id_project CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_project),
    CONSTRAINT fk_user_project_create_user
        FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_project_create_project
        FOREIGN KEY (id_project) REFERENCES PROJECT(id_project)
);

ALTER TABLE PROJECT
    DROP FOREIGN KEY fk_project_creator,
    DROP COLUMN id_creator;

CREATE TABLE USER_EVENT (
    id_user  CHAR(36) NOT NULL,
    id_event CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_event),
    CONSTRAINT fk_user_event_user
        FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_event_event
        FOREIGN KEY (id_event) REFERENCES EVENT(id_event)
);

ALTER TABLE EVENT
    DROP FOREIGN KEY fk_event_creator,
    DROP COLUMN id_creator;

SET FOREIGN_KEY_CHECKS = 1;