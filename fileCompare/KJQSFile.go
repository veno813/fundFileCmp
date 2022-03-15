package filecompare

import (
	"fundFileCmp/logging"
)

//Acconet文件比对
func cmpAdjustShareFile(hsfile []string, ls21file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	hsNum = len(hsfile)
	lsNum = len(ls21file)
	logging.CmpInfo(disCode, "恒生AdjustShare文件数据量", hsNum)
	logging.CmpInfo(disCode, "融先AdjustShare文件数据", lsNum)

	if hsNum < 0 || lsNum < 0 {
		logging.CmpInfo(disCode, "销售商AdjustShare文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hsfile))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls21file))
		return false
	}

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "AdjustShare文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
		//return false
	}
	hsContent := hsfile[:hsNum]
	lsContent := ls21file[:lsNum]
	result = cmpAdjustShareFileContent(hsContent, lsContent, disCode)

	return result
}

func cmpAdjustShareFileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用基金账号，销售商代码，网点代码，交易账号作为key
	for _, v := range hsContent {
		hsConMap[v[22:22+6+3]+v[163:164]] = v
	}
	//加工融先导出文件内容，使用基金账号，销售商代码，网点代码，交易账号作为key
	for _, v := range lsContent {
		lsConMap[v[22:22+6+3]+v[163:164]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "AdjustShare文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "AdjustShare文件数据缺失  基金：", k[0:6], "  销售商：", k[6:])
			//logging.CmpInfo(disCode, "恒生AdjustShare文件数据", v)
			logging.CmpInfo(disCode, "融先系统导出AdjustShare文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "AdjustShare文件存在问题,请核对")
			}
			if v[131:131+16] != lsConMap[k][131:131+16] {
				logging.CmpInfo(disCode, "Acconet文件信息:", k, "剩余份额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[131:131+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][131:131+16])
			}

		}
	}

	for k, _ := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "AdjustShare文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "AdjustShare文件数据缺失  基金：", k[0:6], "  销售商：", k[6:])
			//logging.CmpInfo(disCode, "融先AdjustShare文件数据", v)
			logging.CmpInfo(disCode, "恒生系统导出AdjustShare文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "AdjustShare文件核对完毕")
	}

	return result
}
