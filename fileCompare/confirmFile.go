package filecompare

import (
	"fundFileCmp/logging"
	"strconv"
	"strings"
)

//全量文件比对，有不一样就提示出来
func AllCmpFile(hsfile []string, lsfile []string, disCode string, fileType string) bool {
	//比对文件头
	var hsNum, lsNum int
	var num int

	hsNum = len(hsfile)
	lsNum = len(lsfile)
	if hsNum == lsNum {
		num = hsNum
	} else {
		logging.CmpInfo(disCode, fileType, "文件行数不一致")
		logging.CmpInfo(disCode, "恒生行数:", hsNum)
		logging.CmpInfo(disCode, "融先内容:", lsNum)
		//	return false
	}

	for i := 0; i < num; i++ {
		if hsfile[i] != lsfile[i] {
			logging.CmpInfo(disCode, fileType, "文件内容不一致")
			logging.CmpInfo(disCode, "请雅琴姐查看第", i, "行内容")
			logging.CmpInfo(disCode, "恒生内容:", hsfile[i])
			logging.CmpInfo(disCode, "融先内容:", lsfile[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	return true
}

//02文件比对
func cmp02File(hs02file []string, ls02file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hs02file) < 36 || len(ls02file) < 36 {
		logging.CmpInfo(disCode, "销售商02文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hs02file))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls02file))
		return false
	}
	for i := 0; i < 37; i++ {
		if hs02file[i] != ls02file[i] {
			logging.CmpInfo(disCode, "文件头02内容不一致")
			logging.CmpInfo(disCode, "请雅琴姐查看第", i, "行内容")
			logging.CmpInfo(disCode, "恒生内容:", hs02file[i])
			logging.CmpInfo(disCode, "融先内容:", ls02file[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	//37行为文件内容数量
	hsNum, _ = strconv.Atoi(hs02file[37])
	lsNum, _ = strconv.Atoi(ls02file[37])
	logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
	logging.CmpInfo(disCode, "融先导出数量:", lsNum)

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "02文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
	}
	hs02Content := hs02file[38 : 38+hsNum]
	ls02Content := ls02file[38 : 38+lsNum]
	//比较文件内容
	//logging.CmpInfo(disCode, "开始比对02文件")
	//defer logging.CmpInfo(disCode, "02文件比对完毕")
	result = cmp02FileContent(hs02Content, ls02Content, disCode)

	return result
}

func cmp02FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用申请单号作为key
	for _, v := range hsContent {
		hsConMap[v[0:24]] = v
	}
	//加工融先导出文件内容，使用申请单号作为key
	for _, v := range lsContent {
		lsConMap[v[0:24]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "02文件存在问题,请核对")
			}
			if lsConMap[k] == "" {
				if result {
					result = false
					logging.CmpInfo(disCode, "02文件存在问题,请核对")
				}
				logging.CmpInfo(disCode, "02文件数据缺失")
				logging.CmpInfo(disCode, "恒生份额", v)
				logging.CmpInfo(disCode, "融先系统导出02文件无法找到匹配数据")
				continue
			}
			//交易确认日期
			if v[24:24+8] != lsConMap[k][24:24+8] {
				logging.CmpInfo(disCode, "02文件申请单号:", k, "交易确认日期不一致")
				logging.CmpInfo(disCode, "恒生数据", v[24:24+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][24:24+8])
			}
			if v[32:32+4] != lsConMap[k][32:32+4] {
				logging.CmpInfo(disCode, "02文件申请单号:", k, "返回代码不一致")
				logging.CmpInfo(disCode, "恒生数据", v[32:32+4])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][32:32+4])
			}
			//基金账号暂不稽核，需要稽核后修改比对条件v[38:38+17+9+3+12] != lsConMap[k][38:38+17+9+3+12]
			if v[36:36+17+9+3+12] != lsConMap[k][36:36+17+9+3+12] {
				logging.CmpInfo(disCode, "02文件申请单号:", k, "基础信息不一致(交易账号,销售商代码,业务代码,基金账号)")
				logging.CmpInfo(disCode, "恒生数据", v[36:36+17+9+3+12])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][36:36+17+9+3+12])
			}
			if v[77:77+1] != lsConMap[k][77:77+1] {
				logging.CmpInfo(disCode, "02文件申请单号:", k, "开户渠道标识不一致")
				logging.CmpInfo(disCode, "恒生数据", v[77:77+1])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][77:77+1])
			}
			//确认单号暂不稽核
			if v[98:98+8+6] != lsConMap[k][98:98+8+6] {
				logging.CmpInfo(disCode, "02文件申请单号:", k, "确认日期或确认时间不一致")
				logging.CmpInfo(disCode, "恒生数据", v[98:98+8+6])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][98:98+8+6])
			}
			if v[112:112+9+1] != lsConMap[k][112:112+9+1] {
				logging.CmpInfo(disCode, "02文件申请单号:", k, "网点号或TA发起标识不一致")
				logging.CmpInfo(disCode, "恒生数据", v[112:112+9+1])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][112:112+9+1])
			}
			//if v[122:122+1+30+120+1+12+8+4+17] != lsConMap[k][122:122+1+30+120+1+12+8+4+17]
			if v[122:122+1+30+120+1+12+8+4+17] != lsConMap[k][122:122+1+30+120+1+12+8+4+17] {
				logging.CmpInfo(disCode, "02文件申请单号:", k, "客户基础信息不一致(证件类型,证件号码,名称,客户类型,简称,凭证号,地区号,对方销售商交易账号)")
				logging.CmpInfo(disCode, "恒生数据", v[122:122+1+30+120+1+12+8+4+17])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][122:122+1+30+120+1+12+8+4+17])
			}

			if v[315:] != lsConMap[k][315:] {
				logging.CmpInfo(disCode, "02文件申请单号:", k, "其他信息不一致(操作网点,摘要说明,客户编号,冻结原因,冻结截止日期,出错详细信息)")
				logging.CmpInfo(disCode, "恒生数据", v[315:])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][315:])
			}

		}
	}

	//比对恒生缺失份额数据
	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "02文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "02文件数据缺失")
			logging.CmpInfo(disCode, "融先份额", v)
			logging.CmpInfo(disCode, "恒生系统导出02文件无法找到匹配数据")
			continue
		}
	}
	if !result {
		logging.CmpInfo(disCode, "02文件核对完毕")
	}

	return result
}

