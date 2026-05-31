CREATE TABLE USER (
    id_user CHAR(36) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_user VARCHAR(255) NOT NULL,
    first_name VARCHAR(150) NULL,
    last_name VARCHAR(150) NULL,
    upcycling_score INT NOT NULL DEFAULT 0,
    premium TINYINT(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (id_user)
);

CREATE TABLE ROLE (
    id_role CHAR(36) NOT NULL,
    name_role VARCHAR(100) NOT NULL UNIQUE,
    PRIMARY KEY (id_role)
);

CREATE TABLE DOCUMENT (
    id_document CHAR(36) NOT NULL,
    id_user CHAR(36) NOT NULL,
    category VARCHAR(120) NULL,
    link VARCHAR(2048) NULL,
    PRIMARY KEY (id_document),
    CONSTRAINT fk_document_user
      FOREIGN KEY (id_user) REFERENCES USER(id_user)
      ON DELETE CASCADE
);

CREATE TABLE ANNOUNCEMENT (
    id_announcement CHAR(36) NOT NULL,
    address_annoucement VARCHAR(255) NULL,
    city VARCHAR(120) NULL,
    postal VARCHAR(30) NULL,
    description_annoucement TEXT NULL,
    availability_date DATE NOT NULL,
    price INT NULL,
    request TINYINT(1) NOT NULL DEFAULT 0,
    state_annoucement VARCHAR(60) NULL,
    PRIMARY KEY (id_announcement)
);

CREATE TABLE LOCKER (
    id_locker CHAR(36) NOT NULL,
    access_code VARCHAR(120) NULL,
    PRIMARY KEY (id_locker)
);

CREATE TABLE PAYEMENT (
    id_payement CHAR(36) NOT NULL,
    amount_cents INT NOT NULL,
    currency CHAR(3) NOT NULL,
    provider_payement VARCHAR(60) NULL,
    provider_ref VARCHAR(255) NULL,
    status_payement VARCHAR(60) NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    paid_at DATETIME NULL,
    PRIMARY KEY (id_payement)
);

CREATE TABLE SUBSCRIPTION_PLANS (
    id_plan CHAR(36) NOT NULL,
    price_cents INT NOT NULL,
    interval_unit VARCHAR(30) NOT NULL,
    interval_count INT NOT NULL DEFAULT 1,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    PRIMARY KEY (id_plan)
);

CREATE TABLE SUBSCRIPTION (
    id_subscription CHAR(36) NOT NULL,
    start_timestamp DATETIME NULL,
    end_timestamp DATETIME NULL,
    cancelled TINYINT(1) NOT NULL DEFAULT 0,
    cancelled_at DATETIME NULL,
    PRIMARY KEY (id_subscription)
);

CREATE TABLE ADVICE (
    id_advice CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    tag VARCHAR(120) NULL,
    date_advice DATE NULL,
    PRIMARY KEY (id_advice)
);

CREATE TABLE TOPIC (
    id_topic CHAR(36) NOT NULL,
    PRIMARY KEY (id_topic)
);

CREATE TABLE POST (
    id_post CHAR(36) NOT NULL,
    PRIMARY KEY (id_post)
);

CREATE TABLE PROJECT (
    id_project CHAR(36) NOT NULL,
    start_date_project DATE NULL,
    end_date DATE NULL,
    capacity INT NULL,
    PRIMARY KEY (id_project)
);

CREATE TABLE FORMATION (
    id_formation CHAR(36) NOT NULL,
    PRIMARY KEY (id_formation)
);

CREATE TABLE EVENT (
    id_event CHAR(36) NOT NULL,
    PRIMARY KEY (id_event)
);  

CREATE TABLE TICKET (
    id_ticket CHAR(36) NOT NULL,
    timestamp_ticket DATETIME NULL,
    PRIMARY KEY (id_ticket)
);

CREATE TABLE NOTIFICATION (
    id_notification CHAR(36) NOT NULL,
    view TINYINT(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (id_notification)
);

/* ================================================================================ */

CREATE TABLE REPORT_ANNOUCEMENT (
    id_report_annoucement CHAR(36) NOT NULL,
    id_user CHAR(36) NULL,
    id_announcement CHAR(36) NULL,
    reason TEXT NULL,
    PRIMARY KEY (id_report_annoucement),
    CONSTRAINT fk_report_announcement_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_report_announcement
    FOREIGN KEY (id_announcement) REFERENCES ANNOUNCEMENT(id_announcement)
);

CREATE TABLE REPORT_TOPIC (
    id_report_topic CHAR(36) NOT NULL,
    id_user CHAR(36) NULL,
    id_topic CHAR(36) NULL,
    PRIMARY KEY (id_report_topic),
    CONSTRAINT fk_report_topic_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_report_topic
    FOREIGN KEY (id_topic) REFERENCES TOPIC(id_topic)
);

CREATE TABLE REPORT_POST (
    id_report_post CHAR(36) NOT NULL,
    id_user CHAR(36) NULL,
    id_post CHAR(36) NULL,
    PRIMARY KEY (id_report_post),
    CONSTRAINT fk_report_post_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_report_post
    FOREIGN KEY (id_post) REFERENCES POST(id_post)
);

CREATE TABLE USER_ROLE (
    id_user CHAR(36) NOT NULL,
    id_role CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_role),
    CONSTRAINT fk_user_role_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_role_role
    FOREIGN KEY (id_role) REFERENCES ROLE(id_role)
);

CREATE TABLE USER_ANNOUNCEMENT (
    id_user CHAR(36) NOT NULL,
    id_announcement CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_announcement),
    CONSTRAINT fk_user_announcement_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_announcement_announcement
    FOREIGN KEY (id_announcement) REFERENCES ANNOUNCEMENT(id_announcement)
);

