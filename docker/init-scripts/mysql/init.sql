-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS ky_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE ky_admin;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `password` varchar(100) NOT NULL COMMENT '密码哈希',
    `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
    `mobile` varchar(20) DEFAULT NULL COMMENT '手机号',
    `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
    `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    KEY `idx_email` (`email`),
    KEY `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 创建角色表
CREATE TABLE IF NOT EXISTS `roles` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(50) NOT NULL COMMENT '角色名称',
    `code` varchar(50) NOT NULL COMMENT '角色代码',
    `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 创建权限表
CREATE TABLE IF NOT EXISTS `permissions` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(50) NOT NULL COMMENT '权限名称',
    `code` varchar(50) NOT NULL COMMENT '权限代码',
    `description` varchar(255) DEFAULT NULL COMMENT '权限描述',
    `parent_id` bigint(20) unsigned DEFAULT NULL COMMENT '父权限ID',
    `type` varchar(20) NOT NULL COMMENT '权限类型: menu-菜单, button-按钮, api-接口',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 0-禁用, 1-启用',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`),
    KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 创建用户角色关联表
CREATE TABLE IF NOT EXISTS `user_roles` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_role` (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 创建角色权限关联表
CREATE TABLE IF NOT EXISTS `role_permissions` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
    `permission_id` bigint(20) unsigned NOT NULL COMMENT '权限ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_role_permission` (`role_id`,`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 创建操作日志表
CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '用户ID',
    `ip` varchar(50) DEFAULT NULL COMMENT 'IP地址',
    `method` varchar(10) DEFAULT NULL COMMENT '请求方法',
    `path` varchar(255) DEFAULT NULL COMMENT '请求路径',
    `status` int(11) DEFAULT NULL COMMENT '响应状态码',
    `latency` bigint(20) DEFAULT NULL COMMENT '响应时间(ms)',
    `user_agent` varchar(255) DEFAULT NULL COMMENT '用户代理',
    `request` text COMMENT '请求内容',
    `response` text COMMENT '响应内容',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- 创建系统配置表
CREATE TABLE IF NOT EXISTS `system_configs` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `key` varchar(50) NOT NULL COMMENT '配置键',
    `value` text NOT NULL COMMENT '配置值',
    `description` varchar(255) DEFAULT NULL COMMENT '配置描述',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 插入默认管理员用户（密码为 admin123）
INSERT INTO `users` (`username`, `password`, `email`, `nickname`, `status`) VALUES
('admin', '$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2', 'admin@kysion.com', '管理员', 1);

-- 插入默认角色
INSERT INTO `roles` (`name`, `code`, `description`, `status`) VALUES
('管理员', 'admin', '系统管理员，拥有所有权限', 1),
('普通用户', 'user', '普通用户，拥有基本权限', 1);

-- 将管理员用户关联到管理员角色
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES
(1, 1);

-- 插入基本权限
INSERT INTO `permissions` (`name`, `code`, `description`, `parent_id`, `type`, `status`) VALUES
('系统管理', 'system:manage', '系统管理权限', NULL, 'menu', 1),
('用户管理', 'user:manage', '用户管理权限', 1, 'menu', 1),
('用户查询', 'user:read', '查询用户信息', 2, 'api', 1),
('用户创建', 'user:create', '创建新用户', 2, 'api', 1),
('用户编辑', 'user:update', '编辑用户信息', 2, 'api', 1),
('用户删除', 'user:delete', '删除用户', 2, 'api', 1),
('角色管理', 'role:manage', '角色管理权限', 1, 'menu', 1),
('权限管理', 'permission:manage', '权限管理权限', 1, 'menu', 1);

-- 将权限关联到管理员角色
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8);

-- 将基本查询权限关联到普通用户角色
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(2, 3); 