CREATE TABLE `order` (
     `order_id` BINARY(16),
     `cost` INTEGER NOT NULL,
     `created_at` DATETIME NOT NULL,
     `updated_at` DATETIME NOT NULL,
     `deleted_at` DATETIME,
     PRIMARY KEY (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `order_item` (
    `order_id` BINARY(16),
    `menu_item_id` BINARY(16),
    `quantity` INTEGER NOT NULL,
    PRIMARY KEY (order_id, menu_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci