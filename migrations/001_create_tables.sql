-- Migration: Create Tables
-- Created: 2025-11-06
-- Description: Initial database schema for Green Light Backend

-- Create users table first (no foreign keys)
CREATE TABLE IF NOT EXISTS `users` (
  `user_id` VARCHAR(36) NOT NULL,
  `email` VARCHAR(255) NOT NULL UNIQUE,
  `password_hash` VARCHAR(255) NOT NULL,
  `role` ENUM('admin', 'editor') NOT NULL DEFAULT 'editor',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`user_id`),
  INDEX `idx_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create categories table (no foreign keys)
CREATE TABLE IF NOT EXISTS `categories` (
  `category_id` VARCHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `slug` VARCHAR(255) NOT NULL UNIQUE,
  `description` TEXT,
  `is_active` BOOLEAN NOT NULL DEFAULT TRUE,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`category_id`),
  UNIQUE INDEX `idx_categories_slug` (`slug`),
  INDEX `idx_categories_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create products table last (has foreign key to categories)
CREATE TABLE IF NOT EXISTS `products` (
  `product_id` VARCHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `slug` VARCHAR(255) NOT NULL,
  `sku` VARCHAR(100) NOT NULL,
  `short_desc` VARCHAR(500),
  `description` TEXT,
  `price` DECIMAL(10,2) NOT NULL,
  `stock` INT NOT NULL DEFAULT 0,
  `thumbnail_url` VARCHAR(500),
  `gallery` JSON,
  `category_id` VARCHAR(36) NOT NULL,
  `is_active` BOOLEAN NOT NULL DEFAULT TRUE,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`product_id`),
  UNIQUE INDEX `idx_products_slug` (`slug`),
  INDEX `idx_products_sku` (`sku`),
  INDEX `idx_products_category_id` (`category_id`),
  INDEX `idx_products_is_active` (`is_active`),
  INDEX `idx_products_price` (`price`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add foreign key constraint separately to ensure correct relationship
-- Products reference Categories (not the other way around!)
ALTER TABLE `products`
ADD CONSTRAINT `fk_products_category` 
  FOREIGN KEY (`category_id`) 
  REFERENCES `categories`(`category_id`) 
  ON DELETE CASCADE 
  ON UPDATE CASCADE;

-- Insert default admin user
-- Email: admin@example.com
-- Password: admin123
INSERT INTO `users` (`user_id`, `email`, `password_hash`, `role`, `created_at`, `updated_at`)
VALUES (
  'a628f74e-6a75-4bc8-a287-51264a87b6fd',
  'admin@example.com',
  '$2a$10$2.Od0DDLmSf3LGRn21/wmu2owU7lnvYFoI4tSCzXNeGYNl7GE88F6',
  'admin',
  NOW(3),
  NOW(3)
)
ON DUPLICATE KEY UPDATE `email` = `email`;

