package filecompare

import (
	"fundFileCmp/logging"
	"strconv"
)

//02文件比对
func cmp02File(hs02file []string, ls02file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hs02file) < 36 || len(ls02file) < 36 {
		logging.Error(disCode, "销售商文件存在问题")
		logging.Error("恒生文件行数:", len(hs02file))
		logging.Error("融先文件行数:", len(ls02file))
		return false
	}
	for i := 0; i < 37; i++ {
		if hs02file[i] != ls02file[i] {
			logging.Error("文件头内容不一致")
			logging.Error("请雅琴姐查看第", i, "内容")
			logging.Error("恒生内容:", hs02file[i])
			logging.Error("融先内容:", ls02file[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	//37行为文件内容数量
	hsNum, _ = strconv.Atoi(hs02file[37])
	lsNum, _ = strconv.Atoi(ls02file[37])

	if hsNum != lsNum {
		logging.Error("02文件恒生导出数量与融先导出数量不一致")
		logging.Error("恒生导出数量:", hsNum)
		logging.Error("融先导出数量:", lsNum)
	}
	hs02Content := hs02file[37+1 : 37+hsNum]
	ls02Content := ls02file[37+1 : 37+lsNum]
	//比较文件内容
	result = cmpFileContent(hs02Content, ls02Content)

	return result
}

func cmpFileContent(hsContent []string, lsContent []string) bool {
	return true
}
