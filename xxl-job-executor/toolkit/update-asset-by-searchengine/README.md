# 🚁update-asset-by-searchengine，使用搜索引擎监控资产更新

## 简介

🌟一款[soapffz](https://github.com/soapffz)个人自用的漏扫及告警流程，用于批量刷洞

## 🥐模块功能

从网络空间搜索引擎下载每日更新数据更新到数据库，数据库模型为[soapffz/myman/bounty_asset](https://github.com/soapffz/myman/blob/main/bounty-database/bounty_asset.sql)

目前使用[360quake](https://quake.360.cn/quake/#/index)API进行数据查询

配合[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)设置定时任务可实现自动化指定资产每日更新

 - 获取到没有的ip:port时则插入，已存在的数据则更新时间(设置了ip:port唯一联合索引)

## 🍣使用方法

### 使用前置条件

1.需要已有按数据库模型[soapffz/myman/bounty_asset](https://github.com/soapffz/myman/blob/main/bounty-database/bounty_asset.sql)创建的bounty数据库及相应表

2.在configs/文件夹中复制一份config-example.toml修改为config.toml，按照自己配置修改

### 快速使用方法

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


### 完整参数

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

 - 2022-10-05
    - 1.[add] 新增-ua参数，一键更新配置中所有资产，配合-da下载所有数据参数，可一键下载并更新所有app的当日数据
    - 2.[update] 更新了代码架构、一些小的测试用例，将单次最大下载放宽至10000条数据，不传入关联app关键词将提示，但不会阻止程序运行

 - 2022-10-02，根据[quake_go](https://github.com/360quake/quake_go)项目更改架构
    - 修复：根据原作者代码修改相关字段
    - 还存在的问题：见[issue](https://github.com/360quake/quake_go/issues/14)，待作者修复此bug后再对应修改，本地使用可先强行注释本地库文件中相关代码
    - 还需优化的点：对于循环次数的判断，还没有太好的方法直接读取到行数或者数据返回大小

 - 2022-09-21，更新代码架构，把常用函数全部放到[go_common_functions](go_common_functions),并更换使用[quake_go](https://github.com/360quake/quake_go)，避免了命令执行保存到本地又再次解析本地文件的麻烦

 - 2022-09-12，首次添加代码，使用的是(quake_rs](https://github.com/360quake/quake_rs)