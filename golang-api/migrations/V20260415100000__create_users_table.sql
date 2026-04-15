-- =============================================================================
-- Migration: create users table
-- Version:   V20260415100000
-- Date:      2026-04-15 10:00:00
-- Author:    joaquingomezaj@gmail.com
-- =============================================================================

CREATE TABLE users (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at  DATETIME(3) NULL,
    updated_at  DATETIME(3) NULL,
    deleted_at  DATETIME(3) NULL,
    first_name  VARCHAR(191) NOT NULL,
    last_name   VARCHAR(191) NOT NULL,
    email       VARCHAR(191) NOT NULL,
    password    VARCHAR(191) NOT NULL,
    role        VARCHAR(50)  NOT NULL DEFAULT 'user',
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    UNIQUE KEY uk_users_email (email),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
