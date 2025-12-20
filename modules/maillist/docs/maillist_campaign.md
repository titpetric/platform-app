# Maillist Campaign

Table containing email campaigns for mailing lists.

| Name         | Type                             | Key | Comment                                                      |
|--------------|----------------------------------|-----|--------------------------------------------------------------|
| id           | char(26)                         | PRI | ULID identifier for the campaign                             |
| maillist_id  | char(26)                         |     | ULID of the mailing list this campaign belongs to            |
| subject      | varchar(255)                     |     | Email subject line                                           |
| body         | text                             |     | Content of the email message                                 |
| status       | enum('draft','scheduled','sent') |     | Current status of the campaign                               |
| scheduled_at | timestamp                        |     | Optional timestamp when the campaign is scheduled to be sent |
| sent_at      | timestamp                        |     | Timestamp when the campaign was actually sent                |
| created_at   | timestamp                        |     | Timestamp when the campaign was created                      |
| updated_at   | timestamp                        |     | Timestamp when the campaign was last updated                 |