CREATE TABLE ANNOUNCEMENT_LOCKER (
    id_announcement CHAR(36) NOT NULL,
    id_locker CHAR(36) NOT NULL UNIQUE,
    PRIMARY KEY (id_announcement),
    CONSTRAINT fk_announcement_locker_announcement
    FOREIGN KEY (id_announcement) REFERENCES ANNOUNCEMENT(id_announcement),
    CONSTRAINT fk_announcement_locker_locker
    FOREIGN KEY (id_locker) REFERENCES LOCKER(id_locker)
);

CREATE TABLE USER_ADVICE (
    id_user CHAR(36) NOT NULL,
    id_advice CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_advice),
    CONSTRAINT fk_user_advice_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_advice_advice
    FOREIGN KEY (id_advice) REFERENCES ADVICE(id_advice)
);

CREATE TABLE USER_TOPIC (
    id_user CHAR(36) NOT NULL,
    id_topic CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_topic),
    CONSTRAINT fk_user_topic_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_topic_topic
    FOREIGN KEY (id_topic) REFERENCES TOPIC(id_topic)
);

CREATE TABLE TOPIC_POST (
    id_topic CHAR(36) NOT NULL,
    id_post CHAR(36) NOT NULL,
    PRIMARY KEY (id_topic, id_post),
    CONSTRAINT fk_topic_post_topic
    FOREIGN KEY (id_topic) REFERENCES TOPIC(id_topic),
    CONSTRAINT fk_topic_post_post
    FOREIGN KEY (id_post) REFERENCES POST(id_post)
);

CREATE TABLE USER_PROJECT_CREATE (
    id_user CHAR(36) NOT NULL,
    id_project CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_project),
    CONSTRAINT fk_user_project_create_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_project_create_project
    FOREIGN KEY (id_project) REFERENCES PROJECT(id_project)
);

CREATE TABLE USER_PROJECT_INSCRIPTION (
    id_user CHAR(36) NOT NULL,
    id_project CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_project),
    CONSTRAINT fk_user_project_inscription_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_project_inscription_project
    FOREIGN KEY (id_project) REFERENCES PROJECT(id_project)
);

CREATE TABLE USER_PROJECT_UPDOWN (
    id_user CHAR(36) NOT NULL,
    id_project CHAR(36) NOT NULL,
    updown SMALLINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id_user, id_project),
    CONSTRAINT fk_user_project_updown_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_project_updown_project
    FOREIGN KEY (id_project) REFERENCES PROJECT(id_project)
);

CREATE TABLE USER_FORMATION_FORMATER (
    id_user CHAR(36) NOT NULL,
    id_formation CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_formation),
    CONSTRAINT fk_user_formation_formater_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_formation_formater_formation
    FOREIGN KEY (id_formation) REFERENCES FORMATION(id_formation)
);

CREATE TABLE USER_FORMATION_CREATE (
    id_user CHAR(36) NOT NULL,
    id_formation CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_formation),
    CONSTRAINT fk_user_formation_create_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_formation_create_formation
    FOREIGN KEY (id_formation) REFERENCES FORMATION(id_formation)
);

CREATE TABLE USER_FORMATION_INSCRIPTION (
    id_user CHAR(36) NOT NULL,
    id_formation CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_formation),
    CONSTRAINT fk_user_formation_inscription_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_formation_inscription_formation
    FOREIGN KEY (id_formation) REFERENCES FORMATION(id_formation)
);

