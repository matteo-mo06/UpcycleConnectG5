ALTER TABLE `announcement`
    DROP COLUMN `rejection_reason`;

ALTER TABLE `event`
    DROP COLUMN `status`,
    DROP COLUMN `rejection_reason`;
