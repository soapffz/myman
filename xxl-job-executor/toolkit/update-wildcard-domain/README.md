# 🚁update-wildcard-domain，更新通配符域名

## 功能

解析[bounty-targets-data](https://github.com/arkadiyt/bounty-targets-data)中数据，提取通配符域名并更新到数据库中

配合[xxl-job-executor](https://github.com/soapffz/myman/tree/main/xxl-job-executor)设置定时任务可实现自动化完成通配符域名监测

 - 获取到没有的通配符域名时则插入，已存在的数据则更新时间

## 使用方法

按自己配置填写configs/config-example.toml，修改文件名为config.toml即可使用

## 演示截图

数据库截图-2022年9月10日

![image](./images/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%88%AA%E5%9B%BE-2022%E5%B9%B49%E6%9C%8810%E6%97%A511%E7%82%B946.png)

xxl-job调度任务截图-2022年9月10日

![image](./images/xxl-job%E6%89%A7%E8%A1%8C%E6%88%AA%E5%9B%BE1-2022%E5%B9%B49%E6%9C%8810%E6%97%A511%E7%82%B954.png)

![image](./images/xxl-job%E6%89%A7%E8%A1%8C%E6%88%AA%E5%9B%BE2-2022%E5%B9%B49%E6%9C%8810%E6%97%A511%E7%82%B954.png)

![image](./images/xxl-job%E6%89%A7%E8%A1%8C%E6%88%AA%E5%9B%BE3-2022%E5%B9%B49%E6%9C%8810%E6%97%A511%E7%82%B954.png)