package filecompare

import (
	"fmt"
	fileconfig "fundFileCmp/fileConfig"
	"io/ioutil"
	"log"
	"regexp"
)

//读取文件
func ReadFile(n int, filePath string) {
	fmt.Println(fileconfig.FileType[n], "读取文件中")
	fmt.Println("文件路径", filePath)
	//读取融先下所有文件夹
	getFolderLists(filePath + "/ls/output")
	//读取恒生下所有文件夹
	getFolderLists(filePath + "/hs/output")
}

func getFolderLists(filePath string) {
	files, err := ioutil.ReadDir(filePath)
	agencyCount := 0
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if isDisFolder(file.Name()) {
			fmt.Println(file.Name(), "是销售商文件夹")
			agencyCount++
		} else {
			fmt.Println(file.Name(), "不是销售商文件夹")
		}
	}
	fmt.Println("共有", agencyCount, "个销售商")
}

func ReadFilePath() {

}

func isDisFolder(name string) bool {
	pattern := "\\d+"
	result, _ := regexp.MatchString(pattern, name)
	return result
}
