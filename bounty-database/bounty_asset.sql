# MyMan v0.1
# Copyright (c) 2022-present, soapffz.
CREATE database if NOT EXISTS `bounty` default character set utf8mb4 collate utf8mb4_unicode_ci;

use `bounty`;

SET
  NAMES utf8mb4;

DROP TABLE IF EXISTS `bounty_asset`;

CREATE TABLE `bounty_asset` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `ip` varchar(255) COMMENT '资产ip',
  `port` varchar(255) COMMENT '资产端口',
  `rootdomain` varchar(255) COMMENT '根域名',
  `domain` varchar(255) COMMENT '域名',
  `asset_url` varchar(255) COMMENT '资产url，例如111.111.111.111:8443',
  `full_url` varchar(255) COMMENT 'url全量链接',
  `vuln_url` varchar(255) COMMENT '其他链接，如漏洞链接',
  `relatedapp` varchar(255) COMMENT '关联的app名称',
  `activemark` tinyint(1) DEFAULT 1 COMMENT '是否存活，默认1为有效，失效则置为0',
  `remark` varchar(255) COMMENT '一些自定义备注',
  `addsource` varchar(255) DEFAULT "quake" COMMENT '添加来源,默认为quake',
  `createtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatetime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

# 设置唯一联合索引，保证ip:port字段唯一
alter table
  bounty_asset
add
  unique index_asset(ip, port);

INSERT INTO
  `bounty_asset`(
    `id`,
    `ip`,
    `port`,
    `rootdomain`,
    `domain`,
    `asset_url`,
    `full_url`,
    `vuln_url`,
    `relatedapp`,
    `activemark`,
    `remark`,
    `addsource`,
    `createtime`,
    `updatetime`
  )
VALUES
  (
    1,
    '1.1.1.1',
    '8443',
    'baidu.com',
    'map.baidu.com',
    '1.1.1.1:8443',
    '1.1.1.1:8443/.git',
    'https://www.baidu.com:8443/.git',
    'git_leak',
    1,
    '测试数据',
    'soapffz',
    NOW(),
    NOW()
  );

commit;