# 🚁updata-from-arkdiyt，检测arkadiyt/bounty-targets-data资产信息变动

## 简介

🌟一款[soapffz](https://github.com/soapffz)自用的漏洞赏金资产监测及信息收集框架，基于分布式任务框架[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)，愿景是自动化实现赏金资产的监测及前期信息收集

可能会有部分漏扫功能，在开始使用之前，请务必阅读并同意[免责声明](https://github.com/soapffz/myman/blob/main/Disclaimer.md)中的条款，否则请勿下载安装使用本项目中的所有文件

## 🥐模块功能

原理简述：定时拉取[arkadiyt/bounty-targets-data](https://github.com/arkadiyt/bounty-targets-data)，使用git获取最新一次变动的时间及文件，解析对应变动的文件更新到数据库中

配合[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)设置定时任务可实现自动化完成[arkadiyt/bounty-targets-data](https://github.com/arkadiyt/bounty-targets-data)赏金资产监测

## 💫使用方法

前置条件：
   1. 在mysql数据库中导入[db_model/arkadiyt_data_model.sql](https://github.com/soapffz/myman/tree/main/xxl-job-executor/toolkit/updata-from-arkadiyt/db_model/arkadiyt_data_model.sql)文件
   2. 在[configs](https://github.com/soapffz/myman/tree/main/xxl-job-executor/toolkit/updata-from-arkadiyt/configs/)文件夹中复制一份`config-example.toml`修改为`config.toml`，按照自己配置修改即可使用

## 💥演示截图

*2022年9月10日，旧版只更新通配符域名资产的截图*

数据库截图-2022年9月10日

![image](./images/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%88%AA%E5%9B%BE-2022%E5%B9%B49%E6%9C%8810%E6%97%A511%E7%82%B946.png)

xxl-job调度任务截图-2022年9月10日

![image](./images/xxl-job%E6%89%A7%E8%A1%8C%E6%88%AA%E5%9B%BE1-2022%E5%B9%B49%E6%9C%8810%E6%97%A511%E7%82%B954.png)

![image](./images/xxl-job%E6%89%A7%E8%A1%8C%E6%88%AA%E5%9B%BE2-2022%E5%B9%B49%E6%9C%8810%E6%97%A511%E7%82%B954.png)

![image](./images/xxl-job%E6%89%A7%E8%A1%8C%E6%88%AA%E5%9B%BE3-2022%E5%B9%B49%E6%9C%8810%E6%97%A511%E7%82%B954.png)

## 🍝更新日志

 - 2022-10-16
    - [update]项目由update-wildcard-domain只监测[bounty-targets-data](https://github.com/arkadiyt/bounty-targets-data)库中的[wildcards.txt](https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/wildcards.txt)通配符文件变为监测整个项目，重新设计数据库文件，根据git文件获取提交的时间，配合xxl-job定时任务,解析对应资产到数据库中

 - 2022-09-10
    - [add] 首次添加代码，功能为从[bounty-targets-data](https://github.com/arkadiyt/bounty-targets-data)库中的[wildcards.txt](https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/wildcards.txt)文件中监测更新通配符域名