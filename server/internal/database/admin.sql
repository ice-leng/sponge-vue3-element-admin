/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80036 (8.0.36)
 Source Host           : localhost:3306
 Source Schema         : hyperf

 Target Server Type    : MySQL
 Target Server Version : 80036 (8.0.36)
 File Encoding         : 65001

 Date: 10/11/2024 00:54:23
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for t_config
-- ----------------------------
DROP TABLE IF EXISTS `t_config`;
CREATE TABLE `t_config` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置名称',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '描述',
  `key` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置键',
  `value` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置值',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置';

-- ----------------------------
-- Records of t_config
-- ----------------------------
BEGIN;
INSERT INTO `t_config` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `description`, `key`, `value`) VALUES (2, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, '图片域名', '图片域名', 'imageDomain', 'http://127.0.0.1:9501');
COMMIT;

-- ----------------------------
-- Table structure for t_menu
-- ----------------------------
DROP TABLE IF EXISTS `t_menu`;
CREATE TABLE `t_menu` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `parent_id` int NOT NULL DEFAULT '0' COMMENT '父级',
  `name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单名称',
  `type` enum('CATALOG','MENU','BUTTON','EXTLINK') COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '菜单类型(CATALOG-菜单；MENU-目录；BUTTON-按钮；EXTLINK-外链)',
  `path` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '路由路径',
  `component` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '组件路径(vue页面完整路径，省略.vue后缀)',
  `perm` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '权限标识',
  `sort` int NOT NULL DEFAULT '1' COMMENT '排序',
  `visible` tinyint NOT NULL COMMENT '显示状态',
  `icon` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单图标',
  `redirect` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '跳转路径',
  `always_show` tinyint NOT NULL DEFAULT '0' COMMENT '始终显示',
  `keep_alive` tinyint NOT NULL DEFAULT '1' COMMENT '始终显示',
  `params` json DEFAULT NULL COMMENT '路由参数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单管理';

-- ----------------------------
-- Records of t_menu
-- ----------------------------
BEGIN;
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (1, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 0, '系统管理', 'CATALOG', '/system', 'Layout', '', 1, 1, 'system', 'platform', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (2, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, '管理员管理', 'MENU', 'system/platform', 'system/platform/index', '', 1, 1, 'el-icon-User', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (3, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 2, '管理员新增', 'BUTTON', '', '', 'sys:platform:add', 1, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (4, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 2, '管理员编辑', 'BUTTON', '', '', 'sys:platform:edit', 2, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (5, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 2, '管理员删除', 'BUTTON', '', '', 'sys:platform:delete', 3, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (6, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, '角色管理', 'MENU', 'system/role', 'system/role/index', '', 2, 1, 'role', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (7, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 6, '角色新增', 'BUTTON', '', '', 'sys:role:add', 1, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (8, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 6, '角色编辑', 'BUTTON', '', '', 'sys:role:edit', 2, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (9, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 6, '角色删除', 'BUTTON', '', '', 'sys:role:delete', 3, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (10, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 6, '分配权限', 'BUTTON', '', '', 'sys:role:permission', 4, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (11, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, '菜单管理', 'MENU', 'system/menu', 'system/menu/index', '', 3, 1, 'menu', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (12, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 11, '菜单新增', 'BUTTON', '', '', 'sys:menu:add', 1, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (13, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 11, '菜单编辑', 'BUTTON', '', '', 'sys:menu:edit', 2, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (14, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 11, '菜单删除', 'BUTTON', '', '', 'sys:menu:delete', 3, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (15, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, '系统配置', 'MENU', 'system/config', 'system/config/index', '', 4, 1, 'setting', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (16, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 15, '配置新增', 'BUTTON', '', '', 'sys:config:add', 1, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (17, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 15, '配置编辑', 'BUTTON', '', '', 'sys:config:edit', 2, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (18, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 15, '配置删除', 'BUTTON', '', '', 'sys:config:delete', 3, 1, '', '', 0, 1, NULL);
INSERT INTO `t_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `parent_id`, `name`, `type`, `path`, `component`, `perm`, `sort`, `visible`, `icon`, `redirect`, `always_show`, `keep_alive`, `params`) VALUES (19, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 2, '重置密码', 'BUTTON', '', '', 'sys:platform:password:reset', 4, 1, '', '', 0, 1, NULL);
COMMIT;

