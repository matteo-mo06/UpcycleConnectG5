ALTER TABLE USER ADD COLUMN announcement_count INT NOT NULL DEFAULT 0;

UPDATE USER u SET announcement_count = (
    SELECT COUNT(*) FROM USER_ANNOUNCEMENT WHERE id_user = u.id_user
);
