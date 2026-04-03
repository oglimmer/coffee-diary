ALTER TABLE users ADD COLUMN oidc_sub VARCHAR(255) AFTER username;
ALTER TABLE users ADD UNIQUE KEY uk_users_oidc_sub (oidc_sub);
ALTER TABLE users MODIFY COLUMN password VARCHAR(255) NULL;
