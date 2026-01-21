-- Table to store mailing lists
CREATE TABLE `maillist` (
    `id` CHAR(26) NOT NULL COMMENT 'ULID identifier for the mailing list',
    `name` VARCHAR(255) NOT NULL COMMENT 'Human-readable name of the mailing list',
    `slug` VARCHAR(100) NOT NULL COMMENT 'URL-friendly identifier for subscription links',
    `description` TEXT COMMENT 'Optional description of the mailing list',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp when the mailing list was created',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Timestamp when the mailing list was last updated',
    PRIMARY KEY (`id`)
) COMMENT='Table containing mailing list definitions';

-- Table to store members/subscribers of mailing lists
CREATE TABLE `maillist_member` (
    `id` CHAR(26) NOT NULL COMMENT 'ULID identifier for the member record',
    `maillist_id` CHAR(26) NOT NULL COMMENT 'ULID of the mailing list this member is subscribed to',
    `email` VARCHAR(255) NOT NULL COMMENT 'Email address of the subscriber',
    `name` VARCHAR(255) COMMENT 'Optional name of the subscriber',
    `subscribed_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp when the subscriber joined the list',
    `unsubscribed_at` TIMESTAMP NULL COMMENT 'Timestamp when the subscriber unsubscribed, if applicable',
    `source` VARCHAR(100) COMMENT 'Optional source of subscription (HTML snippet, API, etc.)',
    PRIMARY KEY (`id`)
) COMMENT='Table containing subscribers for mailing lists';

-- Table to store mailing list campaigns / messages
CREATE TABLE `maillist_campaign` (
    `id` CHAR(26) NOT NULL COMMENT 'ULID identifier for the campaign',
    `maillist_id` CHAR(26) NOT NULL COMMENT 'ULID of the mailing list this campaign belongs to',
    `subject` VARCHAR(255) NOT NULL COMMENT 'Email subject line',
    `body` TEXT NOT NULL COMMENT 'Content of the email message',
    `status` ENUM('draft','scheduled','sent') NOT NULL DEFAULT 'draft' COMMENT 'Current status of the campaign',
    `scheduled_at` TIMESTAMP NULL COMMENT 'Optional timestamp when the campaign is scheduled to be sent',
    `sent_at` TIMESTAMP NULL COMMENT 'Timestamp when the campaign was actually sent',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp when the campaign was created',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Timestamp when the campaign was last updated',
    PRIMARY KEY (`id`)
) COMMENT='Table containing email campaigns for mailing lists';

-- Table to store activity logs for members (opens, clicks, etc.)
CREATE TABLE `maillist_member_activity` (
    `id` CHAR(26) NOT NULL COMMENT 'ULID identifier for the activity record',
    `maillist_member_id` CHAR(26) NOT NULL COMMENT 'ULID of the member performing the activity',
    `campaign_id` CHAR(26) NOT NULL COMMENT 'ULID of the campaign related to this activity',
    `activity_type` ENUM('open','click','bounce','unsubscribe') NOT NULL COMMENT 'Type of activity',
    `activity_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Timestamp of the activity',
    PRIMARY KEY (`id`)
) COMMENT='Table recording activity of mailing list members';
