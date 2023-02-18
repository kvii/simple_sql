CREATE TABLE
  `test_table_name` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `field_one` int NOT NULL COMMENT '字段一',
    `field_two` varchar(45) NOT NULL COMMENT '字段二',
    `created_at` datetime DEFAULT NULL COMMENT '创建于',
    PRIMARY KEY (`id`),
    UNIQUE KEY `id_UNIQUE` (`id`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '测试用数据库';