//04文件比对
//todo：存在问题，需要确认最终字段数调整
func cmp04File(hs04file []string, ls04file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hs04file) < 127 || len(ls04file) < 127 {
		logging.CmpInfo(disCode, "销售商04文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hs04file))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls04file))
		return false
	}
	for i := 0; i < 126; i++ {
		if hs04file[i] != ls04file[i] {
			logging.CmpInfo(disCode, "04文件头内容不一致")
			logging.CmpInfo(disCode, "请雅琴姐查看第", i, "行内容")
			logging.CmpInfo(disCode, "恒生内容:", hs04file[i])
			logging.CmpInfo(disCode, "融先内容:", ls04file[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	//126行为文件内容数量
	hsNum, _ = strconv.Atoi(hs04file[131])
	lsNum, _ = strconv.Atoi(ls04file[131])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "04文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
	}
	hs04Content := hs04file[132 : 132+hsNum]
	ls04Content := ls04file[132 : 132+lsNum]
	//比较文件内容
	//logging.CmpInfo(disCode, "开始比对02文件")
	//defer logging.CmpInfo(disCode, "02文件比对完毕")
	result = cmp04FileContent(hs04Content, ls04Content, disCode)

	return result
}

func cmp04FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用申请单号作为key
	for _, v := range hsContent {
		hsConMap[v[0:24]] = v
	}
	//加工融先导出文件内容，使用申请单号作为key
	for _, v := range lsContent {
		lsConMap[v[0:24]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		if v == lsConMap[k] {
			logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "04文件存在问题,请核对")
			}
			if lsConMap[k] == "" {
				if result {
					result = false
					logging.CmpInfo(disCode, "04文件存在问题,请核对")
				}
				logging.CmpInfo(disCode, "04文件数据缺失")
				logging.CmpInfo(disCode, "恒生份额", v)
				logging.CmpInfo(disCode, "融先系统导出04文件无法找到匹配数据")
				continue
			}
			//交易确认日期
			if v[24:24+8+3] != lsConMap[k][24:24+8+3] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "交易确认日期或结算币种不一致")
				logging.CmpInfo(disCode, "恒生数据", v[24:24+8+3])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][24:24+8+3])
			}
			if v[35:35+16+16] != lsConMap[k][35:35+16+16] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "确认金额/份额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[35:35+16+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][35:35+16+16])
			}
			//交易时间不进行稽核
			if v[67:67+6+1+8] != lsConMap[k][67:67+6+1+8] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "基金信息/巨额赎回标志/交易发生日期不一致")
				logging.CmpInfo(disCode, "恒生数据", v[67:67+6+1+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][67:67+6+1+8])
			}
			if v[88:88+4] != lsConMap[k][88:88+4] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "返回代码不一致")
				logging.CmpInfo(disCode, "恒生数据", v[88:88+4])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][88:88+4])
			}
			if v[92:92+17+9] != lsConMap[k][92:92+17+9] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "交易账号或销售商代码不一致")
				logging.CmpInfo(disCode, "恒生数据", v[92:92+17+9])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][92:92+17+9])
			}
			if v[118:118+16+16+3+12] != lsConMap[k][118:118+16+16+3+12] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "申请金额/份额,业务代码,基金账号不一致")
				logging.CmpInfo(disCode, "恒生数据", v[118:118+16+16+3+12])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][118:118+16+16+3+12])
			}
			/*
			if v[185:185+1+5+19+4+8] != lsConMap[k][185:185+1+5+19+4+8] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "杂七杂八的数据1不一致(完全处理标识,折扣率,资金账号,交易地区号,下发日期)")
				logging.CmpInfo(disCode, "恒生数据", v[185:185+1+5+19+4+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][185:185+1+5+19+4+8])
			}
			*/

			if v[222:222+10] != lsConMap[k][222:222+10] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "手续费不一致")
				logging.CmpInfo(disCode, "恒生数据", v[222:222+10])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][222:222+10])
			}
			if v[232:232+10] != lsConMap[k][232:232+10] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "代理费不一致")
				logging.CmpInfo(disCode, "恒生数据", v[232:232+10])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][232:232+10])
			}
			if v[242:242+7] != lsConMap[k][242:242+7] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "单位净值不一致")
				logging.CmpInfo(disCode, "恒生数据", v[242:242+7])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][242:242+7])
			}
			/*
			if v[315:] != lsConMap[k][315:] {
				logging.CmpInfo(disCode, "04文件申请单号:", k, "其他信息不一致(操作网点,摘要说明,客户编号,冻结原因,冻结截止日期,出错详细信息)")
				logging.CmpInfo(disCode, "恒生数据", v[315:])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][315:])
			}
			*/

		}
	}
	//比对恒生缺失份额数据
	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "04文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "04文件数据缺失")
			logging.CmpInfo(disCode, "融先份额", v)
			logging.CmpInfo(disCode, "恒生系统导出04文件无法找到匹配数据")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "04文件核对完毕")
	}

	return true
}

