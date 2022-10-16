# 🚁monsgdata-and-scanwithnuclei

## 🌚简介

🌟一款[soapffz](https://github.com/soapffz)自用的漏洞赏金资产监测及信息收集框架，基于分布式任务框架[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)，愿景是自动化实现赏金资产的监测及前期信息收集

可能会有部分漏扫功能，在开始使用之前，请务必阅读并同意[免责声明](https://github.com/soapffz/myman/blob/main/Disclaimer.md)中的条款，否则请勿下载安装使用本项目中的所有文件

## 🥩模块功能

原理简述：监测从[pdata-by-searchengine](https://github.com/soapffz/myman/tree/main/xxl-job-executor/toolkit/updata-by-searchengine)模块中更新的网络空间搜索引擎数据，并使用nuclei扫描

配合[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)设置定时任务可实现自动化完成资产扫描并生成报告方便提交

## 🍣使用方法

前置条件：
   1. 在mysql数据库中导入[db_model/searchengine_asset_model.sql](https://github.com/soapffz/myman/tree/main/xxl-job-executor/toolkit/monsgdata-and-scanwithnuclei/db_model/searchengine_asset_model.sql)文件
   2. 在[configs](https://github.com/soapffz/myman/tree/main/xxl-job-executor/toolkit/monsgdata-and-scanwithnuclei/configs/)文件夹中复制一份`config-example.toml`修改为`config.toml`，按照自己配置修改即可使用

### 🌹快速使用方法

```
-vp string，扫描指定app关键词，需与数据库中相同

-sall，扫描配置文件中的所有关联app类型，默认关闭

-gen，是否生成漏洞提交报告模版，默认关闭

```

 - -sall和-vp参数共用时，以-sall为准，-gen参数可在两种模式下均使用

## 🧆演示截图

## 🍝更新日志

 - 2022-10-16
      - [update] 更新模块名称`scan-when-asset-add`为`monsgdata-and-scanwithnuclei`，去掉connectmysqldb模块（已集成到[soapffz/common-go-functions](https://github.com/soapffz/common-go-functions/blob/main/pkg/getmysqldbconnbygorm.go)）模块中，优化数据创表文件并放到本模块中

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