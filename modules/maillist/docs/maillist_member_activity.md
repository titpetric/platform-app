# Maillist Member Activity

Table recording activity of mailing list members.

| Name               | Type                                        | Key | Comment                                       |
|--------------------|---------------------------------------------|-----|-----------------------------------------------|
| id                 | char(26)                                    | PRI | ULID identifier for the activity record       |
| maillist_member_id | char(26)                                    |     | ULID of the member performing the activity    |
| campaign_id        | char(26)                                    |     | ULID of the campaign related to this activity |
| activity_type      | enum('open','click','bounce','unsubscribe') |     | Type of activity                              |
| activity_at        | timestamp                                   |     | Timestamp of the activity                     |
