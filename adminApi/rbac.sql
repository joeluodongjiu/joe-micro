/*
 Navicat Premium Data Transfer

 Source Server         : 开发环境
 Source Server Type    : MySQL
 Source Server Version : 50723
 Source Host           : 192.168.0.162:3306
 Source Schema         : rbac

 Target Server Type    : MySQL
 Target Server Version : 50723
 File Encoding         : 65001

 Date: 31/07/2019 19:21:07
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_menu`;
CREATE TABLE `admin_menu` (
  `id` varchar(32) NOT NULL COMMENT '主键',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deletedAt` datetime DEFAULT NULL COMMENT '删除时间',
  `name` varchar(64) NOT NULL COMMENT '菜单名称',
  `sequence` int(11) DEFAULT '0' COMMENT '排序值',
  `url` varchar(128) NOT NULL,
  `parent_id` varchar(32) DEFAULT NULL,
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1.启用   2.禁用',
  `menu_type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1.模块  2.菜单   3.操作',
  `code` varchar(32) NOT NULL COMMENT '菜单代码\n',
  `operate_type` enum('read','write') NOT NULL DEFAULT 'read',
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_menu_id_uindex` (`id`),
  UNIQUE KEY `admin_menu_name_uindex` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';

-- ----------------------------
-- Records of admin_menu
-- ----------------------------
BEGIN;
INSERT INTO `admin_menu` VALUES ('1156076295484538880', '2019-07-30 13:50:24', '2019-07-30 13:50:47', NULL, '系统管理', 1, '', '2147483647', 1, 1, 'Sys', 'read');
INSERT INTO `admin_menu` VALUES ('1156080136950493184', '2019-07-30 13:52:38', '2019-07-30 14:02:14', NULL, '菜单管理', 20, '/menu', '1156076295484538880', 1, 2, 'Menu', 'read');
INSERT INTO `admin_menu` VALUES ('1156080137621581824', '2019-07-30 13:52:38', '2019-07-30 14:02:15', NULL, '角色管理', 30, '/role', '1156076295484538880', 1, 2, 'Role', 'read');
INSERT INTO `admin_menu` VALUES ('1156080138087149568', '2019-07-30 13:52:39', '2019-07-30 14:02:14', NULL, '分配角色菜单', 6, '/role/setrole', '1156080137621581824', 1, 3, 'RoleSetrolemenu', 'read');
INSERT INTO `admin_menu` VALUES ('1156080138628214784', '2019-07-30 13:52:39', '2019-07-30 14:02:15', NULL, '后台用户管理', 40, '/admins', '1156076295484538880', 1, 2, 'Admins', 'read');
INSERT INTO `admin_menu` VALUES ('1156080139085393920', '2019-07-30 13:52:39', '2019-07-30 14:02:14', NULL, '分配角色', 6, '/admins/setrole', '1156080138628214784', 1, 3, 'AdminsSetrole', 'read');
INSERT INTO `admin_menu` VALUES ('2147483647', '2019-07-30 11:38:36', '2019-07-30 14:22:23', NULL, 'TOP', 1, '', '', 1, 1, 'TOP', 'read');
COMMIT;

-- ----------------------------
-- Table structure for admin_role
-- ----------------------------
DROP TABLE IF EXISTS `admin_role`;
CREATE TABLE `admin_role` (
  `id` varchar(32) NOT NULL COMMENT '主键',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deletedAt` datetime DEFAULT NULL COMMENT '删除时间',
  `name` varchar(64) NOT NULL COMMENT '角色名称',
  `memo` int(11) DEFAULT NULL COMMENT '备注',
  `sequence` int(11) NOT NULL DEFAULT '0' COMMENT '排序值',
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_role_name_uindex` (`name`),
  UNIQUE KEY `admin_role_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='后台用户表';

-- ----------------------------
-- Records of admin_role
-- ----------------------------
BEGIN;
INSERT INTO `admin_role` VALUES ('1', '2019-07-26 15:12:37', '2019-07-26 15:12:37', NULL, 'admin', NULL, 0);
COMMIT;

-- ----------------------------
-- Table structure for admin_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_menu`;
CREATE TABLE `admin_role_menu` (
  `id` varchar(32) NOT NULL COMMENT '主键',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deletedAt` datetime DEFAULT NULL COMMENT '删除时间',
  `role_id` varchar(32) NOT NULL COMMENT '角色id',
  `menu_id` varchar(32) NOT NULL COMMENT '菜单id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_role_menu_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单管理表';

-- ----------------------------
-- Table structure for admin_user
-- ----------------------------
DROP TABLE IF EXISTS `admin_user`;
CREATE TABLE `admin_user` (
  `id` varchar(32) NOT NULL COMMENT '用户uid\n',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deletedAt` datetime DEFAULT NULL COMMENT '删除时间',
  `username` varchar(32) NOT NULL COMMENT '用户名',
  `password` varchar(64) NOT NULL COMMENT '密码',
  `status` tinyint(4) DEFAULT '1' COMMENT '用户状态  1.启用  2.禁用',
  `real_name` varchar(64) DEFAULT NULL COMMENT '真实姓名',
  `email` varchar(64) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(32) DEFAULT NULL COMMENT '手机号',
  `salt` varchar(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_user_uid_uindex` (`id`),
  UNIQUE KEY `admin_user_user_name_uindex` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='后台用户表';

-- ----------------------------
-- Records of admin_user
-- ----------------------------
BEGIN;
INSERT INTO `admin_user` VALUES ('1', '2019-07-26 15:03:21', '2019-07-31 19:03:29', NULL, 'superAdmin', '6d16bc74bb1ff640b83c1d67376e26905d1f981', 1, '', '', '0', '6d16bc7');
INSERT INTO `admin_user` VALUES ('1156520597047054336', '2019-07-31 19:02:53', '2019-07-31 19:07:15', NULL, '测试1', 'p4lx1np03d6c499ebcff930f4ebfb492e3d500a', 1, '', '', '', 'p4lx1np');
INSERT INTO `admin_user` VALUES ('1156521949193695232', '2019-07-31 19:08:15', '2019-07-31 19:08:15', NULL, '1', 'juc8m8ce683b4eff7b20a0ed9bd22a66a7fb6c9', 1, '2', 'string', 'string', 'juc8m8c');
INSERT INTO `admin_user` VALUES ('2', '2019-07-31 16:26:04', '2019-07-31 19:03:29', NULL, 'joe', '6d16bc7ca5e58628a458ad50c0e57c15efd9ee1', 1, '2', '', '', '6d16bc7');
COMMIT;

-- ----------------------------
-- Table structure for admin_user_roles
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_roles`;
CREATE TABLE `admin_user_roles` (
  `id` varchar(32) NOT NULL COMMENT '主键',
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deletedAt` datetime DEFAULT NULL COMMENT '删除时间',
  `user_id` varchar(32) NOT NULL DEFAULT '0' COMMENT '用户id',
  `role_id` varchar(32) NOT NULL DEFAULT '0' COMMENT '角色ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `admin_user_roles_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

SET FOREIGN_KEY_CHECKS = 1;
