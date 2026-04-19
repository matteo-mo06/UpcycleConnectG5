CREATE TABLE PROFESSIONAL_REQUEST (
    id_request CHAR(36) NOT NULL,
    id_user CHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME NULL,
    PRIMARY KEY (id_request),
    CONSTRAINT fk_professional_request_user
        FOREIGN KEY (id_user) REFERENCES USER(id_user)
);

INSERT INTO ROLE (id_role, name_role) VALUES
('77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab', 'admin'),
('11111111-1111-1111-1111-000000000002', 'artisan'),
('11111111-1111-1111-1111-000000000003', 'user'),
('11111111-1111-1111-1111-000000000004', 'salarie');