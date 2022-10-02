# MyMan v0.1
# Copyright (c) 2022-present, soapffz.
CREATE database if NOT EXISTS `bounty` default character set utf8mb4 collate utf8mb4_unicode_ci;

use `bounty`;

SET
  NAMES utf8mb4;

DROP TABLE IF EXISTS `bounty_asset_wildcarddomain`;

CREATE TABLE `bounty_asset_wildcarddomain` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `rootdomain` varchar(255) unique NOT NULL COMMENT '根域名',
  `createtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatetime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `remark` varchar(255) NULL COMMENT '一些自定义备注',
  `activemark` tinyint(1) NULL DEFAULT 1 COMMENT '是否有效，默认1为有效，失效则置为0',
  `addsource` varchar(255) NULL DEFAULT "bounty-targets-data" COMMENT '添加来源,默认为bounty-targets-data',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

INSERT INTO
  `bounty_asset_wildcarddomain`(
    `id`,
    `rootdomain`,
    `createtime`,
    `updatetime`,
    `remark`,
    `activemark`,
    `addsource`
  )
VALUES
  (
    1,
    'baidu.com',
    NOW(),
    NOW(),
    '测试数据',
    0,
    'soapffz'
  );

commit;