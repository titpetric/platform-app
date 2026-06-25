# User Passkey

User Passkey.

| Name             | Type     | Key | Comment          |
|------------------|----------|-----|------------------|
| id               | varchar  | PRI | ID               |
| user_id          | varchar  | MUL | User ID          |
| credential_id    | blob     | MUL | Credential ID    |
| public_key       | blob     |     | Public Key       |
| attestation_type | varchar  |     | Attestation Type |
| transport        | varchar  |     | Transport        |
| sign_count       | bigint   |     | Sign Count       |
| created_at       | datetime |     | Created At       |
