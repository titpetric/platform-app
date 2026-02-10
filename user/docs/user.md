# User

Stores user profile information using ULID as primary key.

| Name       | Type         | Key | Comment                               |
|------------|--------------|-----|---------------------------------------|
| id         | char(26)     | PRI | Primary key: ULID string              |
| full_name  | varchar(100) |     | User full name                        |
| username   | varchar(100) |     | User username                         |
| slug       | varchar(100) |     | URL-friendly username slug            |
| deleted_at | timestamp    | MUL | Soft delete timestamp, NULL if active |
| created_at | timestamp    |     | Record creation timestamp             |
| updated_at | timestamp    |     | Record update timestamp               |
