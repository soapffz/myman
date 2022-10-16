# 🚁update-asset-by-searchengine，使用搜索引擎监控资产更新

## 🌚简介

🌟一款[soapffz](https://github.com/soapffz)自用的漏洞赏金资产监测及信息收集框架，基于分布式任务框架[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)，愿景是自动化实现赏金资产的监测及前期信息收集

可能会有部分漏扫功能，在开始使用之前，请务必阅读并同意[免责声明](https://github.com/soapffz/myman/blob/main/Disclaimer.md)中的条款，否则请勿下载安装使用本项目中的所有文件


## 🥐模块功能

原理简述：从网络空间搜索引擎下载每日更新数据更新到数据库

目前使用[360quake](https://quake.360.cn/quake/#/index)API进行数据查询

配合[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)设置定时任务可实现自动化指定资产每日更新

 - 获取到没有的ip:port时则插入，已存在的数据则更新时间(设置了ip:port唯一联合索引)


## 🍣使用方法

前置条件：
   1. 在mysql数据库中导入[db_model/searchengine_asset_model.sql](https://github.com/soapffz/myman/tree/main/xxl-job-executor/toolkit/updata-by-searchengine/db_model/searchengine_asset_model.sql)文件
   2. 在[configs](https://github.com/soapffz/myman/tree/main/xxl-job-executor/toolkit/updata-by-searchengine/configs/)文件夹中复制一份`config-example.toml`修改为`config.toml`，按照自己配置修改即可使用


### 🌹快速使用方法

比[quake_go](https://github.com/360quake/quake_go)的使用方法多添加了三个参数

```
-relatedapp,-rp,相关的app,与数据库模型中的relatedapp字段对应[用于单条查询]

-downall_flag,-da，是否下载查询到的所有数据，默认为否，若开启单次最多下载10000条数据[用于单条查询和批量更新]

--updateallapp_flag,-ua，是否更新配置文件中所有app，默认为否，若开启除-da参数外无视其他语句
```

 - 默认请求从今天的0点到使用时间的数据，若需自定义修改时间可注释掉源码中相关语句
 - 使用search指定查询语句及-rp指定关联app时，默认请求10条数据，添加-da选项后每个app都下载当天的所有数据（最多10000条）
 - 使用-ua参数更新config.toml文件中的所有app，默认请求10条数据，添加-da选项后每个app都下载当天的所有数据（最多10000条）

额外提醒：
 - 重复运行同一指令也不要怕
   - quake对于同一语句多次查询的数据，只会对新查询的数据计算API调用次数
   - 本程序使用[gorm](https://gorm.io/zh_CN/docs/index.html)库进行数据库插入时候会根据ip:port唯一联合索引处理重复




### 🌞完整参数

```
usage: update_asset_by_quake [option] [args] [-da,downall_flag bool=false] [-e,end_time time=2022-10-05 22:23:44] [-fe,field string] [-ft,file_txt string] [-h,help bool] [-ic,ignore_cache bool=false] [-rp,relatedapp string] [-sz,size string=10] [-st,start string=0] [-s,start_time time=2022-01-01] [-ua,updateallapp_flag bool=false]

positional options:
       option               [string]                    init,info,search,host
       args                 [string]                    query value,example port:443

options:
  -da, --downall_flag       [bool=false]                -da download all data,default false
   -e, --end_time           [time=2022-10-05 22:23:44]  -e time to end time flag
  -fe, --field              [string]                    -fe swich body,title,host,html_hash,x_powered_by  to show infomation
  -ft, --file_txt           [string]                    -ft ./file.txt file to query search
   -h, --help               [bool]                      show usage
  -ic, --ignore_cache       [bool=false]                -ic true or false,default false
  -rp, --relatedapp         [string]                    -rp related app 
  -sz, --size               [string=10]                 -sz to size number 
  -st, --start              [string=0]                  -st to start number
   -s, --start_time         [time=2022-01-01]           -s time flag , default time is time.now.year
  -ua, --updateallapp_flag  [bool=false]                -ua update all app by query_statement
```

## 🎂演示截图

## 🥃更新日志

 - 2022-10-16
    - [update] 模块从`update-asset-by-searchengine`更名为`updata-by-searchengine`
    - [update] 去掉connectmysqldb模块（已集成到[soapffz/common-go-functions](https://github.com/soapffz/common-go-functions/blob/main/pkg/getmysqldbconnbygorm.go)）模块中
    - [update] 优化数据创表文件并放到本模块中

 - 2022-10-06
    - [update] 更新了对于传入参数的判断逻辑[#4240d9e](https://github.com/soapffz/myman/commit/4240d9e0e0f1a9821a3e97c5e1d6e9f1314d8522)

 - 2022-10-05
    - [add] 新增-ua参数，一键更新配置中所有资产，配合-da下载所有数据参数，可一键下载并更新所有app的当日数据
    - [update] 更新了代码架构、一些小的测试用例，将单次最大下载放宽至10000条数据，不传入关联app关键词将提示，但不会阻止程序运行

 - 2022-10-02
    - [update] 根据[quake_go](https://github.com/360quake/quake_go)项目更改架构
    - [fix] 根据原作者代码修改相关字段
    - [known_issue] 见[issue](https://github.com/360quake/quake_go/issues/14)，待作者修复此bug后再对应修改，本地使用可先强行注释本地库文件中相关代码
    - [known_issue] 对于循环次数的判断，还没有太好的方法直接读取到行数或者数据返回大小

 - 2022-09-21
    - [update] 更新代码架构，把常用函数全部放到[go_common_functions](go_common_functions),并更换使用[quake_go](https://github.com/360quake/quake_go)，避免了命令执行保存到本地又再次解析本地文件的麻烦

 - 2022-09-12
    - [add] 首次添加代码，使用的是(quake_rs](https://github.com/360quake/quake_rs)