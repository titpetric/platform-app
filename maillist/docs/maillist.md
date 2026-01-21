# Maillist

Table containing mailing list definitions.

| Name        | Type         | Key | Comment                                          |
|-------------|--------------|-----|--------------------------------------------------|
| id          | char(26)     | PRI | ULID identifier for the mailing list             |
| name        | varchar(255) |     | Human-readable name of the mailing list          |
| slug        | varchar(100) |     | URL-friendly identifier for subscription links   |
| description | text         |     | Optional description of the mailing list         |
| created_at  | timestamp    |     | Timestamp when the mailing list was created      |
| updated_at  | timestamp    |     | Timestamp when the mailing list was last updated |
