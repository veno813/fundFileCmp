package logging

import (
	"fmt"
	"fundFileCmp/file"
	"os"
	"time"
)

//获取日志路径
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", "runtime/", "logs/")
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		"Log",
		time.Now().Format("20060102"),
		"txt")
}

func openLogFile(filename, filepath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err :%v", err)
	}
	src := dir + "/" + filepath
	perm := file.CheckPermission(src)
	if perm {
		return nil, fmt.Errorf("文件权限检查失败,%s", src)
	}
	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("不存在改文件路径 %s,错误原因 %s", src, err)
	}

	f, err := file.Open(src+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("文件打开失败,%s", err)
	}

	return f, nil
}