CREATE TABLE USER_FORMATION_INSCRIPTION_PAYEMENT (
    id_user CHAR(36) NOT NULL,
    id_formation CHAR(36) NOT NULL,
    id_payement CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_formation, id_payement),
    CONSTRAINT fk_ufip_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_ufip_formation
    FOREIGN KEY (id_formation) REFERENCES FORMATION(id_formation),
    CONSTRAINT fk_ufip_payement
    FOREIGN KEY (id_payement) REFERENCES PAYEMENT(id_payement)
);

CREATE TABLE USER_EVENT (
    id_user CHAR(36) NOT NULL,
    id_event CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_event),
    CONSTRAINT fk_user_event_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_event_event
    FOREIGN KEY (id_event) REFERENCES EVENT(id_event)
);

CREATE TABLE USER_EVENT_INSCRIPTION (
    id_user CHAR(36) NOT NULL,
    id_event CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_event),
    CONSTRAINT fk_user_event_inscription_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_event_inscription_event
    FOREIGN KEY (id_event) REFERENCES EVENT(id_event)
);

CREATE TABLE USER_TICKET (
    id_user CHAR(36) NOT NULL,
    id_ticket CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_ticket),
    CONSTRAINT fk_user_ticket_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_ticket_ticket
    FOREIGN KEY (id_ticket) REFERENCES TICKET(id_ticket)
);

CREATE TABLE USER_NOTIFICATION (
    id_user CHAR(36) NOT NULL,
    id_notification CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_notification),
    CONSTRAINT fk_user_notification_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_notification_notification
    FOREIGN KEY (id_notification) REFERENCES NOTIFICATION(id_notification)
);

CREATE TABLE USER_PAYEMENT (
    id_user CHAR(36) NOT NULL,
    id_payement CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_payement),
    CONSTRAINT fk_user_payement_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_payement_payement
    FOREIGN KEY (id_payement) REFERENCES PAYEMENT(id_payement)
);

CREATE TABLE PAYEMENT_FORMATION (
    id_payement CHAR(36) NOT NULL,
    id_formation CHAR(36) NOT NULL,
    PRIMARY KEY (id_payement, id_formation),
    CONSTRAINT fk_payement_formation_payement
    FOREIGN KEY (id_payement) REFERENCES PAYEMENT(id_payement),
    CONSTRAINT fk_payement_formation_formation
    FOREIGN KEY (id_formation) REFERENCES FORMATION(id_formation)
);

CREATE TABLE PAYEMENT_EVENT (
    id_payement CHAR(36) NOT NULL,
    id_event CHAR(36) NOT NULL,
    PRIMARY KEY (id_payement, id_event),
    CONSTRAINT fk_payement_event_payement
    FOREIGN KEY (id_payement) REFERENCES PAYEMENT(id_payement),
    CONSTRAINT fk_payement_event_event
    FOREIGN KEY (id_event) REFERENCES EVENT(id_event)
);

CREATE TABLE PAYEMENT_PROJECT (
    id_payement CHAR(36) NOT NULL,
    id_project CHAR(36) NOT NULL,
    PRIMARY KEY (id_payement, id_project),
    CONSTRAINT fk_payement_project_payement
    FOREIGN KEY (id_payement) REFERENCES PAYEMENT(id_payement),
    CONSTRAINT fk_payement_project_project
    FOREIGN KEY (id_project) REFERENCES PROJECT(id_project)
);

CREATE TABLE USER_SUBSCRIPTION (
    id_user CHAR(36) NOT NULL,
    id_subscription CHAR(36) NOT NULL,
    PRIMARY KEY (id_user, id_subscription),
    CONSTRAINT fk_user_subscription_user
    FOREIGN KEY (id_user) REFERENCES USER(id_user),
    CONSTRAINT fk_user_subscription_subscription
    FOREIGN KEY (id_subscription) REFERENCES SUBSCRIPTION(id_subscription)
);

CREATE TABLE SUB_SUB_PLAN (
    id_subscription CHAR(36) NOT NULL,
    id_plan CHAR(36) NOT NULL,
    PRIMARY KEY (id_subscription, id_plan),
    CONSTRAINT fk_sub_sub_plan_subscription
    FOREIGN KEY (id_subscription) REFERENCES SUBSCRIPTION(id_subscription),
    CONSTRAINT fk_sub_sub_plan_plan
    FOREIGN KEY (id_plan) REFERENCES SUBSCRIPTION_PLANS(id_plan)
);