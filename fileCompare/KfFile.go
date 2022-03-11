package filecompare

import (
	"fundFileCmp/logging"
	"strconv"
)

//Acconet文件比对
func cmpAcconetFile(hsfile []string, ls21file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hsfile) < 2 || len(ls21file) < 2 {
		logging.CmpInfo(disCode, "销售商acconet文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hsfile))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls21file))
		return false
	}

	//第1行为文件内容数量
	hsNum, _ = strconv.Atoi(hsfile[0])
	lsNum, _ = strconv.Atoi(ls21file[0])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "acconet文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
		//return false
	}
	hsContent := hsfile[1 : 1+hsNum]
	lsContent := ls21file[1 : 1+lsNum]
	result = cmpAcconetFileContent(hsContent, lsContent, disCode)

	return result
}

func cmpAcconetFileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用基金账号，销售商代码，网点代码，交易账号作为key
	for _, v := range hsContent {
		hsConMap[v[0:12+9+9]+v[42:42+17]] = v
	}
	//加工融先导出文件内容，使用基金账号，销售商代码，网点代码，交易账号作为key
	for _, v := range lsContent {
		lsConMap[v[0:12+9+9]+v[42:42+17]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "Acconet文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "Acconet文件数据缺失", k)
			logging.CmpInfo(disCode, "恒生Acconet文件数据", v)
			logging.CmpInfo(disCode, "融先系统导出Acconet文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "Acconet文件存在问题,请核对")
			}

			if v[59:59+1] != lsConMap[k][59:59+1] {
				logging.CmpInfo(disCode, "Acconet文件信息:", k, "网点标识不一致")
				logging.CmpInfo(disCode, "恒生数据", v[59:59+1])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][59:59+1])
			}
			if v[60:60+1] != lsConMap[k][60:60+1] {
				logging.CmpInfo(disCode, "Acconet文件信息:", k, "默认分红方式不一致")
				logging.CmpInfo(disCode, "恒生数据", v[15:15+80])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][15:15+80])
			}
			if v[61:61+200+40+20] != lsConMap[k][61:61+200+40+20] {
				logging.CmpInfo(disCode, "Acconet文件信息:", k, "收款银行相关信息不一致")
				logging.CmpInfo(disCode, "恒生数据", v[61:61+200+40+20])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][61:61+200+40+20])
			}
			if v[321:] != lsConMap[k][321:] {
				logging.CmpInfo(disCode, "Acconet文件信息:", k, "其他信息不一致")
				logging.CmpInfo(disCode, "恒生数据", v[321:])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][321:])
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "Acconet文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "Acconet文件数据缺失", k)
			logging.CmpInfo(disCode, "融先Acconet文件数据", v)
			logging.CmpInfo(disCode, "恒生系统导出Acconet文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "Acconet文件核对完毕")
	}

	return result
}

//Accorequest文件比对
func cmpAccoReqFile(hsfile []string, ls21file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hsfile) < 2 || len(ls21file) < 2 {
		logging.CmpInfo(disCode, "销售商Accorequest文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hsfile))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls21file))
		return false
	}

	//第1行为文件内容数量
	hsNum, _ = strconv.Atoi(hsfile[0])
	lsNum, _ = strconv.Atoi(ls21file[0])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "Accorequest文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
		//return false
	}
	hsContent := hsfile[1 : 1+hsNum]
	lsContent := ls21file[1 : 1+lsNum]
	result = cmpAccorequestFileContent(hsContent, lsContent, disCode)

	return result
}

func cmpAccorequestFileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用申请单号作为key
	for _, v := range hsContent {
		hsConMap[v[17:17+24+9]] = v
	}
	//加工融先导出文件内容，使用申请单号作为key
	for _, v := range lsContent {
		lsConMap[v[17:17+24+9]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "Accorequest文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "Accorequest文件数据缺失", k)
			logging.CmpInfo(disCode, "恒生Accorequest文件数据", v)
			logging.CmpInfo(disCode, "融先系统导出Accorequest文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "Accorequest文件存在问题,请核对")
			}

			if v[:3] != lsConMap[k][:3] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "业务代码1不一致")
				logging.CmpInfo(disCode, "恒生数据", v[:3])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][:3])
			}
			if v[3:3+8] != lsConMap[k][3:3+8] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "交易申请日期不一致")
				logging.CmpInfo(disCode, "恒生数据", v[3:3+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][3:3+8])
			}
			if v[44:44+4] != lsConMap[k][44:44+4] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "地区编号不一致")
				logging.CmpInfo(disCode, "恒生数据", v[44:44+4])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][44:44+4])
			}
			if v[48:48+9] != lsConMap[k][48:48+9] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "网点号码不一致")
				logging.CmpInfo(disCode, "恒生数据", v[48:48+9])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][48:48+9])
			}
			if v[77:77+200+20+1+3+40] != lsConMap[k][77:77+200+20+1+3+40] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "投资人基本信息(名称,证件号码等)不一致")
				logging.CmpInfo(disCode, "恒生数据", v[77:77+200+20+1+3+40])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][77:77+200+20+1+3+40])
			}
			if v[341:341+60+60+3+40] != lsConMap[k][341:341+60+60+3+40] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "法人/经办人信息不一致")
				logging.CmpInfo(disCode, "恒生数据", v[341:341+60+60+3+40])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][341:341+60+60+3+40])
			}
			if v[504:504+300+6+40+40+40] != lsConMap[k][504:504+300+6+40+40+40] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "基本信息1不一致")
				logging.CmpInfo(disCode, "恒生数据", v[504:504+300+6+40+40+40])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][504:504+300+6+40+40+40])
			}
			if v[930:930+40+8+1+3+5+16] != lsConMap[k][930:930+40+8+1+3+5+16] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "基本信息2不一致")
				logging.CmpInfo(disCode, "恒生数据", v[930:930+40+8+1+3+5+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][930:930+40+8+1+3+5+16])
			}
			if v[1004:1004+17] != lsConMap[k][1004:1004+17] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "交易账号不一致")
				logging.CmpInfo(disCode, "恒生数据", v[1004:1004+17])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][1004:1004+17])
			}
			if v[1020:1332] != lsConMap[k][1020:1332] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "杂七杂八信息1不一致")
				logging.CmpInfo(disCode, "恒生数据", v[1020:1332])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][1020:1332])
			}
			if v[1332:1332+3+1+9+8] != lsConMap[k][1332:1332+3+1+9+8] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "杂七杂八信息2不一致")
				logging.CmpInfo(disCode, "恒生数据", v[1332:1332+3+1+9+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][1332:1332+3+1+9+8])
			}
			if v[1435:] != lsConMap[k][1435:] {
				logging.CmpInfo(disCode, "Accorequest文件信息:", k, "杂七杂八信息3不一致")
				logging.CmpInfo(disCode, "恒生数据", v[1435:])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][1435:])
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "Acconet文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "Acconet文件数据缺失", k)
			logging.CmpInfo(disCode, "融先Acconet文件数据", v)
			logging.CmpInfo(disCode, "恒生系统导出Acconet文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "21文件核对完毕")
	}

	return result
}

//confirm文件比对
func cmpKFCfmFile(hsfile []string, ls21file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hsfile) < 2 || len(ls21file) < 2 {
		logging.CmpInfo(disCode, "销售商confirm文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hsfile))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls21file))
		return false
	}

	//第1行为文件内容数量
	hsNum, _ = strconv.Atoi(hsfile[0])
	lsNum, _ = strconv.Atoi(ls21file[0])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "confirm文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
		//return false
	}
	hsContent := hsfile[1 : 1+hsNum]
	lsContent := ls21file[1 : 1+lsNum]
	result = cmpCfmFileContent(hsContent, lsContent, disCode)

	return result
}

func cmpCfmFileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用申请单号作为key
	for _, v := range hsContent {
		hsConMap[v[76:76+24]] = v
	}
	//加工融先导出文件内容，使用申请单号作为key
	for _, v := range lsContent {
		lsConMap[v[76:76+24]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "confirm文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "confirm文件数据缺失", k)
			logging.CmpInfo(disCode, "恒生confirm文件数据", v)
			logging.CmpInfo(disCode, "融先系统导出confirm文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "confirm文件存在问题,请核对")
			}

			if v[:8] != lsConMap[k][:8] {
				logging.CmpInfo(disCode, "confirm文件信息:", k, "确认日期不一致")
				logging.CmpInfo(disCode, "恒生数据", v[:8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][:8])
			}
			if v[34:34+3] != lsConMap[k][34:34+3] {
				logging.CmpInfo(disCode, "confirm文件信息:", k, "业务代码不一致")
				logging.CmpInfo(disCode, "恒生数据", v[34:34+3])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][34:34+3])
			}
			if v[49:49+12+6+1+8] != lsConMap[k][49:49+12+6+1+8] {
				logging.CmpInfo(disCode, "confirm文件信息:", k, "基金账号,基金代码,份额类别,申请日期不一致")
				logging.CmpInfo(disCode, "恒生数据", v[49:49+12+6+1+8])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][49:49+12+6+1+8])
			}
			if v[100:100+9+9+9+3] != lsConMap[k][100:100+9+9+9+3] {
				logging.CmpInfo(disCode, "confirm文件信息:", k, "销售商代码,网点代码,币种不一致")
				logging.CmpInfo(disCode, "恒生数据", v[100:100+9+9+9+3])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][100:100+9+9+9+3])
			}
			if v[130:130+16+16] != lsConMap[k][130:130+16+16] {
				logging.CmpInfo(disCode, "confirm文件信息:", k, "确认金额/份额不一致")
				logging.CmpInfo(disCode, "恒生数据", v[130:130+16+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][130:130+16+16])
			}
			if v[242:242+6+1+12+9+9] != lsConMap[k][242:242+6+1+12+9+9] {
				logging.CmpInfo(disCode, "confirm文件信息:", k, "转换相关字段不一致")
				logging.CmpInfo(disCode, "恒生数据", v[242:242+6+1+12+9+9])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][242:242+6+1+12+9+9])
			}
			if v[279:279+4] != lsConMap[k][279:279+4] {
				logging.CmpInfo(disCode, "confirm文件信息:", k, "返回代码不一致")
				logging.CmpInfo(disCode, "恒生数据", v[279:279+4])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][279:279+4])
			}
			/*
				if v[283:283+20+1+1+8+24] != lsConMap[k][283:283+20+1+1+8+24] {
					logging.CmpInfo(disCode, "confirm文件信息:", k, "乱七八糟1字段不一致")
					logging.CmpInfo(disCode, "恒生数据", v[283:283+20+1+1+8+24])
					logging.CmpInfo(disCode, "融先数据", lsConMap[k][283:283+20+1+1+8+24])
				}
			*/
			if v[337:337+17+16] != lsConMap[k][337:337+17+16] {
				logging.CmpInfo(disCode, "confirm文件信息:", k, "基金账号,成交价不一致")
				logging.CmpInfo(disCode, "恒生数据", v[337:337+17+16])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][337:337+17+16])
			}
			if v[386:] != lsConMap[k][386:] {
				logging.CmpInfo(disCode, "confirm文件信息:", k, "其他参数不一致")
				logging.CmpInfo(disCode, "恒生数据", v[386:])
				logging.CmpInfo(disCode, "融先数据", lsConMap[k][386:])
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "confirm文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "confirm文件数据缺失", k)
			logging.CmpInfo(disCode, "融先confirm文件数据", v)
			logging.CmpInfo(disCode, "恒生系统导出confirm文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "Acconet文件核对完毕")
	}

	return result
}
