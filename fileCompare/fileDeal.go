package filecompare

import (
	"fmt"
	fileconfig "fundFileCmp/fileConfig"
	"fundFileCmp/logging"
	"io/ioutil"
	"log"
	"regexp"
)

//最多支持500个销售商
var hsDisList [500]string
var lsDisList [500]string

var hsDisCount, lsDisCount int

//读取文件
func ReadFile(n int, filePath string) {
	fmt.Println(fileconfig.FileType[n], "读取文件中")
	fmt.Println("文件路径", filePath)
	//读取融先下所有文件夹
	logging.Info("开始读取数据")
	hsDisCount = getFolderLists(filePath+"/ls/output", &hsDisList)
	//读取恒生下所有文件夹
	lsDisCount = getFolderLists(filePath+"/hs/output", &lsDisList)
	//比对销售商文件数
	cmpDisFolder()
	//比对每个销售商下文件
}

func getFolderLists(filePath string, DisList *[500]string) int {
	files, err := ioutil.ReadDir(filePath)
	agencyCount := 0
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if isDisFolder(file.Name()) {
			DisList[agencyCount] = file.Name()
			agencyCount++
		} else {
			//fmt.Println(file.Name(), "不是销售商文件夹")
		}
	}
	return agencyCount
}

func ReadFilePath() {

}

func isDisFolder(name string) bool {
	pattern := "\\d+"
	result, _ := regexp.MatchString(pattern, name)
	return result
}

func cmpDisFolder() {
	fmt.Println("hs", hsDisCount)
	fmt.Println("ls", lsDisCount)
	var result map[string]int = make(map[string]int)
	var hsDif []string = make([]string, 0)
	var lsDif []string = make([]string, 0)
	//使用map存储，1：仅恒生存在；2：仅融先存在；9：两边均有
	if hsDisCount != lsDisCount {
		logging.Info("融先导出销售商个数:", lsDisCount)
		logging.Info("恒生导出销售商个数:", hsDisCount)
		logging.Info("销售商数量不一致")
	}
	for _, v := range hsDisList {
		if v == "" {
			break
		}
		result[v] = 1
	}
	for _, v := range lsDisList {
		if result[v] == 1 {
			result[v] = 9
		} else {
			result[v] = 2
		}
	}

	//输出比对后结果
	for k, v := range result {
		switch v {
		case 1:
			hsDif = append(hsDif, k)
		case 2:
			lsDif = append(lsDif, k)
		case 9:
		}
	}
	logging.Info("恒生多出销售商代码", hsDif)
	logging.Info("融先多出销售商代码", lsDif)
}