//05文件比对
func cmp05File(hs05file []string, ls05file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hs05file) < 32 || len(ls05file) < 32 {
		logging.CmpInfo(disCode, "销售商05文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hs05file))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls05file))
		return false
	}
	for i := 0; i < 30; i++ {
		if hs05file[i] != ls05file[i] {
			logging.CmpInfo(disCode, "05文件头内容不一致")
			logging.CmpInfo(disCode, "请雅琴姐查看第", i, "行内容")
			logging.CmpInfo(disCode, "恒生内容:", hs05file[i])
			logging.CmpInfo(disCode, "融先内容:", ls05file[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	hsNum, _ = strconv.Atoi(hs05file[30])
	lsNum, _ = strconv.Atoi(ls05file[30])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "05文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
	}
	if hsNum == 0 && lsNum == 0 {
		return true
	}
	if hsNum == 0 || lsNum == 0 {
		logging.CmpInfo(disCode, "存在一方05文件为空,不再进行比对")
		return false
	}

	hs05Content := hs05file[31 : 31+hsNum]
	ls05Content := ls05file[31 : 31+lsNum]
	//fmt.Println(hs05Content)
	//fmt.Println(ls05Content)
	//比较文件内容
	//logging.CmpInfo(disCode, "开始比对02文件")
	//defer logging.CmpInfo(disCode, "02文件比对完毕")
	result = cmp05FileContent(hs05Content, ls05Content, disCode)

	return result
}

func cmp05FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用剩余份额+确认日期+六要素+数据明细标识作为key
	for _, v := range hsContent {
		hsConMap[v[0:109]+v[146:147]] = v
	}
	//加工融先导出文件内容，使用剩余份额+确认日期+六要素+数据明细标识作为key
	for _, v := range lsContent {
		lsConMap[v[0:109]+v[146:147]] = v
	}

	//根据加工后key进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		//比对融先缺失数据
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "05文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "05文件数据缺失")
			logging.CmpInfo(disCode, "恒生份额", v)
			logging.CmpInfo(disCode, "融先系统导出05文件无法找到匹配数据")
			continue
		}

		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "05文件存在问题,请核对")
			}
			//05文件整体比较，暂不做字段划分
			logging.CmpInfo(disCode, "05以下导出数据存在问题")
			logging.CmpInfo(disCode, "恒生内容:", v)
			logging.CmpInfo(disCode, "融先内容:", lsConMap[k])

		}
	}

	//比对恒生缺失份额数据
	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "05文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "05文件数据缺失")
			logging.CmpInfo(disCode, "融先份额", v)
			logging.CmpInfo(disCode, "恒生系统导出05文件无法找到匹配赎回")
			continue
		}
	}
	if !result {
		logging.CmpInfo(disCode, "05文件核对完毕")
	}

	return result
}

