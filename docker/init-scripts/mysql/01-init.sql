-- 初始化数据库
USE ky_admin;

-- 创建初始角色表
CREATE TABLE IF NOT EXISTS `roles` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_roles_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 创建初始权限表
CREATE TABLE IF NOT EXISTS `permissions` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '权限名称',
  `description` varchar(255) DEFAULT NULL COMMENT '权限描述',
  `module` varchar(50) NOT NULL COMMENT '模块名称',
  `action` varchar(50) NOT NULL COMMENT '操作类型',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_permissions_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 创建角色权限关联表
CREATE TABLE IF NOT EXISTS `role_permissions` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int(11) unsigned NOT NULL COMMENT '角色ID',
  `permission_id` int(11) unsigned NOT NULL COMMENT '权限ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_permission` (`role_id`,`permission_id`),
  KEY `idx_permission_id` (`permission_id`),
  CONSTRAINT `fk_rp_permission_id` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_rp_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码哈希',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) DEFAULT NULL COMMENT '电子邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号码',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
  `last_login` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_username` (`username`),
  UNIQUE KEY `idx_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 创建用户角色关联表
CREATE TABLE IF NOT EXISTS `user_roles` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned NOT NULL COMMENT '用户ID',
  `role_id` int(11) unsigned NOT NULL COMMENT '角色ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_role` (`user_id`,`role_id`),
  KEY `idx_ur_role_id` (`role_id`),
  CONSTRAINT `fk_ur_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_ur_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 插入初始角色
INSERT INTO `roles` (`name`, `description`, `status`) VALUES
('admin', '系统管理员', 1),
('user', '普通用户', 1);

-- 插入初始权限
INSERT INTO `permissions` (`name`, `description`, `module`, `action`) VALUES
('user:view', '查看用户信息', 'user', 'view'),
('user:create', '创建用户', 'user', 'create'),
('user:update', '更新用户信息', 'user', 'update'),
('user:delete', '删除用户', 'user', 'delete'),
('role:view', '查看角色信息', 'role', 'view'),
('role:create', '创建角色', 'role', 'create'),
('role:update', '更新角色信息', 'role', 'update'),
('role:delete', '删除角色', 'role', 'delete');

-- 给admin角色分配所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 1, id FROM `permissions`;

-- 给user角色分配部分权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 2, id FROM `permissions` WHERE `name` IN ('user:view', 'role:view');

-- 创建管理员用户 (密码: admin123)
INSERT INTO `users` (`username`, `password`, `nickname`, `email`, `status`)
VALUES ('admin', '$2a$10$N/EGBt7HJcfM4VpCBChMI.Jc.G4GDjxNC09hNFGwuFUmYfN9OZiT.', '管理员', 'admin@kysion.com', 1);

-- 创建测试用户 (密码: user123)
INSERT INTO `users` (`username`, `password`, `nickname`, `email`, `status`)
VALUES ('user', '$2a$10$RTpWzHrGe7lRRzZ7IWgG8.Y1sSNECJ25XuZ/JHcpfQO/jER7VsING', '测试用户', 'user@kysion.com', 1);

-- 给用户分配角色
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES (1, 1); -- admin -> admin
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES (2, 2); -- user -> user 