DELETE rp FROM ROLE_PERMISSION rp
INNER JOIN `permission` p ON p.id_permission = rp.permission_id
WHERE p.name_permission IN ('access_admin', 'access_artisan', 'access_salarie');

DELETE FROM `permission` WHERE name_permission IN ('access_admin', 'access_artisan', 'access_salarie');
