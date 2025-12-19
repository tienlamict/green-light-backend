-- Migration: Add Product Variants Support
-- Created: 2025-11-06
-- Description: Support multiple SKUs per product (variants: color, size, etc.)

-- Create product_variants table
CREATE TABLE IF NOT EXISTS `product_variants` (
  `variant_id` VARCHAR(36) NOT NULL,
  `product_id` VARCHAR(36) NOT NULL,
  `sku` VARCHAR(100) NOT NULL,
  `name` VARCHAR(255) NOT NULL COMMENT 'Variant name: e.g., "Red - Large", "Blue - Small"',
  `attributes` JSON COMMENT 'Variant attributes: {"color": "red", "size": "L"}',
  `price` DECIMAL(10,2) COMMENT 'Override product price if different',
  `stock` INT NOT NULL DEFAULT 0,
  `is_active` BOOLEAN NOT NULL DEFAULT TRUE,
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`variant_id`),
  UNIQUE INDEX `idx_variants_sku` (`sku`),
  INDEX `idx_variants_product_id` (`product_id`),
  INDEX `idx_variants_is_active` (`is_active`),
  CONSTRAINT `fk_variants_product`
    FOREIGN KEY (`product_id`)
    REFERENCES `products`(`product_id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Update products table: Make SKU nullable (optional) since variants will have SKUs
ALTER TABLE `products`
MODIFY COLUMN `sku` VARCHAR(100) NULL COMMENT 'Base SKU (optional if using variants)';

-- Add index for variant lookup
CREATE INDEX IF NOT EXISTS `idx_products_use_variants` ON `products`(`product_id`);

