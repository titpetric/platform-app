# User Auth

User Auth.

| Name               | Type     | Key | Comment            |
|--------------------|----------|-----|--------------------|
| user_id            | varchar  | PRI | User ID            |
| email              | varchar  | MUL | Email              |
| password           | varchar  |     | Password           |
| created_at         | datetime |     | Created At         |
| updated_at         | datetime |     | Updated At         |
| activated_at       | datetime |     | Activated At       |
| activation_token   | varchar  | MUL | Activation Token   |
| activation_sent_at | datetime |     | Activation Sent At |
