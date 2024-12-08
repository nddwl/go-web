CREATE DATABASE IF NOT EXISTS 'order';
USE `order`;
CREATE TABLE `user`(
                       `id` int unsigned NOT NULL AUTO_INCREMENT,
                       `created_at` datetime DEFAULT NULL ,
                       `updated_at` datetime DEFAULT NULL ,
                       `deleted_at` datetime DEFAULT NULL,
                       `uid` bigint NOT NULL ,
                       `name` varchar(10) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                       `avatar` varchar(255) CHARSET utf8mb4 DEFAULT NULL COLLATE Utf8mb4_General_Ci,
                       `email` varchar(40) CHARSET utf8mb4 DEFAULT NULL COLLATE Utf8mb4_General_Ci,
                       `phone` varchar(11) CHARSET utf8mb4 DEFAULT NULL COLLATE Utf8mb4_General_Ci,
                       `exp` int DEFAULT 0,
                       `coin` int DEFAULT 0,
                       `status` tinyint unsigned DEFAULT 0,
                       `role` tinyint unsigned DEFAULT 0,
                       PRIMARY KEY (`id`),
                       UNIQUE KEY (`uid`) USING BTREE ,
                       UNIQUE KEY (`name`) USING BTREE ,
                       UNIQUE KEY (`email`) USING BTREE
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;
CREATE TABLE `password`(
                           `id` int unsigned NOT NULL AUTO_INCREMENT,
                           `created_at` datetime DEFAULT NULL ,
                           `updated_at` datetime DEFAULT NULL ,
                           `deleted_at` datetime DEFAULT NULL,
                           `uid` bigint NOT NULL ,
                           `pwd_hash` varchar(64) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                           PRIMARY KEY (`id`),
                           UNIQUE (`uid`) USING BTREE
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;

CREATE TABLE `passport`(
                           `id` int unsigned NOT NULL AUTO_INCREMENT,
                           `created_at` datetime DEFAULT NULL ,
                           `updated_at` datetime DEFAULT NULL ,
                           `deleted_at` datetime DEFAULT NULL,
                           `uid` bigint NOT NULL ,
                           `token` varchar(36) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                           `ip` varchar(15) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                           `device_id` varchar(36) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                           `ua` text CHARSET utf8mb4 DEFAULT NULL COLLATE Utf8mb4_General_Ci,
                           PRIMARY KEY (`id`),
                           INDEX (`uid`) USING BTREE ,
                           UNIQUE (`token`) USING BTREE ,
                           UNIQUE (`device_id`) USING BTREE
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;

CREATE TABLE `user_sign`(
                            `id` int unsigned NOT NULL AUTO_INCREMENT,
                            `created_at` datetime DEFAULT NULL ,
                            `updated_at` datetime DEFAULT NULL ,
                            `deleted_at` datetime DEFAULT NULL,
                            `uid` bigint NOT NULL ,
                            `status` tinyint unsigned DEFAULT 0,
                            `reward` varchar(255) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                            PRIMARY KEY (`id`),
                            INDEX (`uid`) USING BTREE
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;

CREATE TABLE `activity`(
                           `id` int unsigned NOT NULL AUTO_INCREMENT,
                           `created_at` datetime DEFAULT NULL ,
                           `updated_at` datetime DEFAULT NULL ,
                           `deleted_at` datetime DEFAULT NULL,
                           `uuid`  bigint NOT NULL ,
                           `name` varchar(50) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                           `url` varchar(255) CHARSET utf8mb4 DEFAULT NULL COLLATE Utf8mb4_General_Ci,
                           `type` tinyint unsigned DEFAULT 0,
                           `status` tinyint  unsigned DEFAULT 0,
                           `cost` int unsigned DEFAULT 0,
                           `info` varchar(255) CHARSET utf8mb4 DEFAULT NULL COLLATE Utf8mb4_General_Ci,
                           PRIMARY KEY (`id`),
                           UNIQUE (`uuid`) USING BTREE
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;

CREATE TABLE `prize`(
                        `id` int unsigned NOT NULL AUTO_INCREMENT,
                        `created_at` datetime DEFAULT NULL ,
                        `updated_at` datetime DEFAULT NULL ,
                        `deleted_at` datetime DEFAULT NULL,
                        `activity_uuid` bigint NOT NULL ,
                        `uuid`  bigint NOT NULL ,
                        `name` varchar(10) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                        `type` tinyint unsigned DEFAULT 0,
                        `value` varchar(20) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                        `initial_stock` int unsigned DEFAULT 0,
                        `stock` int  unsigned DEFAULT 0,
                        `score` int unsigned NOT NULL ,
                        PRIMARY KEY (`id`),
                        INDEX (`activity_uuid`) USING BTREE,
                        UNIQUE KEY (`uuid`) USING BTREE
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;

CREATE TABLE `activity_record`(
                                  `id` int unsigned NOT NULL AUTO_INCREMENT,
                                  `created_at` datetime DEFAULT NULL ,
                                  `updated_at` datetime DEFAULT NULL ,
                                  `deleted_at` datetime DEFAULT NULL,
                                  `uid` bigint NOT NULL ,
                                  `activity_uuid`  bigint NOT NULL ,
                                  `prize_uuid`  bigint NOT NULL ,
                                  `remark` varchar(255)CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                                  PRIMARY KEY (`id`),
                                  INDEX (`uid`) USING BTREE
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;

CREATE TABLE `post`(
                       `id` int unsigned NOT NULL AUTO_INCREMENT,
                       `created_at` datetime DEFAULT NULL ,
                       `updated_at` datetime DEFAULT NULL ,
                       `deleted_at` datetime DEFAULT NULL,
                       `uuid` bigint NOT NULL ,
                       `title` varchar(50) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                       `uid` bigint NOT NULL ,
                       `category` tinyint unsigned NOT NULL ,
                       `top_category` tinyint unsigned NOT NULL ,
                       `summary` varchar(255) CHARSET utf8mb4 DEFAULT NULL COLLATE Utf8mb4_General_Ci,
                       `content` text CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                       `source` tinyint unsigned DEFAULT 0,
                       PRIMARY KEY (`id`),
                       UNIQUE KEY (`uuid`) USING BTREE ,
                       INDEX (`uid`) USING BTREE ,
                       FULLTEXT (`title`)
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;

CREATE TABLE `post_tag`(
                           `id` int unsigned NOT NULL AUTO_INCREMENT,
                           `created_at` datetime DEFAULT NULL ,
                           `updated_at` datetime DEFAULT NULL ,
                           `deleted_at` datetime DEFAULT NULL,
                           `uuid` bigint NOT NULL ,
                           `name` varchar(10) CHARSET utf8mb4 NOT NULL COLLATE Utf8mb4_General_Ci,
                           PRIMARY KEY (`id`),
                           UNIQUE KEY (`uuid`) USING BTREE,
                           UNIQUE KEY (`name`) USING BTREE
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;

CREATE TABLE `post_association`(
                                   `id` int unsigned NOT NULL AUTO_INCREMENT,
                                   `created_at` datetime DEFAULT NULL ,
                                   `updated_at` datetime DEFAULT NULL ,
                                   `deleted_at` datetime DEFAULT NULL,
                                   `post_tag_uuid` bigint NOT NULL ,
                                   `post_uuid` bigint NOT NULL ,
                                   PRIMARY KEY (`id`),
                                   INDEX (`post_tag_uuid`) USING BTREE ,
                                   INDEX (`post_uuid`) USING BTREE
)ENGINE InnoDB DEFAULT CHARSET utf8mb4 COLLATE Utf8mb4_General_Ci ;