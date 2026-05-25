ALTER TABLE `announcement`
    ADD COLUMN `rejection_reason` text DEFAULT NULL;

ALTER TABLE `event`
    ADD COLUMN `status` varchar(20) NOT NULL DEFAULT 'approved',
    ADD COLUMN `rejection_reason` text DEFAULT NULL;
