ALTER TABLE LOCKER DROP INDEX locker_number, DROP COLUMN locker_number;
DELETE FROM LOCKER WHERE locker_number > 5;
