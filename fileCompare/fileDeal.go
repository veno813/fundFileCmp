package filecompare

import (
	"bufio"
	"fmt"
	fileconfig "fundFileCmp/fileConfig"
	"fundFileCmp/logging"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"time"
)

//最多支持500个销售商
var hsDisList [500]string
var lsDisList [500]string
var cmpDisList [500]string

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
	cmpDisFile(filePath)
	//开始比对相同销售商下的各个文件
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

//路径先写死，有需要再调整
func ReadFilePath() {

}

//判断是否为销售商文件夹
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
	var i = 0
	//使用map存储，1：仅恒生存在；2：仅融先存在；9：两边均有
	if hsDisCount != lsDisCount {
		logging.Info("融先导出销售商个数:", lsDisCount)
		logging.Info("恒生导出销售商个数:", hsDisCount)
		logging.Error("销售商数量不一致")
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
			cmpDisList[i] = k
			i++
		}
	}
	logging.Error("恒生多出销售商代码", hsDif)
	logging.Error("融先多出销售商代码", lsDif)
}

func cmpDisFile(filePath string) {
	start := time.Now()
	for k, v := range cmpDisList {
		if v == "" {
			break
		}
		logging.Info("开始", v, "销售商数据比对")
		cmpSingelDisFile(filePath, v)
		fmt.Println("第", k+1, "个销售商文件比对开始")
	}
	cost := time.Since(start)
	fmt.Println("总计用时:", cost)
}

//比对单个销售商下所有文件内容
func cmpSingelDisFile(filePath string, disCode string) {
	//恒生文件夹路径
	var hsFilePath string = filePath + "/hs/output/" + disCode + "/20210712"
	//融先文件夹路径
	var lsFilePath string = filePath + "/ls/output/" + disCode
	//比对确认文件夹
	var confirmPath string = "/Confirm/"
	//比对行情文件夹
	//var funddayPath string = "/FundDay"
	var hsFileList []string
	var lsFileList []string
	var cmpFileList []string

	hsFileList = getFileList(hsFilePath + confirmPath)
	lsFileList = getFileList(lsFilePath + confirmPath)

	//fmt.Println("恒生文件列表", hsFileList)
	//fmt.Println("融先文件列表", lsFileList)
	cmpFileList = cmpDisFileNum(hsFileList, lsFileList)

	//开始比对单个文件
	for _, v := range cmpFileList {
		cmpCfmFile(hsFilePath+confirmPath+v, lsFilePath+confirmPath+v)
	}

}

//获取文件列表
func getFileList(filePath string) []string {
	files, err := ioutil.ReadDir(filePath)
	var FileList []string
	if err != nil {
		logging.Fatal("读取文件错误")
	}
	//跳过.OK文件
	for _, file := range files {
		if path.Ext(file.Name()) == ".ok" || path.Ext(file.Name()) == ".OK" {
			continue
		}
		FileList = append(FileList, file.Name())
	}
	return FileList
}

//比对缺失文件
func cmpDisFileNum(hsFileList []string, lsFileList []string) (cmpFileList []string) {
	//fmt.Println("hs", hsFileList)
	//fmt.Println("ls", lsFileList)
	var result map[string]int = make(map[string]int)
	var hsDif []string = make([]string, 0)
	var lsDif []string = make([]string, 0)
	//使用map存储，1：仅恒生存在；2：仅融先存在；9：两边均有
	if len(hsFileList) != len(lsFileList) {
		logging.Info("融先确认文件数:", len(lsFileList))
		logging.Info("恒生确认文件数:", len(hsFileList))
		logging.Error("确认文件数量不一致")
	}
	for _, v := range hsFileList {
		if v == "" {
			break
		}
		result[v] = 1
	}
	for _, v := range lsFileList {
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
			cmpFileList = append(cmpFileList, k)
		}
	}
	logging.Error("恒生多出确认文件", hsDif)
	logging.Error("融先多出确认文件", lsDif)

	return cmpFileList
}

func cmpCfmFile(hsfile string, lsfile string) bool {
	//fmt.Println(hsfile)
	//fmt.Println(lsfile)
	file, err := os.Open(hsfile)
	if err != nil {
		logging.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		fmt.Println(lineText)
	}

	return true
}
