# Maillist Member

Table containing subscribers for mailing lists.

| Name            | Type         | Key | Comment                                                   |
|-----------------|--------------|-----|-----------------------------------------------------------|
| id              | char(26)     | PRI | ULID identifier for the member record                     |
| maillist_id     | char(26)     |     | ULID of the mailing list this member is subscribed to     |
| email           | varchar(255) |     | Email address of the subscriber                           |
| name            | varchar(255) |     | Optional name of the subscriber                           |
| subscribed_at   | timestamp    |     | Timestamp when the subscriber joined the list             |
| unsubscribed_at | timestamp    |     | Timestamp when the subscriber unsubscribed, if applicable |
| source          | varchar(100) |     | Optional source of subscription (HTML snippet, API, etc.) |