//10文件比对
func cmp10File(hs10file []string, ls10file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hs10file) < 70 || len(ls10file) < 70 {
		logging.CmpInfo(disCode, "销售商10文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hs10file))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls10file))
		return false
	}
	for i := 0; i < 68; i++ {
		if hs10file[i] != ls10file[i] {
			logging.CmpInfo(disCode, "文件头10内容不一致")
			logging.CmpInfo(disCode, "请雅琴姐查看第", i, "行内容")
			logging.CmpInfo(disCode, "恒生内容:", hs10file[i])
			logging.CmpInfo(disCode, "融先内容:", ls10file[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	//37行为文件内容数量
	hsNum, _ = strconv.Atoi(hs10file[68])
	lsNum, _ = strconv.Atoi(ls10file[68])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "10文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
	}
	hs10Content := hs10file[69 : 69+hsNum]
	ls10Content := ls10file[69 : 69+lsNum]
	//比较文件内容
	//logging.CmpInfo(disCode, "开始比对02文件")
	//defer logging.CmpInfo(disCode, "02文件比对完毕")
	result = cmp10FileContent(hs10Content, ls10Content, disCode)

	return result
}

func cmp10FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用基金代码作为key
	for _, v := range hsContent {
		hsConMap[v[48:48+6]] = v
	}
	//加工融先导出文件内容，使用基金代码作为key
	for _, v := range lsContent {
		lsConMap[v[48:48+6]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "10文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "10文件数据缺失")
			logging.CmpInfo(disCode, "恒生10文件数据", v)
			logging.CmpInfo(disCode, "融先系统导出10文件无法找到匹配赎回")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "10文件存在问题,请核对")
			}

			//交易确认日期
			if v[:16+16+16] != lsConMap[k][:16+16+16] {
				logging.CmpInfo(disCode, "10文件基金代码:", k, "申购金额/费用/份额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[:16+16+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][:16+16+16])
			}
			if v[54:54+16+16+16] != lsConMap[k][54:54+16+16+16] {
				logging.CmpInfo(disCode, "10文件基金代码:", k, "返回代码不一致")
				logging.CmpInfo(disCode, "恒生数据", v[54:54+16+16+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][54:54+16+16+16])
			}
			//检查销售商代码
			if strings.Replace(v[102:102+9], " ", "", -1) != disCode || strings.Replace(lsConMap[k][102:102+9], " ", "", -1) != disCode {
				logging.CmpInfo(disCode, "10文件基金代码:", k, "销售商代码错误")
				logging.CmpInfo(disCode, "恒生销售商代码", strings.Replace(v[102:102+9], " ", "", -1))
				logging.CmpInfo(disCode, "融先销售商代码", strings.Replace(lsConMap[k][102:102+9], " ", "", -1))
			}

			if v[111:111+16] != lsConMap[k][111:111+16] {
				logging.CmpInfo(disCode, "10文件基金代码:", k, "退款金额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[111:111+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][111:111+16])
			}
			if v[127:127+16] != lsConMap[k][127:127+16] {
				logging.CmpInfo(disCode, "10文件基金代码:", k, "基金总份数不一致")
				logging.CmpInfo(disCode, "恒生数据", v[127:127+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][127:127+16])
			}
			if v[143:143+16+16+16+10] != lsConMap[k][143:143+16+16+16+10] {
				logging.CmpInfo(disCode, "10文件基金代码:", k, "分红不一致")
				logging.CmpInfo(disCode, "恒生数据", v[143:143+16+16+16+10])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][143:143+16+16+16+10])
			}
			if v[201:201+16] != lsConMap[k][201:201+16] {
				logging.CmpInfo(disCode, "10文件基金代码:", k, "划回资金不一致")
				logging.CmpInfo(disCode, "恒生数据", v[201:201+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][201:201+16])
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "10文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "10文件数据缺失,缺失基金代码", k)
			logging.CmpInfo(disCode, "融先10文件数据", v)
			logging.CmpInfo(disCode, "恒生系统导出10文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "10文件核对完毕")
	}

	return result
}

//11文件比对
func cmp11File(hs11file []string, ls11file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hs11file) < 18 || len(ls11file) < 18 {
		logging.CmpInfo(disCode, "销售商11文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hs11file))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls11file))
		return false
	}
	for i := 0; i < 16; i++ {
		if hs11file[i] != ls11file[i] {
			logging.CmpInfo(disCode, "文件头11内容不一致")
			logging.CmpInfo(disCode, "请雅琴姐查看第", i, "行内容")
			logging.CmpInfo(disCode, "恒生内容:", hs11file[i])
			logging.CmpInfo(disCode, "融先内容:", ls11file[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	//37行为文件内容数量
	hsNum, _ = strconv.Atoi(hs11file[17])
	lsNum, _ = strconv.Atoi(ls11file[17])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "11文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
	}
	hs11Content := hs11file[18 : 18+hsNum]
	ls11Content := ls11file[18 : 18+lsNum]
	result = cmp11FileContent(hs11Content, ls11Content, disCode)

	return result
}

func cmp11FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用基金代码+销售商代码+业务代码作为key
	for _, v := range hsContent {
		hsConMap[v[:18]] = v
	}
	//加工融先导出文件内容，使用基金代码作为key
	for _, v := range lsContent {
		lsConMap[v[:18]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "11文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "11文件数据缺失")
			logging.CmpInfo(disCode, "恒生11文件数据", v)
			logging.CmpInfo(disCode, "融先系统导出11文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "11文件存在问题,请核对")
			}

			//交易笔数
			if v[18:18+8] != lsConMap[k][18:18+8] {
				logging.CmpInfo(disCode, "11文件信息:", k, "交易笔数不一致")
				logging.CmpInfo(disCode, "恒生数据", v[18:18+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][18:18+8])
			}
			if v[26:26+8] != lsConMap[k][26:26+8] {
				logging.CmpInfo(disCode, "11文件信息:", k, "日期不一致")
				logging.CmpInfo(disCode, "恒生数据", v[26:26+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][26:26+8])
			}
			if v[34:34+16] != lsConMap[k][34:34+16] {
				logging.CmpInfo(disCode, "11文件信息:", k, "确认份数不一致")
				logging.CmpInfo(disCode, "恒生数据", v[34:34+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][34:34+16])
			}
			if v[50:50+16] != lsConMap[k][50:50+16] {
				logging.CmpInfo(disCode, "11文件信息:", k, "确认金额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[50:50+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][50:50+16])
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "11文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "11文件数据缺失,缺失基金代码", k)
			logging.CmpInfo(disCode, "融先11文件数据", v)
			logging.CmpInfo(disCode, "恒生系统导出11文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "11文件核对完毕")
	}

	return result
}

//12文件比对
func cmp12File(hs12file []string, ls12file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hs12file) < 22 || len(ls12file) < 22 {
		logging.CmpInfo(disCode, "销售商12文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hs12file))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls12file))
		return false
	}
	for i := 0; i < 20; i++ {
		if hs12file[i] != ls12file[i] {
			logging.CmpInfo(disCode, "文件头12内容不一致")
			logging.CmpInfo(disCode, "请雅琴姐查看第", i, "行内容")
			logging.CmpInfo(disCode, "恒生内容:", hs12file[i])
			logging.CmpInfo(disCode, "融先内容:", ls12file[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	//37行为文件内容数量
	hsNum, _ = strconv.Atoi(hs12file[21])
	lsNum, _ = strconv.Atoi(ls12file[21])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "12文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
	}
	hs12Content := hs12file[22 : 22+hsNum]
	ls12Content := ls12file[22 : 22+lsNum]
	result = cmp12FileContent(hs12Content, ls12Content, disCode)

	return result
}

func cmp12FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用基金代码+销售商代码+业务代码作为key
	for _, v := range hsContent {
		hsConMap[v[:18]] = v
	}
	//加工融先导出文件内容，使用基金代码作为key
	for _, v := range lsContent {
		lsConMap[v[:18]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "12文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "12文件数据缺失")
			logging.CmpInfo(disCode, "恒生12文件数据", v)
			logging.CmpInfo(disCode, "融先系统导出12文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "12文件存在问题,请核对")
			}

			//交易笔数
			if v[18:18+8] != lsConMap[k][18:18+8] {
				logging.CmpInfo(disCode, "12文件信息:", k, "业务笔数不一致")
				logging.CmpInfo(disCode, "恒生数据", v[18:18+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][18:18+8])
			}
			if v[26:26+8] != lsConMap[k][26:26+8] {
				logging.CmpInfo(disCode, "12文件信息:", k, "日期不一致")
				logging.CmpInfo(disCode, "恒生数据", v[26:26+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][26:26+8])
			}
			if v[34:34+16] != lsConMap[k][34:34+16] {
				logging.CmpInfo(disCode, "12文件信息:", k, "确认失败份数不一致")
				logging.CmpInfo(disCode, "恒生数据", v[34:34+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][34:34+16])
			}
			if v[50:50+16] != lsConMap[k][50:50+16] {
				logging.CmpInfo(disCode, "12文件信息:", k, "确认成功份数不一致")
				logging.CmpInfo(disCode, "恒生数据", v[50:50+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][50:50+16])
			}
			if v[66:66+16] != lsConMap[k][66:66+16] {
				logging.CmpInfo(disCode, "12文件信息:", k, "确认失败金额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[66:66+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][66:66+16])
			}
			if v[82:82+16] != lsConMap[k][82:82+16] {
				logging.CmpInfo(disCode, "12文件信息:", k, "确认成功金额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[82:82+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][82:82+16])
			}
			if v[98:98+6] != lsConMap[k][98:98+6] {
				logging.CmpInfo(disCode, "12文件信息:", k, "失败笔数不一致")
				logging.CmpInfo(disCode, "恒生数据", v[98:98+6])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][98:98+6])
			}
			if v[104:104+6] != lsConMap[k][104:104+6] {
				logging.CmpInfo(disCode, "12文件信息:", k, "成功笔数不一致")
				logging.CmpInfo(disCode, "恒生数据", v[104:104+6])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][104:104+6])
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "12文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "12文件数据缺失,缺失基金代码", k)
			logging.CmpInfo(disCode, "融先12文件数据", v)
			logging.CmpInfo(disCode, "恒生系统导出12文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "12文件核对完毕")
	}

	return result
}

//25文件比对
func cmp25File(hs25file []string, ls25file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hs25file) < 29 || len(ls25file) < 29 {
		logging.CmpInfo(disCode, "销售商25文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hs25file))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls25file))
		return false
	}
	for i := 0; i < 27; i++ {
		if hs25file[i] != ls25file[i] {
			logging.CmpInfo(disCode, "文件头25内容不一致")
			logging.CmpInfo(disCode, "请雅琴姐查看第", i, "行内容")
			logging.CmpInfo(disCode, "恒生内容:", hs25file[i])
			logging.CmpInfo(disCode, "融先内容:", ls25file[i])
			return false
		}
		//fmt.Println(hs02file[i])
	}
	//37行为文件内容数量
	hsNum, _ = strconv.Atoi(hs25file[28])
	lsNum, _ = strconv.Atoi(ls25file[28])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "25文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
	}
	hs25Content := hs25file[29 : 29+hsNum]
	ls25Content := ls25file[29 : 29+lsNum]
	result = cmp25FileContent(hs25Content, ls25Content, disCode)

	return result
}

func cmp25FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用发生日期+基金代码+业务代码+资金类型作为key
	for _, v := range hsContent {
		hsConMap[v[20:20+8]+v[70:70+6]+v[100:100+3+3]] = v
	}
	//加工融先导出文件内容，使用基金代码作为key
	for _, v := range lsContent {
		lsConMap[v[20:20+8]+v[70:70+6]+v[100:100+3+3]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "25文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "25文件数据缺失")
			logging.CmpInfo(disCode, "恒生25文件数据", v)
			logging.CmpInfo(disCode, "融先系统导出25文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "25文件存在问题,请核对")
			}

			if v[28:28+8+3+1+30] != lsConMap[k][28:28+8+3+1+30] {
				logging.CmpInfo(disCode, "25文件信息:", k, "确认日期,币种,交易所标识,资金账号不一致")
				logging.CmpInfo(disCode, "恒生数据", v[28:28+8+3+1+30])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][28:28+8+3+1+30])
			}
			if v[76:76+9+6+9] != lsConMap[k][76:76+9+6+9] {
				logging.CmpInfo(disCode, "12文件信息:", k, "销售商代码,席位代码,网点代码不一致")
				logging.CmpInfo(disCode, "恒生数据", v[76:76+9+6+9])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][76:76+9+6+9])
			}
			if v[106:106+16] != lsConMap[k][106:106+16] {
				logging.CmpInfo(disCode, "25文件信息:", k, "确认金额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[106:106+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][106:106+16])
			}
			if v[122:122+1] != lsConMap[k][122:122+1] {
				logging.CmpInfo(disCode, "25文件信息:", k, "收付标志不一致")
				logging.CmpInfo(disCode, "恒生数据", v[122:122+1])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][122:122+1])
			}
			if v[123:] != lsConMap[k][123:] {
				logging.CmpInfo(disCode, "25文件信息:", k, "清算/交收日期,明细标识不一致")
				logging.CmpInfo(disCode, "恒生数据", v[123:])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][123:])
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "25文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "25文件数据缺失,缺失基金代码", k)
			logging.CmpInfo(disCode, "融先25文件数据", v)
			logging.CmpInfo(disCode, "恒生系统导出25文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "25文件核对完毕")
	}

	return result
}
