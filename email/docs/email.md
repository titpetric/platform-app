# Email

| Name        | Type     | Key | Comment     |
|-------------|----------|-----|-------------|
| id          | TEXT     | PRI | ID          |
| recipient   | TEXT     |     | Recipient   |
| subject     | TEXT     |     | Subject     |
| body        | TEXT     |     | Body        |
| status      | TEXT     |     | Status      |
| error       | TEXT     |     | Error       |
| created_at  | DATETIME |     | Created At  |
| sent_at     | DATETIME |     | Sent At     |
| retry_count | INTEGER  |     | Retry Count |
| retry_error | TEXT     |     | Retry Error |
| retry_at    | DATETIME |     | Retry At    |
