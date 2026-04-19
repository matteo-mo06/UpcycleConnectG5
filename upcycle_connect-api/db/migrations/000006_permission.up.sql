CREATE TABLE PERMISSION (
    id         CHAR(36)     NOT NULL,
    name       VARCHAR(100) NOT NULL UNIQUE,
    domain     VARCHAR(50)  NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE ROLE_PERMISSION (
    id            CHAR(36) NOT NULL,
    role_id       CHAR(36) NOT NULL,
    permission_id CHAR(36) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_role_permission (role_id, permission_id),
    CONSTRAINT fk_rp_role
        FOREIGN KEY (role_id) REFERENCES ROLE(id_role),
    CONSTRAINT fk_rp_permission
        FOREIGN KEY (permission_id) REFERENCES PERMISSION(id)
);

INSERT INTO PERMISSION (id, name, domain) VALUES
    ('3f2d8a1b-4c5e-4f6a-9b0c-1d2e3f4a5b6c', 'manage_users',                 'administration'),
    ('7a8b9c0d-1e2f-4a3b-8c4d-5e6f7a8b9c0d', 'manage_roles',                 'administration'),
    ('2c3d4e5f-6a7b-4c8d-ae0f-1a2b3c4d5e6f', 'manage_categories',            'administration'),
    ('8e9f0a1b-2c3d-4e5f-9a6b-7c8d9e0f1a2b', 'manage_professional_requests', 'administration'),
    ('4f5a6b7c-8d9e-4f0a-8b1c-2d3e4f5a6b7c', 'send_notifications',           'administration'),
    ('1a2b3c4d-5e6f-4a7b-9c8d-0e1f2a3b3c4d', 'manage_payments',              'administration'),
    ('6b7c8d9e-0f1a-4b2c-9d3e-4f5a6b7c8d9e', 'manage_subscriptions',         'administration'),
    ('9c0d1e2f-3a4b-4c5d-8e6f-7a8b9c0d1e2f', 'manage_documents',             'administration');

INSERT INTO PERMISSION (id, name, domain) VALUES
    ('5d6e7f8a-9b0c-4d1e-bf2a-3b4c5d6e7f8a', 'create_announcement',          'announcements'),
    ('0e1f2a3b-4c5d-4e6f-8a7b-8c9d0e1f2a3b', 'manage_announcements',         'announcements'),
    ('7f8a9b0c-1d2e-4f3a-9b4c-5d6e7f8a9b0c', 'buy_announcement',             'announcements');

INSERT INTO PERMISSION (id, name, domain) VALUES
    ('2a3b4c5d-6e7f-4a8b-8c9d-0e1f2a3b4c5d', 'request_locker_deposit',       'lockers'),
    ('8b9c0d1e-2f3a-4b4c-9d5e-6f7a8b9c0d1e', 'pickup_locker',                'lockers'),
    ('4c5d6e7f-8a9b-4c0d-8e1f-2a3b4c5d6e7f', 'manage_lockers',               'lockers');

INSERT INTO PERMISSION (id, name, domain) VALUES
    ('1d2e3f4a-5b6c-4d7e-bf8a-9b0c1d2e3f4a', 'create_event',                 'events'),
    ('7e8f9a0b-1c2d-4e3f-9a4b-5c6d7e8f9a0b', 'manage_events',                'events'),
    ('3f4a5b6c-7d8e-4f9a-8b0c-1d2e3f4a5b6c', 'register_event',               'events');

INSERT INTO PERMISSION (id, name, domain) VALUES
    ('9a0b1c2d-3e4f-4a5b-9c6d-7e8f9a0b1c2d', 'create_formation',             'formations'),
    ('5b6c7d8e-9f0a-4b1c-8d2e-3f4a5b6c7d8e', 'manage_formations',            'formations'),
    ('1c2d3e4f-5a6b-4c7d-8e8f-9a0b1c2d3e4f', 'register_formation',           'formations');

INSERT INTO PERMISSION (id, name, domain) VALUES
    ('6d7e8f9a-0b1c-4d2e-9f3a-4b5c6d7e8f9a', 'create_project',               'projects'),
    ('2e3f4a5b-6c7d-4e8f-8a9b-0c1d2e3f4a5b', 'manage_projects',              'projects');

INSERT INTO PERMISSION (id, name, domain) VALUES
    ('8f9a0b1c-2d3e-4f4a-9b5c-6d7e8f9a0b1c', 'create_topic',                 'forum'),
    ('4a5b6c7d-8e9f-4a0b-8c1d-2e3f4a5b6c7d', 'create_post',                  'forum'),
    ('0b1c2d3e-4f5a-4b6c-9d7e-8f9a0b1c2d3e', 'moderate_forum',               'forum');

INSERT INTO PERMISSION (id, name, domain) VALUES
    ('6c7d8e9f-0a1b-4c2d-8e3f-4a5b6c7d8e9f', 'create_advice',                'conseils'),
    ('2d3e4f5a-6b7c-4d8e-9f0a-0b1c2d3e4f5a', 'manage_advice',                'conseils');

INSERT INTO ROLE_PERMISSION (id, role_id, permission_id)
SELECT UUID(), '77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab', id FROM PERMISSION;

INSERT INTO ROLE_PERMISSION (id, role_id, permission_id)
SELECT UUID(), '11111111-1111-1111-1111-000000000002', id
FROM PERMISSION
WHERE name IN (
    'create_announcement',
    'manage_announcements',
    'create_event',
    'manage_events',
    'register_event',
    'create_formation',
    'manage_formations',
    'register_formation',
    'request_locker_deposit',
    'pickup_locker',
    'create_project',
    'manage_projects',
    'create_topic',
    'create_post',
    'create_advice',
    'manage_advice'
);

INSERT INTO ROLE_PERMISSION (id, role_id, permission_id)
SELECT UUID(), '11111111-1111-1111-1111-000000000004', id
FROM PERMISSION
WHERE name IN (
    'create_announcement',
    'buy_announcement',
    'register_event',
    'register_formation',
    'request_locker_deposit',
    'pickup_locker',
    'create_project',
    'manage_projects',
    'create_topic',
    'create_post',
    'create_advice',
    'manage_advice'
);

INSERT INTO ROLE_PERMISSION (id, role_id, permission_id)
SELECT UUID(), '11111111-1111-1111-1111-000000000003', id
FROM PERMISSION
WHERE name IN (
    'buy_announcement',
    'register_event',
    'register_formation',
    'request_locker_deposit',
    'pickup_locker',
    'create_topic',
    'create_post'
);
