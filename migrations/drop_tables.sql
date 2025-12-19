-- Drop tables in reverse order to avoid foreign key constraints
-- WARNING: This will delete all data!

SET FOREIGN_KEY_CHECKS = 0;

-- Drop foreign key constraints first
ALTER TABLE `products` DROP FOREIGN KEY IF EXISTS `fk_products_category`;
ALTER TABLE `categories` DROP FOREIGN KEY IF EXISTS `fk_products_category`;

-- Drop tables
DROP TABLE IF EXISTS `products`;
DROP TABLE IF EXISTS `categories`;
DROP TABLE IF EXISTS `users`;

SET FOREIGN_KEY_CHECKS = 1;

