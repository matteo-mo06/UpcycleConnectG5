INSERT INTO `permission` (id_permission, name_permission, domain) VALUES
    ('a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'access_admin',   'espaces'),
    ('b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'access_artisan', 'espaces'),
    ('c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'access_salarie', 'espaces');

INSERT INTO ROLE_PERMISSION (id, role_id, permission_id)
SELECT UUID(), '77e2bd9a-a1af-44f7-8f0e-edfb97f9c0ab', id_permission
FROM `permission` WHERE name_permission = 'access_admin';

INSERT INTO ROLE_PERMISSION (id, role_id, permission_id)
SELECT UUID(), '11111111-1111-1111-1111-000000000002', id_permission
FROM `permission` WHERE name_permission = 'access_artisan';

INSERT INTO ROLE_PERMISSION (id, role_id, permission_id)
SELECT UUID(), '11111111-1111-1111-1111-000000000004', id_permission
FROM `permission` WHERE name_permission = 'access_salarie';
