# 🌸 xxl-job-executor说明

 - 基于官方[xxl-job-executor-sample-springboot](https://github.com/xuxueli/xxl-job/tree/2.3.1/xxl-job-executor-samples/xxl-job-executor-sample-springboot)改造，指定自定义镜像，加入自定义工具，重新打包

  - 使用方法：mvn编译官方执行器代码后，复制生成的xxl-job-executor-sample-springboot-*.jar文件到xxl-job-executor文件夹后使用docker打包即可：

```shell
docker build . -t xxl-job-executor:latest
```

执行器启动默认端口都为9999，如果需要启动多个第二个注意端口映射

启动执行器

```shell
docker run --privileged=true -d -p 9999:9999 xxl-job-executor:latest
```

第二个的话就得这样

```shell
docker run --privileged=true -d -p 9998:9999 xxl-job-executor:latest
```