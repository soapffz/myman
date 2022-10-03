# 🚁scan-when-asset-add，数据更新后调用指定模块进行扫描

## 🌚简介

🌟一款[soapffz](https://github.com/soapffz)个人自用的漏扫及告警流程，用于批量刷洞

## 🥩模块功能

传入数据库[soapffz/myman/bounty_asset](https://github.com/soapffz/myman/blob/main/bounty-database/bounty_asset.sql)中的标签，根据标签类型，搜索指定时间内更新的数据，根据标签调用对应扫描模块([自用nuclei脚本](https://github.com/soapffz/myown-nuclei-poc))进行扫描，并进行ip域名解析及根域名权重查询，方便批量刷洞

配合[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)设置定时任务可实现自动化完成定时漏洞监测及通知

## 🥙使用方法

前置条件：在configs/文件夹中复制一份config-example.toml修改为config.toml，按照自己配置修改即可使用

命令行使用:

```
-all
      是否生成漏洞提交报告模版，默认关闭
-vp string
      对应的app关键词，需与数据库中相同
```

## 🧆演示截图

## 🍝更新日志

 - 2022-10-03
    - 1.添加ip解析为网站及查询根域名权重的代码
    - 2.添加通过Server酱推送消息的代码
    - 3.优化代码架构

 - 2022-09-18,首次添加扫描代码

 - 2022-09-21,更新代码架构，添加yongyou_nc的nuclei模版
