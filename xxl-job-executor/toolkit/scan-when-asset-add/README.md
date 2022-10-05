# 🚁scan-when-asset-add，数据更新后调用指定模块进行扫描

## 🌚简介

🌟一款[soapffz](https://github.com/soapffz)个人自用的漏扫及告警流程，用于批量刷洞

## 🥩模块功能

传入数据库[soapffz/myman/bounty_asset](https://github.com/soapffz/myman/blob/main/bounty-database/bounty_asset.sql)中的标签，根据标签类型，搜索指定时间内更新的数据，根据标签调用对应扫描模块([自用nuclei脚本](https://github.com/soapffz/myown-nuclei-poc))进行扫描，并进行ip域名解析及根域名权重查询，方便批量刷洞

配合[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)设置定时任务可实现自动化完成定时漏洞监测及通知

## 🍣使用方法

### 🪷使用前置条件

1.需要已有按数据库模型[soapffz/myman/bounty_asset](https://github.com/soapffz/myman/blob/main/bounty-database/bounty_asset.sql)创建的bounty数据库及相应表

2.在configs/文件夹中复制一份config-example.toml修改为config.toml，按照自己配置修改



### 🌹快速使用方法

```
-vp string，扫描指定app关键词，需与数据库中相同

-sall，扫描配置文件中的所有关联app类型，默认关闭

-gen，是否生成漏洞提交报告模版，默认关闭

```

 - -sall和-vp参数共用时，以-sall为准，-gen参数可在两种模式下均使用

## 🧆演示截图

## 🍝更新日志

 - 2022-10-05
      - [update] 重构扫描所有app部分的功能架构
      - [update] 更新获取nuclei扫出结果解析部分代码

 - 2022-10-04
      - [add] 增加-sall参数，开启时候直接启动所有已知标签的扫描，减少xxl-job添加多个任务的麻烦
      - [update] 解耦扫描函数，传入关联app类型、数据结构体、及nuclei脚本地址，扫描有漏洞后根据ip、域名、关联app名称写入数据库
      - [update] nuclei扫描结果不再输出到文件，直接每次扫描取扫描结果

 - 2022-10-03
    - [add] 添加ip解析为网站及查询根域名权重功能
    - [add] 添加通过Server酱推送消息的代码
    - [update] 优化代码架构

 - 2022-09-21
    - [update] 更新代码架构，添加yongyou_nc的nuclei模版

 - 2022-09-18
    - [add] 首次添加代码