package main

import (
	"context"
	"fmt"

	"log"

	xxl "github.com/xxl-job/xxl-job-executor-go"
)

func main() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://192.168.100.21:9080/xxl-job-admin"),
		xxl.AccessToken("default_token"), //请求令牌(默认为空)
		xxl.ExecutorIp("0.0.0.0"),        //可自动获取
		xxl.ExecutorPort("9998"),         //默认9998（非必填）
		xxl.RegistryKey("golang-jobs"),   //执行器名称
		xxl.SetLogger(&logger{}),         //自定义日志
	)
	exec.Init()
	//设置日志查看handler
	exec.LogHandler(func(req *xxl.LogReq) *xxl.LogRes {
		return &xxl.LogRes{Code: 200, Msg: "", Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "这个是自定义日志handler",
			IsEnd:       true,
		}}
	})
	//注册任务handler
	exec.RegTask("Task_Test", Task_Test)
	log.Fatal(exec.Run())
}

//xxl.Logger接口实现
type logger struct{}

func (l *logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

func Task_Test(cxt context.Context, param *xxl.RunReq) (msg string) {
	fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + " log_id:" + xxl.Int64ToStr(param.LogID))
	return "test done"
}
