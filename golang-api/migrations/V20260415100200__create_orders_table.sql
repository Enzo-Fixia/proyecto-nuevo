-- =============================================================================
-- Migration: create orders table
-- Version:   V20260415100200
-- Date:      2026-04-15 10:02:00
-- Author:    joaquingomezaj@gmail.com
-- =============================================================================

CREATE TABLE orders (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at  DATETIME(3) NULL,
    updated_at  DATETIME(3) NULL,
    deleted_at  DATETIME(3) NULL,
    user_id     BIGINT UNSIGNED NOT NULL,
    total       DECIMAL(12,2)   NOT NULL,
    status      VARCHAR(50)     NOT NULL DEFAULT 'pending',
    INDEX idx_orders_deleted_at (deleted_at),
    INDEX idx_orders_user_id (user_id),
    INDEX idx_orders_status (status),
    CONSTRAINT fk_orders_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
