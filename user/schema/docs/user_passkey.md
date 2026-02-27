# User Passkey

User Passkey.

| Name             | Type      | Key | Comment          |
|------------------|-----------|-----|------------------|
| id               | text      | PRI | ID               |
| user_id          | text      | MUL | User ID          |
| credential_id    | blob      | MUL | Credential ID    |
| public_key       | blob      |     | Public Key       |
| attestation_type | text      |     | Attestation Type |
| transport        | text      |     | Transport        |
| sign_count       | integer   |     | Sign Count       |
| created_at       | timestamp |     | Created At       |
