# MyMan v0.1
# Copyright (c) 2022-present, soapffz.
CREATE database if NOT EXISTS `bounty` default character set utf8mb4 collate utf8mb4_unicode_ci;

use `bounty`;

SET
  NAMES utf8mb4;

CREATE TABLE if not exists `arkadiyt_hackerone` (
  `uid` integer NOT NULL AUTO_INCREMENT COMMENT '自增id，因为原数据中已经有id字段了',
  `allowsbountysplitting` tinyint(1) null COMMENT '是否允许分钱',
  `averagetimetobountyawarded` float(10, 1) COMMENT '平均给钱时间',
  `averagetimetofirstprogramresponse` float(10, 1) COMMENT '平均第一反应时间',
  `averagetimetoreportresolved` float(10, 1) COMMENT '平均报告解决时间',
  `handle` varchar(255) COMMENT '项目简称',
  `id` integer COMMENT '项目id',
  `managedprogram` tinyint(1) null COMMENT '是否由hackerone托管',
  `name` varchar(255) COMMENT '项目名称',
  `offersbounties` tinyint(1) null COMMENT '是否提供赏金',
  `offersswag` tinyint(1) null COMMENT '是否提供礼品',
  `responseefficiencypercentage` integer COMMENT '报告解决率',
  `submissionstate` varchar(255) COMMENT '可提交状态',
  `url` varchar(255) COMMENT '项目链接',
  `website` varchar(255) COMMENT '项目主要站点',
  `inscope` tinyint(1) null COMMENT '是否在范围内',
  `assetidentifier` varchar(255) COMMENT '子资产标识符',
  `assettype` varchar(255) COMMENT '子资产类型',
  `availabilityrequirement` varchar(255) COMMENT '子资产要求',
  `confidentialityrequirement` varchar(255) COMMENT '子资产资质要求',
  `eligibleforbounty` tinyint(1) null COMMENT '是否可以获得赏金',
  `eligibleforsubmission` tinyint(1) null COMMENT '是否可以提交',
  `instruction` varchar(255) COMMENT '子资产介绍',
  `integrityrequirement` varchar(9011) COMMENT '子资产另外要求',
  `maxseverity` varchar(255) COMMENT '子资产最大漏洞等级',
  `createtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatetime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `activemark` tinyint(1) NULL DEFAULT 1 COMMENT '资产是否有效，默认1为有效，失效则置为0',
  `addsource` varchar(255) NULL DEFAULT "bounty-targets-data" COMMENT '添加来源,默认为bounty-targets-data',
  PRIMARY KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE if not exists `arkadiyt_bugcrowd` (
  `id` integer NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `name` varchar(255) COMMENT '项目名称',
  `url` varchar(255) COMMENT '项目链接',
  `allowsdisclosure` tinyint(1) null COMMENT '是否允许披露漏洞',
  `managedbybugcrowd` tinyint(1) null COMMENT '是否由bugcrowd托管',
  `safeharbor` varchar(255) COMMENT '是否为安全港',
  `maxpayout` integer COMMENT '最高漏洞付款金额',
  `inscope` tinyint(1) null COMMENT '是否在范围内',
  `type` varchar(255) COMMENT '子项目类型',
  `target` varchar(255) COMMENT '子项目目标',
  `createtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatetime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `activemark` tinyint(1) NULL DEFAULT 1 COMMENT '资产是否有效，默认1为有效，失效则置为0',
  `addsource` varchar(255) NULL DEFAULT "bounty-targets-data" COMMENT '添加来源,默认为bounty-targets-data',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE if not exists `arkadiyt_intigriti` (
  `uid` integer NOT NULL AUTO_INCREMENT COMMENT '自增id，因为原数据中已经有id字段了',
  `id` varchar(255) COMMENT '项目唯一标识符号',
  `name` varchar(255) COMMENT '项目简称',
  `companyhandle` varchar(255) COMMENT '公司简称',
  `handle` varchar(255) COMMENT '标识',
  `url` varchar(255) COMMENT '项目链接',
  `status` varchar(255) COMMENT '项目状态',
  `confidentialitylevel` varchar(255) COMMENT '项目保密等级',
  `minbounty` float(10, 1) COMMENT '最小赏金',
  `minbountycurrency` varchar(255) COMMENT '最小赏金币种',
  `maxbounty` float(10, 1) COMMENT '最大赏金',
  `maxbountycurrency` varchar(255) COMMENT '最大赏金币种',
  `inscope` tinyint(1) null COMMENT '是否在范围内',
  `type` varchar(255) COMMENT '子项目类型',
  `endpoint` varchar(255) COMMENT '子项目目标',
  `description` varchar(255) COMMENT '子项目描述',
  `createtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatetime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `activemark` tinyint(1) NULL DEFAULT 1 COMMENT '资产是否有效，默认1为有效，失效则置为0',
  `addsource` varchar(255) NULL DEFAULT "bounty-targets-data" COMMENT '添加来源,默认为bounty-targets-data',
  PRIMARY KEY (`uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

CREATE TABLE if not exists `arkadiyt_wildcard` (
  `id` integer NOT NULL AUTO_INCREMENT COMMENT '自增id，因为原数据中已经有id字段了',
  `rootdomainofwildcard` varchar(255) NULL COMMENT '通配符域名的根域名（有些通配符域名是二级域名）',
  `wildcarddomain` varchar(255) unique NULL COMMENT '通配符域名',
  `blackflag` tinyint(1) NULL DEFAULT 0 COMMENT '域名黑名单标记，0为正常，1为黑名单，默认为0',
  `createtime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updatetime` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `activemark` tinyint(1) NULL DEFAULT 1 COMMENT '资产是否有效，默认1为有效，失效则置为0',
  `addsource` varchar(255) NULL DEFAULT "bounty-targets-data" COMMENT '添加来源,默认为bounty-targets-data',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;

# 每次执行都删除索引再重建
# hackerone
DROP INDEX index_h1 ON arkadiyt_hackerone;
# bugcorwd
DROP INDEX index_bc ON arkadiyt_bugcrowd;
# intigriti
DROP INDEX index_ig on arkadiyt_intigriti;

# 设置唯一联合索引，保证字段唯一
# hackerone
alter table
  arkadiyt_hackerone
add
  unique index_h1(handle, name, assetidentifier);

# bugcorwd
alter table
  arkadiyt_bugcrowd
add
  unique index_bc(name, type, target);

# intigriti
alter table
  arkadiyt_intigriti
add
  unique index_ig(id, type, endpoint);

INSERT
  IGNORE INTO `arkadiyt_hackerone`(
    `uid`,
    `allowsbountysplitting`,
    `averagetimetobountyawarded`,
    `averagetimetofirstprogramresponse`,
    `averagetimetoreportresolved`,
    `handle`,
    `id`,
    `managedprogram`,
    `name`,
    `offersbounties`,
    `offersswag`,
    `responseefficiencypercentage`,
    `submissionstate`,
    `url`,
    `website`,
    `inscope`,
    `assetidentifier`,
    `assettype`,
    `availabilityrequirement`,
    `confidentialityrequirement`,
    `eligibleforbounty`,
    `eligibleforsubmission`,
    `instruction`,
    `integrityrequirement`,
    `maxseverity`,
    `createtime`,
    `updatetime`,
    `activemark`,
    `addsource`
  )
VALUES
  (
    1,
    1,
    200.0,
    7.0,
    100.0,
    "test-bounty",
    114514,
    1,
    "this is a test bounty name",
    1,
    0,
    98,
    "open",
    "https://hackerone.com/baidu",
    "https://baidu.com",
    1,
    "114514890013",
    "APPLE_STORE_APP_ID",
    null,
    null,
    1,
    1,
    "https://baidu.com/instructions.html",
    null,
    "critical",
    NOW(),
    NOW(),
    0,
    "soapffz测试数据"
  );

INSERT
  IGNORE INTO `arkadiyt_bugcrowd`(
    `id`,
    `name`,
    `url`,
    `allowsdisclosure`,
    `managedbybugcrowd`,
    `safeharbor`,
    `maxpayout`,
    `inscope`,
    `type`,
    `target`,
    `createtime`,
    `updatetime`,
    `activemark`,
    `addsource`
  )
VALUES
  (
    1,
    "测试bugcrowd项目",
    "https://bugcorwd.com/soaptest",
    0,
    1,
    "full",
    114514,
    1,
    "website",
    "soapffz.com",
    NOW(),
    NOW(),
    0,
    "soapffz测试数据"
  );

INSERT
  IGNORE INTO `arkadiyt_intigriti`(
    `uid`,
    `id`,
    `name`,
    `companyhandle`,
    `handle`,
    `url`,
    `status`,
    `confidentialitylevel`,
    `minbounty`,
    `minbountycurrency`,
    `maxbounty`,
    `maxbountycurrency`,
    `inscope`,
    `type`,
    `endpoint`,
    `description`,
    `createtime`,
    `updatetime`,
    `activemark`,
    `addsource`
  )
VALUES
  (
    1,
    "1111-22222-3333--4444",
    "soapffz intigriti 测试项目",
    "soapcompany",
    "soapintigriti",
    "https://www.intigriti.com/programs/soaposoap",
    "close",
    "public",
    0.0,
    "EUR",
    1000.0,
    "EUR",
    0,
    "url",
    "soapffz.com",
    null,
    NOW(),
    NOW(),
    0,
    "soapffz测试数据"
  );

INSERT
  IGNORE INTO `arkadiyt_wildcard`(
    `id`,
    `rootdomainofwildcard`,
    `wildcarddomain`,
    `blackflag`,
    `createtime`,
    `updatetime`,
    `activemark`,
    `addsource`
  )
VALUES
  (
    1,
    "baidu.com",
    "wenku.baidu.com",
    1,
    NOW(),
    NOW(),
    0,
    "soapffz测试数据"
  );

commit;