-- ----------------------------
-- Table structure for t_migrations
-- ----------------------------
DROP TABLE IF EXISTS `t_migrations`;
CREATE TABLE `t_migrations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `migration` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `batch` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of t_migrations
-- ----------------------------
BEGIN;
INSERT INTO `t_migrations` (`id`, `migration`, `batch`) VALUES (1, '2023_12_16_195120_platform', 1);
INSERT INTO `t_migrations` (`id`, `migration`, `batch`) VALUES (2, '2023_12_31_221138_auth', 1);
COMMIT;

-- ----------------------------
-- Table structure for t_platform
-- ----------------------------
DROP TABLE IF EXISTS `t_platform`;
CREATE TABLE `t_platform` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `username` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账号',
  `password` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `avatar` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'https://foruda.gitee.com/images/1723603502796844527/03cdca2a_716974.gif' COMMENT '头像',
  `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '昵称',
  `mobile` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `gender` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '性别1男2女3保密',
  `role_id` json NOT NULL COMMENT '角色',
  `status` tinyint NOT NULL COMMENT '状态',
  `last_time` datetime DEFAULT NULL COMMENT '上次登录时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员';

-- ----------------------------
-- Records of t_platform
-- ----------------------------
BEGIN;
INSERT INTO `t_platform` (`id`, `created_at`, `updated_at`, `deleted_at`, `username`, `password`, `avatar`, `role_id`, `status`, `last_time`, `mobile`, `nickname`, `gender`) VALUES (1, '2024-11-09 23:56:51', '2024-11-10 00:51:47', NULL, 'admin', '$2y$12$h4UkAJlNkiAuDZguHWEYreIKlv1rnA49QO4uLEipw5TC3KGADHw.W', 'https://foruda.gitee.com/images/1723603502796844527/03cdca2a_716974.gif', '[1]', 1, '2024-11-10 00:51:47', '', '超级管理员', 1);
COMMIT;

-- ----------------------------
-- Table structure for t_role
-- ----------------------------
DROP TABLE IF EXISTS `t_role`;
CREATE TABLE `t_role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `name` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `code` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色编码',
  `sort` int NOT NULL DEFAULT '1' COMMENT '排序',
  `status` tinyint NOT NULL COMMENT '状态',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色管理';

-- ----------------------------
-- Records of t_role
-- ----------------------------
BEGIN;
INSERT INTO `t_role` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `code`, `sort`, `status`) VALUES (1, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, '系统管理员', 'ADMIN', 1, 1);
COMMIT;

-- ----------------------------
-- Table structure for t_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `t_role_menu`;
CREATE TABLE `t_role_menu` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `role_id` int NOT NULL DEFAULT '0' COMMENT '角色ID',
  `menu_id` int NOT NULL DEFAULT '0' COMMENT '菜单ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联';

-- ----------------------------
-- Records of t_role_menu
-- ----------------------------
BEGIN;
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (2, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 2);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (3, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 3);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (4, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 4);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (5, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 5);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (6, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 6);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (7, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 7);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (8, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 8);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (9, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 9);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (10, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 10);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (11, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 11);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (12, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 12);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (13, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 13);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (14, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 14);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (15, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 15);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (16, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 16);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (17, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 17);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (18, '2024-11-09 23:56:51', '2024-11-09 23:56:51', NULL, 1, 18);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (19, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 1);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (20, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 2);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (21, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 3);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (22, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 4);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (23, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 5);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (24, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 19);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (25, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 6);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (26, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 7);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (27, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 8);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (28, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 9);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (29, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 10);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (30, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 11);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (31, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 12);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (32, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 13);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (33, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 14);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (34, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 15);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (35, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 16);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (36, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 17);
INSERT INTO `t_role_menu` (`id`, `created_at`, `updated_at`, `deleted_at`, `role_id`, `menu_id`) VALUES (37, '2024-11-10 00:51:27', '2024-11-10 00:51:27', NULL, 1, 18);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
