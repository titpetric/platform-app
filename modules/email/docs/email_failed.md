# Email Failed

| Name        | Type     | Key | Comment     |
|-------------|----------|-----|-------------|
| id          | TEXT     | PRI | ID          |
| recipient   | TEXT     |     | Recipient   |
| subject     | TEXT     |     | Subject     |
| body        | TEXT     |     | Body        |
| status      | TEXT     |     | Status      |
| created_at  | DATETIME |     | Created At  |
| sent_at     | DATETIME |     | Sent At     |
| error       | TEXT     |     | Error       |
| retry_count | INTEGER  |     | Retry Count |
| last_error  | TEXT     |     | Last Error  |
| last_retry  | DATETIME |     | Last Retry  |
