package filecompare

import (
	"fundFileCmp/logging"
	"strconv"
)

//21文件比对
func cmp21File(hs21file []string, ls21file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hs21file) < 16 || len(ls21file) < 16 {
		logging.CmpInfo(disCode, "销售商21文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hs21file))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls21file))
		return false
	}
	for i := 0; i < 14; i++ {
		if hs21file[i] != ls21file[i] {
			logging.CmpInfo(disCode, "文件头21内容不一致")
			logging.CmpInfo(disCode, "请雅琴姐查看第", i, "行内容")
			logging.CmpInfo(disCode, "恒生内容:", hs21file[i])
			logging.CmpInfo(disCode, "融先内容:", ls21file[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	//37行为文件内容数量
	hsNum, _ = strconv.Atoi(hs21file[15])
	lsNum, _ = strconv.Atoi(ls21file[15])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "21文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
		//return false
	}
	hs21Content := hs21file[16 : 16+hsNum]
	ls21Content := ls21file[16 : 16+lsNum]
	result = cmp21FileContent(hs21Content, ls21Content, disCode)

	return result
}

func cmp21FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用销售商代码作为key
	for _, v := range hsContent {
		hsConMap[v[6:6+9]] = v
	}
	//加工融先导出文件内容，使用销售商代码作为key
	for _, v := range lsContent {
		lsConMap[v[6:6+9]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "21文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "21文件数据缺失")
			logging.CmpInfo(disCode, "恒生21文件数据", v)
			logging.CmpInfo(disCode, "融先系统导出21文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "21文件存在问题,请核对")
			}

			if v[:6] != lsConMap[k][:6] {
				logging.CmpInfo(disCode, "21文件信息:", k, "席位代码不一致")
				logging.CmpInfo(disCode, "恒生数据", v[:6])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][:6])
			}
			if v[15:15+80] != lsConMap[k][15:15+80] {
				logging.CmpInfo(disCode, "21文件信息:", k, "销售商名称不一致")
				logging.CmpInfo(disCode, "恒生数据", v[15:15+80])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][15:15+80])
			}
			if v[95:95+1] != lsConMap[k][95:95+1] {
				logging.CmpInfo(disCode, "21文件信息:", k, "市场转入标识不一致")
				logging.CmpInfo(disCode, "恒生数据", v[95:95+1])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][95:95+1])
			}
			if v[96:96+8] != lsConMap[k][96:96+8] {
				logging.CmpInfo(disCode, "21文件信息:", k, "确认金额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[96:96+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][96:96+8])
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "21文件存在问题,请核对")
			}
			//logging.CmpInfo(disCode, "21文件数据缺失,缺失基金代码", k)
			logging.CmpInfo(disCode, "融先21文件数据", v)
			//logging.CmpInfo(disCode, "恒生系统导出21文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "21文件核对完毕")
	}

	return result
}
