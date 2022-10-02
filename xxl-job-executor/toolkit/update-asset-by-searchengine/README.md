# 🚁update-asset-by-searchengine，使用搜索引擎监控资产更新

## 🥐功能

从网络空间搜索引擎下载每日更新数据更新到数据库

目前使用[360quake](https://quake.360.cn/quake/#/index)API进行数据查询

配合[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)设置定时任务可实现自动化指定资产每日更新

 - 获取到没有的ip:port时则插入，已存在的数据则更新时间(设置了ip:port唯一联合索引)

## 🍣使用方法

与[quake_go](https://github.com/360quake/quake_go)使用方法一致，除此之外还添加了两个参数

```
-relatedapp,-rp,相关的app,与数据库模型中的relatedapp字段对应

-downall_flag,-da，是否下载查询到的所有数据，默认为否，若开启单次最多下载1000条数据
```

按自己配置填写configs/config-example.toml，修改文件名为config.toml

 - 默认请求从今天的0点到使用时间的数据，若需自定义修改时间请自己动手
 - 默认请求10条数据，开启-da(-downall_flag)选项后，下载查询到的所有数据（最多1000条）
 - 重复运行同一指令也不要怕，quake对于同一语句多次查询的数据，只会对新查询的数据计算API调用次数
 - 本程序使用[gorm](https://gorm.io/zh_CN/docs/index.html)库进行数据库插入时候会根据ip:port唯一联合索引处理重复


## 🎂演示截图

## 🥃更新日志

 - 2022-09-21，更新代码架构，把常用函数全部放到[go_common_functions](go_common_functions),并更换使用[quake_go](https://github.com/360quake/quake_go)，避免了命令执行保存到本地又再次解析本地文件的麻烦

 - 2022-09-12，首次添加代码，使用的是(quake_rs](https://github.com/360quake/quake_rs)