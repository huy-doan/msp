-- master_mfa_types
INSERT INTO `master_mfa_types` (id, no, title, is_active, created_at, updated_at)
VALUES
   (1, 'OTP', 1, NOW(), NOW()),
   (2, 'メール', 1, NOW(), NOW()),
   (3, 'SMS', 1, NOW(), NOW());
