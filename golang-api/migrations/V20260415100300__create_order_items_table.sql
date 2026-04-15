-- =============================================================================
-- Migration: create order_items table
-- Version:   V20260415100300
-- Date:      2026-04-15 10:03:00
-- Author:    joaquingomezaj@gmail.com
-- =============================================================================

CREATE TABLE order_items (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at  DATETIME(3) NULL,
    updated_at  DATETIME(3) NULL,
    deleted_at  DATETIME(3) NULL,
    order_id    BIGINT UNSIGNED NOT NULL,
    product_id  BIGINT UNSIGNED NOT NULL,
    quantity    INT             NOT NULL,
    unit_price  DECIMAL(12,2)   NOT NULL,
    INDEX idx_order_items_deleted_at (deleted_at),
    INDEX idx_order_items_order_id (order_id),
    INDEX idx_order_items_product_id (product_id),
    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id) REFERENCES orders(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_order_items_product
        FOREIGN KEY (product_id) REFERENCES products(id)
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
