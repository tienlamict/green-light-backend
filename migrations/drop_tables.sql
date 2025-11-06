-- Drop tables in reverse order to avoid foreign key constraints
-- WARNING: This will delete all data!

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `products`;
DROP TABLE IF EXISTS `categories`;
DROP TABLE IF EXISTS `users`;

SET FOREIGN_KEY_CHECKS = 1;

