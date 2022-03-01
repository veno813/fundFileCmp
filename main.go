package main

import (
	"fmt"
	filecompare "fundFileCmp/fileCompare"
	fileconfig "fundFileCmp/fileConfig"
	"log"

	"github.com/go-ini/ini"
)

type App struct {
	FilePath string
}

var AppSetting = &App{}

func main() {
	Cfg, err := ini.Load("config/app.ini")
	log.Println("开始解析配置文件conf/app.ini")
	defer log.Println("配置文件conf/app.ini解析完成")
	if err != nil {
		log.Fatalf("conf/app.ini配置文件有误,错误信息:%v", err)
	}

	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("conf/app.ini配置文件[app]中内容有误,错误信息:%v", err)
	}

	var n int = 0
	fmt.Println("请输入需要比对的文件类型")
	for i := 1; i <= 5; i++ {
		fmt.Println(i, ":", fileconfig.FileType[i])
	}
	fmt.Scanln(&n)
	fmt.Println("您选择的文件类型为:")
	fmt.Println(n, ":", fileconfig.FileType[n])

	filecompare.ReadFile(n, AppSetting.FilePath)
}
