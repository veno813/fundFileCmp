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
	"strings"
	"time"

	"github.com/axgle/mahonia"
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
	if n == 1 {
		hsDisCount = getFolderLists(filePath+"/ls/output", &hsDisList)
		//读取恒生下所有文件夹
		lsDisCount = getFolderLists(filePath+"/hs/output", &lsDisList)
		//比对销售商文件数
		cmpDisFolder()
		//比对每个销售商下文件
		cmpDisFile(filePath)
		//开始比对相同销售商下的各个文件
	} else if n == 2 {
		//开始比对集中备份文件
		cmpJZBFFile(filePath)
	} else if n == 3 {
		//开始比对客服文件
		cmpKFFile(filePath)
	} else if n == 4 {
		cmpKJQSFile(filePath)
	}
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
	for _, v := range cmpDisList {
		if v == "" {
			break
		}
		logging.Info("开始", v, "销售商数据比对")
		cmpSingelDisFile(filePath, v)
		//fmt.Println("第", k+1, "个销售商文件比对开始")
	}
	cost := time.Since(start)
	fmt.Println("总计用时:", cost)
}

//比对单个销售商下所有文件内容
func cmpSingelDisFile(filePath string, disCode string) {
	//恒生文件夹路径
	var hsFilePath string = filePath + "/hs/output/" + disCode + "/20220309"
	//融先文件夹路径
	var lsFilePath string = filePath + "/ls/output/" + disCode
	//比对确认文件夹
	var confirmPath string = "/Confirm/"
	//比对行情文件夹
	var funddayPath string = "/FundDay/"
	//确认文件夹下数据
	var hsFileList []string
	var lsFileList []string
	var cmpFileList []string
	//行情文件夹下数据

	var hsFunddayFileList []string
	var lsFunddayFileList []string
	var cmpFunddayFileList []string

	hsFileList = getFileList(hsFilePath + confirmPath)
	lsFileList = getFileList(lsFilePath + confirmPath)

	//fmt.Println("恒生文件列表", hsFileList)
	//fmt.Println("融先文件列表", lsFileList)
	cmpFileList = cmpDisFileNum(hsFileList, lsFileList)

	//开始比对单个文件
	for _, v := range cmpFileList {
		logging.Info("开始核对", v, "文件")
		cmpCfmFile(hsFilePath+confirmPath+v, lsFilePath+confirmPath+v, disCode)
		logging.Info(v, "文件核对完成")
	}

	hsFunddayFileList = getFileList(hsFilePath + funddayPath)
	lsFunddayFileList = getFileList(lsFilePath + funddayPath)
	cmpFunddayFileList = cmpDisFileNum(hsFunddayFileList, lsFunddayFileList)

	for _, v := range cmpFunddayFileList {
		logging.Info("开始核对", v, "文件")
		//cmpFunddayFile(hsFilePath+funddayPath+v, lsFilePath+funddayPath+v, disCode)
		logging.Info(v, "文件核对完成")
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

func cmpCfmFile(hsfile string, lsfile string, disCode string) bool {
	//fmt.Println(hsfile)
	//fmt.Println(lsfile)
	var hslineText []string
	var lslineText []string
	//var result bool
	fileH, err := os.Open(hsfile)
	if err != nil {
		logging.Fatal(err)
	}

	scannerHs := bufio.NewScanner(fileH)
	for scannerHs.Scan() {
		//lineText := scanner.Text()
		//gbk读取又乱码，需要转换
		//lineText, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(scannerHs.Text())), simplifiedchinese.GBK.NewEncoder()))
		hslineText = append(hslineText, ConvertToString(scannerHs.Text(), "GBK", "UTF-8"))
	}

	fileL, err := os.Open(lsfile)
	if err != nil {
		logging.Fatal(err)
	}
	scannerLs := bufio.NewScanner(fileL)
	for scannerLs.Scan() {
		//lineText := scanner.Text()
		//lineText := ConvertToString(scannerLs.Text(),"GBK","UTF-8")
		lslineText = append(lslineText, ConvertToString(scannerLs.Text(), "GBK", "UTF-8"))
	}
	if strings.Contains(hsfile, "_02.") {
		//fmt.Println("这个是02文件")
		//fmt.Println(hslineText)
		//fmt.Println(lslineText)
		cmp02File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "_04.") {
		cmp04File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "_05.") {
		//	cmp05File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "_06.") {
		//06后续实现，暂时用全量比对方式比对
		//cmp06File(hslineText, lslineText, disCode)
		//	AllCmpFile(hslineText, lslineText, disCode, "06")
	} else if strings.Contains(hsfile, "_09.") {
		//09采用全量方式比对
		//AllCmpFile(hslineText, lslineText, disCode, "09")
	} else if strings.Contains(hsfile, "_10.") {
		//cmp10File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "_11.") {
		//cmp11File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "_12.") {
		//cmp12File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "_24.") {
		//24采用全量方式比对
		//AllCmpFile(hslineText, lslineText, disCode, "24")
	} else if strings.Contains(hsfile, "_25.") {
		//cmp25File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "_26.") {
		//26采用全量方式比对
		//	AllCmpFile(hslineText, lslineText, disCode, "26")
	} else if strings.Contains(hsfile, "_44.") {
		//44采用全量方式比对
		//AllCmpFile(hslineText, lslineText, disCode, "44")
	} else if strings.Contains(hsfile, "_R2.") {
		//R2采用全量方式比对
		AllCmpFile(hslineText, lslineText, disCode, "R2")
	} else {
		//索引文件进行全量比对
		//有文件缺失，导致索引文件均不一致，先不比对
		AllCmpFile(hslineText, lslineText, disCode, "索引文件")
	}

	return true
}

func cmpFunddayFile(hsfile string, lsfile string, disCode string) bool {
	//fmt.Println(hsfile)
	//fmt.Println(lsfile)
	var hslineText []string
	var lslineText []string
	//var result bool
	fileH, err := os.Open(hsfile)
	if err != nil {
		logging.Fatal(err)
	}

	scannerHs := bufio.NewScanner(fileH)
	for scannerHs.Scan() {
		//lineText := scanner.Text()
		//gbk读取又乱码，需要转换
		//lineText, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(scannerHs.Text())), simplifiedchinese.GBK.NewEncoder()))
		//hslineText = append(hslineText, ConvertToString(scannerHs.Text(), "GBK", "UTF-8"))
		hslineText = append(hslineText, mahonia.NewDecoder("GBK").ConvertString(scannerHs.Text()))
	}

	fileL, err := os.Open(lsfile)
	if err != nil {
		logging.Fatal(err)
	}
	scannerLs := bufio.NewScanner(fileL)
	for scannerLs.Scan() {
		//lineText := scanner.Text()
		//lineText := ConvertToString(scannerLs.Text(),"GBK","UTF-8")
		//lineText :=mahonia.NewDecoder("GBK").ConvertString(scannerLs.Text())
		//lslineText = append(lslineText, ConvertToString(scannerLs.Text(), "GBK", "UTF-8"))
		lslineText = append(lslineText, mahonia.NewDecoder("GBK").ConvertString(scannerLs.Text()))
	}
	if strings.Contains(hsfile, "_07.") {
		//先进行全量比对，拿到数据再改
		AllCmpFile(hslineText, lslineText, disCode, "07")
	} else if strings.Contains(hsfile, "_08.") {
		//公告文件全量比对
		AllCmpFile(hslineText, lslineText, disCode, "08")
	} else if strings.Contains(hsfile, "_21.") {
		cmp21File(hslineText, lslineText, disCode)
	} else {
		//索引文件进行全量比对
		//有文件缺失，导致索引文件均不一致，先不比对
		//AllCmpFile(hslineText, lslineText, disCode, "索引文件")
	}

	return true
}

//GBK转utf8的方法
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func cmpKFFile(filePath string) {
	start := time.Now()

	logging.Info("开始客服文件数据比对")
	cmpSingelKFFile(filePath)
	//fmt.Println("第", k+1, "个销售商文件比对开始")
	cost := time.Since(start)
	fmt.Println("总计用时:", cost)
}

//比对单个销售商下所有文件内容
func cmpSingelKFFile(filePath string) {
	//恒生文件夹路径
	var hsFilePath string = filePath + "/hs/output/KF/20220309/"
	//融先文件夹路径
	var lsFilePath string = filePath + "/ls/output/KF/"

	//客服文件夹下数据
	var hsFileList []string
	var lsFileList []string
	var cmpFileList []string

	hsFileList = getFileList(hsFilePath)
	lsFileList = getFileList(lsFilePath)

	//fmt.Println("恒生文件列表", hsFileList)
	//fmt.Println("融先文件列表", lsFileList)
	cmpFileList = cmpDisFileNum(hsFileList, lsFileList)

	//开始比对单个文件
	for _, v := range cmpFileList {
		logging.Info("开始核对", v, "文件")
		cmpKFSgFile(hsFilePath+v, lsFilePath+v, "KF")
		logging.Info(v, "文件核对完成")
	}

}

func cmpKFSgFile(hsfile string, lsfile string, disCode string) bool {
	//fmt.Println(hsfile)
	//fmt.Println(lsfile)
	var hslineText []string
	var lslineText []string
	//var result bool
	fileH, err := os.Open(hsfile)
	if err != nil {
		logging.Fatal(err)
	}

	scannerHs := bufio.NewScanner(fileH)
	for scannerHs.Scan() {
		//lineText := scanner.Text()
		//gbk读取又乱码，需要转换
		//lineText, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(scannerHs.Text())), simplifiedchinese.GBK.NewEncoder()))
		hslineText = append(hslineText, ConvertToString(scannerHs.Text(), "GBK", "UTF-8"))
	}

	fileL, err := os.Open(lsfile)
	if err != nil {
		logging.Fatal(err)
	}
	scannerLs := bufio.NewScanner(fileL)
	for scannerLs.Scan() {
		//lineText := scanner.Text()
		//lineText := ConvertToString(scannerLs.Text(),"GBK","UTF-8")
		lslineText = append(lslineText, ConvertToString(scannerLs.Text(), "GBK", "UTF-8"))
	}
	if strings.Contains(hsfile, "Acco_") {
		//暂时存在问题，后续进行比对
	} else if strings.Contains(hsfile, "Acconet_") {
		//cmpAcconetFile(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "AccoRequest_") {
		//cmpAccoReqFile(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "Confirm_") {
		cmpKFCfmFile(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "ConfirmDetail_") {
		//生成确认单号规则不一致，暂时无法比对
		//AllCmpFile(hslineText, lslineText, disCode, "09")
	} else if strings.Contains(hsfile, "Dividend_") {
		//cmp10File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "FundInfo_") {
		//cmp11File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "Request_") {
		//cmp12File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "Share_") {
		//太大了，读不了
		//AllCmpFile(hslineText, lslineText, disCode, "24")
	} else if strings.Contains(hsfile, "ShareCurrents_") {
		//cmp25File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "ShareDetail_") {
		//更大了，也对不了
		//AllCmpFile(hslineText, lslineText, disCode, "26")
	}

	return true
}

func cmpJZBFFile(filePath string) {
	start := time.Now()

	logging.Info("开始比对集中备份数据比对")
	cmpSingelJZBFFile(filePath)
	//fmt.Println("第", k+1, "个销售商文件比对开始")
	cost := time.Since(start)
	fmt.Println("总计用时:", cost)
}

//比对单个销售商下所有文件内容
func cmpSingelJZBFFile(filePath string) {
	//恒生文件夹路径
	var hsFilePath string = filePath + "/hs/output/JZBF/20220309/"
	//融先文件夹路径
	var lsFilePath string = filePath + "/ls/output/JZBF/"

	//客服文件夹下数据
	var hsFileList []string
	var lsFileList []string
	var cmpFileList []string

	hsFileList = getFileList(hsFilePath)
	lsFileList = getFileList(lsFilePath)

	//fmt.Println("恒生文件列表", hsFileList)
	//fmt.Println("融先文件列表", lsFileList)
	cmpFileList = cmpDisFileNum(hsFileList, lsFileList)

	//开始比对单个文件
	for _, v := range cmpFileList {
		logging.Info("开始核对", v, "文件")
		if strings.Contains(v, "_92.") {
			cmpJZBFSgFile(hsFilePath+v, lsFilePath+v, "JZBF")
		}
		logging.Info(v, "文件核对完成")
	}

}

func cmpJZBFSgFile(hsfile string, lsfile string, disCode string) bool {
	//fmt.Println(hsfile)
	//fmt.Println(lsfile)
	var hslineText []string
	var lslineText []string
	//var result bool
	fileH, err := os.Open(hsfile)
	if err != nil {
		logging.Fatal(err)
	}

	scannerHs := bufio.NewScanner(fileH)
	for scannerHs.Scan() {
		hslineText = append(hslineText, scannerHs.Text())
		//gbk读取又乱码，需要转换
		//lineText, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(scannerHs.Text())), simplifiedchinese.GBK.NewEncoder()))
		//hslineText = append(hslineText, ConvertToString(scannerHs.Text(), "GBK", "UTF-8"))
	}

	fileL, err := os.Open(lsfile)
	if err != nil {
		logging.Fatal(err)
	}
	scannerLs := bufio.NewScanner(fileL)
	for scannerLs.Scan() {
		lslineText = append(lslineText, scannerLs.Text())
		//lineText := ConvertToString(scannerLs.Text(),"GBK","UTF-8")
		//lslineText = append(lslineText, ConvertToString(scannerLs.Text(), "GBK", "UTF-8"))
	}
	if strings.Contains(hsfile, "_92.") {
		//暂时存在问题，后续进行比对
		cmp92File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "_94.") {
		//暂时存在问题，后续进行比对
		cmp94File(hslineText, lslineText, disCode)
	} else if strings.Contains(hsfile, "_T1.") {
		//暂时存在问题，后续进行比对
		//cmpT1File(hslineText, lslineText, disCode)
	}

	return true
}

func cmpKJQSFile(filePath string) {
	start := time.Now()

	logging.Info("开始比对会计清算数据比对")
	cmpSingelKJQSFile(filePath)
	//fmt.Println("第", k+1, "个销售商文件比对开始")
	cost := time.Since(start)
	fmt.Println("总计用时:", cost)
}

//比对单个销售商下所有文件内容
func cmpSingelKJQSFile(filePath string) {
	//恒生文件夹路径
	var hsFilePath string = filePath + "/hs/output/KJQS/20220309/"
	//融先文件夹路径
	var lsFilePath string = filePath + "/ls/output/KJQS/"

	//客服文件夹下数据
	var hsFileList []string
	var lsFileList []string
	var cmpFileList []string

	hsFileList = getFileList(hsFilePath)
	lsFileList = getFileList(lsFilePath)

	//fmt.Println("恒生文件列表", hsFileList)
	//fmt.Println("融先文件列表", lsFileList)
	cmpFileList = cmpDisFileNum(hsFileList, lsFileList)

	//开始比对单个文件
	for _, v := range cmpFileList {
		logging.Info("开始核对", v, "文件")
		cmpKJQSSgFile(hsFilePath+v, lsFilePath+v, "KJQS")
		logging.Info(v, "文件核对完成")
	}

}

func cmpKJQSSgFile(hsfile string, lsfile string, disCode string) bool {
	//fmt.Println(hsfile)
	//fmt.Println(lsfile)
	var hslineText []string
	var lslineText []string
	//var result bool
	fileH, err := os.Open(hsfile)
	if err != nil {
		logging.Fatal(err)
	}

	scannerHs := bufio.NewScanner(fileH)
	for scannerHs.Scan() {
		hslineText = append(hslineText, scannerHs.Text())
		//gbk读取又乱码，需要转换
		//lineText, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(scannerHs.Text())), simplifiedchinese.GBK.NewEncoder()))
		//hslineText = append(hslineText, ConvertToString(scannerHs.Text(), "GBK", "UTF-8"))
	}

	fileL, err := os.Open(lsfile)
	if err != nil {
		logging.Fatal(err)
	}
	scannerLs := bufio.NewScanner(fileL)
	for scannerLs.Scan() {
		lslineText = append(lslineText, scannerLs.Text())
		//lineText := ConvertToString(scannerLs.Text(),"GBK","UTF-8")
		//lslineText = append(lslineText, ConvertToString(scannerLs.Text(), "GBK", "UTF-8"))
	}
	if strings.Contains(hsfile, "AdjustShare") {
		//暂时存在问题，后续进行比对
		cmpAdjustShareFile(hslineText, lslineText, disCode)
	}

	return true
}
