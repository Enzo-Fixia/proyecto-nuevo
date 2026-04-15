-- =============================================================================
-- Migration: create products table
-- Version:   V20260415100100
-- Date:      2026-04-15 10:01:00
-- Author:    joaquingomezaj@gmail.com
-- =============================================================================

CREATE TABLE products (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at  DATETIME(3) NULL,
    updated_at  DATETIME(3) NULL,
    deleted_at  DATETIME(3) NULL,
    name        VARCHAR(191) NOT NULL,
    description TEXT,
    price       DECIMAL(12,2) NOT NULL,
    stock       INT NOT NULL DEFAULT 0,
    category    VARCHAR(191),
    image_url   VARCHAR(500),
    user_id     BIGINT UNSIGNED,
    INDEX idx_products_deleted_at (deleted_at),
    INDEX idx_products_user_id (user_id),
    CONSTRAINT fk_products_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
