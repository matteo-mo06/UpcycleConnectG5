ALTER TABLE LOCKER ADD COLUMN locker_number INT NULL;

UPDATE LOCKER SET locker_number = 1 WHERE id_locker = 'c46811ae-8fda-4b2c-9757-fc4d4e56ef08';
UPDATE LOCKER SET locker_number = 2 WHERE id_locker = '90c96690-e4ff-4711-b900-89be2fe9dcdb';
UPDATE LOCKER SET locker_number = 3 WHERE id_locker = '4fade1b3-d6b6-4635-b505-0820f31b41f2';
UPDATE LOCKER SET locker_number = 4 WHERE id_locker = 'fbe9708e-3919-4b15-b59f-767821715847';
UPDATE LOCKER SET locker_number = 5 WHERE id_locker = '59adc1b0-5760-425f-81b7-3e61755cf3b0';

INSERT INTO LOCKER (id_locker, locker_number) VALUES
(UUID(), 6), (UUID(), 7), (UUID(), 8), (UUID(), 9), (UUID(), 10),
(UUID(), 11), (UUID(), 12), (UUID(), 13), (UUID(), 14), (UUID(), 15),
(UUID(), 16), (UUID(), 17), (UUID(), 18), (UUID(), 19), (UUID(), 20);

ALTER TABLE LOCKER ADD UNIQUE locker_number_unique (locker_number);
