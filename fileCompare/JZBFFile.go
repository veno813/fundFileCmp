package filecompare

import (
	"fundFileCmp/logging"
	"strconv"
	"strings"
)

//92件比对
func cmp92File(hsfile []string, ls21file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hsfile) < 77 || len(ls21file) < 98 {
		logging.CmpInfo(disCode, "集中备份92文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hsfile))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls21file))
		return false
	}

	//第1行为文件内容数量
	hsNum, _ = strconv.Atoi(hsfile[76])
	lsNum, _ = strconv.Atoi(ls21file[97])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "集中备份92文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
		//return false
	}
	hsContent := hsfile[77 : 77+hsNum]
	lsContent := ls21file[98 : 98+lsNum]
	result = cmp92FileContent(hsContent, lsContent, disCode)

	return result
}

func cmp92FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用基金账号，销售商代码，网点代码，交易账号作为key
	for _, v := range hsContent {
		hsConMap[v[171:171+24]+v[807:807+9]] = v
	}
	//加工融先导出文件内容，使用基金账号，销售商代码，网点代码，交易账号作为key
	for _, v := range lsContent {
		lsConMap[v[403:403+24]+v[1284:1284+9]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "集中备份92文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "集中备份92文件数据缺失", k)
			logging.CmpInfo(disCode, "恒生92文件数据", ConvertToString(v, "GBK", "UTF-8"))
			logging.CmpInfo(disCode, "融先系统导出92文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "92文件存在问题,请核对")
			}

			if strings.TrimSpace(v[0:120]) != strings.TrimSpace(lsConMap[k][0:300]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "地址不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[0:120]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][0:300]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[120:120+30]) != strings.TrimSpace(lsConMap[k][300:300+40]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "法人证件号码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[120:120+30]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][300:300+40]), "GBK", "UTF-8"))
			}
			/*
				if strings.TrimSpace(v[150:150+1]) != strings.TrimSpace(lsConMap[k][340:340+3]) {
					logging.CmpInfo(disCode, "92文件信息:", k, "法人证件类型不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[150:150+1]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][340:340+3]), "GBK", "UTF-8"))
				}
			*/
			if strings.TrimSpace(v[151:151+20]) != strings.TrimSpace(lsConMap[k][343:343+60]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "法人姓名不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[151:151+20]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][343:343+60]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[325:325+1]) != strings.TrimSpace(lsConMap[k][557:557+1]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "证件类型不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[325:325+1]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][557:557+1]), "GBK", "UTF-8"))
			}
			/*
				if strings.TrimSpace(v[326:326+19]) != strings.TrimSpace(lsConMap[k][560:560+40]) {
					logging.CmpInfo(disCode, "92文件信息:", k, "资金账号不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[326:326+19]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][560:560+40]), "GBK", "UTF-8"))
				}
			*/
			if strings.TrimSpace(v[345:345+4]) != strings.TrimSpace(lsConMap[k][600:600+4]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "地区编码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[345:345+4]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][600:600+4]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[349:349+8]) != strings.TrimSpace(lsConMap[k][604:604+8]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "交易确认日期不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[349:349+8]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][604:604+8]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[359:359+40]) != strings.TrimSpace(lsConMap[k][615:615+40]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "邮箱地址不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[359:359+40]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][615:615+40]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[519:519+120]) != strings.TrimSpace(lsConMap[k][846:846+200]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "投资者姓名不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[519:519+120]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][846:846+200]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[639:639+9]) != strings.TrimSpace(lsConMap[k][1046:1046+9]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "网点号码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[639:639+9]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1046:1046+9]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[690:690+24]) != strings.TrimSpace(lsConMap[k][1115:1115+24]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "原申请单号不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[690:690+24]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1115:1115+24]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[714:714+8+6]) != strings.TrimSpace(lsConMap[k][1139:1139+8+6]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "申请日期/时间不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[714:714+8+6]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1139:1139+8+6]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[728:728+1]) != strings.TrimSpace(lsConMap[k][1153:1153+1]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "客户类型不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[728:728+1]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1153:1153+1]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[729:729+6]) != strings.TrimSpace(lsConMap[k][1154:1154+6]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "邮政编码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[729:729+6]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1154:1154+6]), "GBK", "UTF-8"))
			}
			/*
				if strings.TrimSpace(v[735:735+30+1+20]) != strings.TrimSpace(lsConMap[k][1160:1160+40+3+60]) {
					logging.CmpInfo(disCode, "92文件信息:", k, "经办人信息不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[735:735+30+1+20]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1160:1160+40+3+60]), "GBK", "UTF-8"))
				}
			*/
			if strings.TrimSpace(v[786:786+4]) != strings.TrimSpace(lsConMap[k][1263:1263+4]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "返回代码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[786:786+4]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1263:1263+4]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[790:790+17]) != strings.TrimSpace(lsConMap[k][1267:1267+17]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "交易账号不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[790:790+17]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1267:1267+17]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[807:807+9]) != strings.TrimSpace(lsConMap[k][1284:1284+9]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "销售商代码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[807:807+9]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1284:1284+9]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[816:816+12]) != strings.TrimSpace(lsConMap[k][1293:1293+20]) {
				logging.CmpInfo(disCode, "92文件信息:", k, "投资人姓名简称不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[816:816+12]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][1293:1293+20]), "GBK", "UTF-8"))
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "92文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "92文件数据缺失", k)
			logging.CmpInfo(disCode, "融先92文件数据", ConvertToString(v, "GBK", "UTF-8"))
			logging.CmpInfo(disCode, "恒生系统导出92文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "92文件核对完毕")
	}

	return result
}

//94件比对
func cmp94File(hsfile []string, lsfile []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hsfile) < 75 || len(lsfile) < 119 {
		logging.CmpInfo(disCode, "集中备份94文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hsfile))
		logging.CmpInfo(disCode, "融先文件行数:", len(lsfile))
		return false
	}

	//第1行为文件内容数量
	hsNum, _ = strconv.Atoi(hsfile[74])
	lsNum, _ = strconv.Atoi(lsfile[118])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "集中备份94文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
		//return false
	}
	hsContent := hsfile[75 : 75+hsNum]
	lsContent := lsfile[119 : 119+lsNum]
	result = cmp94FileContent(hsContent, lsContent, disCode)

	return result
}

func cmp94FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用申请单号作为key
	for _, v := range hsContent {
		hsConMap[v[:24]] = v
	}
	//加工融先导出文件内容，使用申请单号作为key
	for _, v := range lsContent {
		lsConMap[v[:24]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "集中备份94文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "集中备份94文件数据缺失", k)
			logging.CmpInfo(disCode, "恒生94文件数据", ConvertToString(v, "GBK", "UTF-8"))
			logging.CmpInfo(disCode, "融先系统导出94文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "94文件存在问题,请核对")
			}
			/*
				if strings.TrimSpace(v[24:24+1]) != strings.TrimSpace(lsConMap[k][24:24+1]) {
					logging.CmpInfo(disCode, "94文件信息:", k, "分红方式不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[24:24+1]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][24:24+1]), "GBK", "UTF-8"))
				}
			*/
			if strings.TrimSpace(v[25:25+5]) != strings.TrimSpace(lsConMap[k][25:25+5]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "佣金折扣率不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[25:25+5]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][25:25+5]), "GBK", "UTF-8"))
			}

			if strings.TrimSpace(v[30:30+4]) != strings.TrimSpace(lsConMap[k][30:30+4]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "地区编号不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[30:30+4]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][30:30+4]), "GBK", "UTF-8"))
			}

			if strings.TrimSpace(v[34:34+8]) != strings.TrimSpace(lsConMap[k][34:34+8]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "确认日期不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[34:34+8]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][34:34+8]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[42:42+6]) != strings.TrimSpace(lsConMap[k][42:42+6]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "证件类型不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[42:42+6]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][42:42+6]), "GBK", "UTF-8"))
			}
			/*
				tmpFeeHs, _ := strconv.Atoi(v[51 : 51+10])
				tmpFeeLs, _ := strconv.Atoi(lsConMap[k][51 : 51+16])
				if tmpFeeHs != tmpFeeLs {
					logging.CmpInfo(disCode, "94文件信息:", k, "手续费不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[51:51+10]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][51:51+16]), "GBK", "UTF-8"))
				}
				tmpFee1Hs, _ := strconv.Atoi(v[61 : 61+10])
				tmpFee1Ls, _ := strconv.Atoi(lsConMap[k][67 : 67+16])
				if tmpFee1Hs != tmpFee1Ls {
					logging.CmpInfo(disCode, "94文件信息:", k, "代理费不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[61:61+10]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][67:67+16]), "GBK", "UTF-8"))
				}

				tmpFeeTolHs, _ := strconv.Atoi(v[71 : 71+10])
				tmpFeeTolLs, _ := strconv.Atoi(lsConMap[k][83 : 83+16])
				if tmpFeeTolHs != tmpFeeTolLs {
					logging.CmpInfo(disCode, "94文件信息:", k, "总费用不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[71:71+10]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][83:83+16]), "GBK", "UTF-8"))
				}
			*/

			if strings.TrimSpace(v[90:90+16+16]) != strings.TrimSpace(lsConMap[k][108:108+16+16]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "确认份额/金额不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[90:90+16+16]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][108:108+16+16]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[122:122+6]) != strings.TrimSpace(lsConMap[k][140:140+6]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "基金代码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[122:122+6]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][140:140+6]), "GBK", "UTF-8"))
			}

			tmpInrlHs, _ := strconv.Atoi(v[128 : 128+10])
			tmpInrlLs, _ := strconv.Atoi(lsConMap[k][146 : 146+16])
			if tmpInrlHs != tmpInrlLs {
				logging.CmpInfo(disCode, "94文件信息:", k, "利息不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[128:128+10]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][146:146+16]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[138:138+1]) != strings.TrimSpace(lsConMap[k][162:162+1]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "巨额赎回标识不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[138:138+1]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][162:162+1]), "GBK", "UTF-8"))
			}
			/*
					if strings.TrimSpace(v[139:139+7]) != strings.TrimSpace(lsConMap[k][163+9:163+16]) {
						logging.CmpInfo(disCode, "94文件信息:", k, "网点号码不一致")
						logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[139:139+7]), "GBK", "UTF-8"))
						logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][163+9:163+16]), "GBK", "UTF-8"))
					}

				if strings.TrimSpace(v[146:146+9]) != strings.TrimSpace(lsConMap[k][179:179+9]) {
					logging.CmpInfo(disCode, "94文件信息:", k, "网点代码不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[146:146+9]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][179:179+9]), "GBK", "UTF-8"))
				}
			*/
			if strings.TrimSpace(v[199:199+8+6]) != strings.TrimSpace(lsConMap[k][232:232+8+6]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "申请日期/时间不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[199:199+8+6]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][232:232+8+6]), "GBK", "UTF-8"))
			}
			tmpOthFeeHs, _ := strconv.Atoi(v[213 : 213+10])
			tmpOthFeeLs, _ := strconv.Atoi(lsConMap[k][246 : 246+16])
			if tmpOthFeeHs != tmpOthFeeLs {
				logging.CmpInfo(disCode, "94文件信息:", k, "其他费用1/2/3不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[213:213+10]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][246:246+16]), "GBK", "UTF-8"))
			}

			if strings.TrimSpace(v[223:213+10+16+16]) != strings.TrimSpace(lsConMap[k][246+16:246+16+16+16]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "其他费用1/2/3不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[223:213+10+16+16]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][246+16:246+16+16+16]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[255:255+9]) != strings.TrimSpace(lsConMap[k][294:294+9]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "对方销售商不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[255:255+9]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][294:294+9]), "GBK", "UTF-8"))
			}

			if strings.TrimSpace(v[264:264+1]) != strings.TrimSpace(lsConMap[k][303:303+1]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "客户类型不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[264:264+1]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][303:303+1]), "GBK", "UTF-8"))
			}

			if strings.TrimSpace(v[273:273+4]) != strings.TrimSpace(lsConMap[k][312:312+4]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "返回代码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[273:273+4]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][312:312+4]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[277:277+17]) != strings.TrimSpace(lsConMap[k][316:316+17]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "交易账号不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[277:277+17]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][316:316+17]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[294:294+9]) != strings.TrimSpace(lsConMap[k][333:333+9]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "销售商代码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[294:294+9]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][333:333+9]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[303:303+16+16]) != strings.TrimSpace(lsConMap[k][342:342+16+16]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "申请份额/金额不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[303:303+16+16]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][342:342+16+16]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[335:335+3]) != strings.TrimSpace(lsConMap[k][374:374+3]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "业务代码不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[335:335+3]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][374:374+3]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[386:386+9+17+12+4]) != strings.TrimSpace(lsConMap[k][425:425+9+17+12+4]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "对方信息(网点,基金账号,交易账号,地区编码)不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[386:386+9+17+12+4]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][425:425+9+17+12+4]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[428:428+16]) != strings.TrimSpace(lsConMap[k][467:467+16]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "对方确认份额不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[428:428+16]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][467:467+16]), "GBK", "UTF-8"))
			}
			tmpNavHs, _ := strconv.Atoi(v[444 : 444+7])
			tmpNavLs, _ := strconv.Atoi(lsConMap[k][483 : 483+16])
			if tmpNavHs != tmpNavLs {
				logging.CmpInfo(disCode, "94文件信息:", k, "对方单位净值不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[444:444+7]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][483:483+17]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[451:451+7]) != strings.TrimSpace(lsConMap[k][499:499+7]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "实际折扣率不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[451:451+7]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][499:499+7]), "GBK", "UTF-8"))
			}
			/*
				tmpTrfFeeHs, _ := strconv.Atoi(v[488 : 488+10])
				tmpTrfFeeLs, _ := strconv.Atoi(lsConMap[k][566 : 566+16])
				if tmpTrfFeeHs != tmpTrfFeeLs {
					logging.CmpInfo(disCode, "94文件信息:", k, "过户费不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[488:488+10]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][566:566+16]), "GBK", "UTF-8"))
				}
			*/
			if strings.TrimSpace(v[520:520+16]) != strings.TrimSpace(lsConMap[k][604:604+16]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "利息转份额不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[428:428+16]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][604:604+16]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[536:536+16]) != strings.TrimSpace(lsConMap[k][620:620+16]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "未付收益不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[536:536+16]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][620:620+16]), "GBK", "UTF-8"))
			}
			/*
				if strings.TrimSpace(v[554:554+16]) != strings.TrimSpace(lsConMap[k][604:604+16]) {
					logging.CmpInfo(disCode, "94文件信息:", k, "后收费不一致")
					logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[554:554+16]), "GBK", "UTF-8"))
					logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][604:604+16]), "GBK", "UTF-8"))
				}
			*/
			if strings.TrimSpace(v[570:570+16]) != strings.TrimSpace(lsConMap[k][686:686+16]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "补差费不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[570:570+16]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][686:686+16]), "GBK", "UTF-8"))
			}
			if strings.TrimSpace(v[586:586+16]) != strings.TrimSpace(lsConMap[k][702:702+16]) {
				logging.CmpInfo(disCode, "94文件信息:", k, "未付收益不一致")
				logging.CmpInfo(disCode, "恒生数据", ConvertToString(strings.TrimSpace(v[586:586+16]), "GBK", "UTF-8"))
				logging.CmpInfo(disCode, "融先数据", ConvertToString(strings.TrimSpace(lsConMap[k][702:702+16]), "GBK", "UTF-8"))
			}

		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "94文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "94文件数据缺失", k)
			logging.CmpInfo(disCode, "融先94文件数据", ConvertToString(v, "GBK", "UTF-8"))
			logging.CmpInfo(disCode, "恒生系统导出94文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "94文件核对完毕")
	}

	return result
}

//T1件比对
func cmpT1File(hsfile []string, ls21file []string, disCode string) bool {
	//比对文件头
	var hsNum, lsNum int
	var result bool

	if len(hsfile) < 38 || len(ls21file) < 56 {
		logging.CmpInfo(disCode, "集中备份T1文件存在问题")
		logging.CmpInfo(disCode, "恒生文件行数:", len(hsfile))
		logging.CmpInfo(disCode, "融先文件行数:", len(ls21file))
		return false
	}

	//第1行为文件内容数量
	hsNum, _ = strconv.Atoi(hsfile[37])
	lsNum, _ = strconv.Atoi(ls21file[55])

	if hsNum != lsNum {
		logging.CmpInfo(disCode, "集中备份T1文件恒生导出数量与融先导出数量不一致")
		logging.CmpInfo(disCode, "恒生导出数量:", hsNum)
		logging.CmpInfo(disCode, "融先导出数量:", lsNum)
		//return false
	}
	hsContent := hsfile[38 : 38+hsNum]
	lsContent := ls21file[56 : 56+lsNum]
	result = cmpT1FileContent(hsContent, lsContent, disCode)

	return result
}

func cmpT1FileContent(hsContent []string, lsContent []string, disCode string) bool {
	//选取两套导出文件中可以确定唯一性字段为map的key，行内容为值，然后进行比对
	var hsConMap map[string]string = make(map[string]string)
	var lsConMap map[string]string = make(map[string]string)
	var result bool = true

	//加工恒生导出文件内容，使用基金账号，销售商代码，网点代码，交易账号作为key
	for _, v := range hsContent {
		hsConMap[v[:6]] = v
	}
	//加工融先导出文件内容，使用基金账号，销售商代码，网点代码，交易账号作为key
	for _, v := range lsContent {
		lsConMap[v[:6]] = v
	}

	//根据申请单号进行匹配
	for k, v := range hsConMap {
		//logging.Info("开始比对第", k, "行数据")
		if lsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "集中备份T1文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "集中备份T1文件数据缺失", k)
			logging.CmpInfo(disCode, "恒生T1文件数据", ConvertToString(v, "GBK", "UTF-8"))
			logging.CmpInfo(disCode, "融先系统导出T1文件无法找到匹配数据")
			continue
		}
		if v == lsConMap[k] {
			//logging.CmpInfo(disCode, "申请单号:", k, "核对通过")
		} else {
			if result {
				result = false
				logging.CmpInfo(disCode, "T1文件存在问题,请核对")
			}
		}
	}

	for k, v := range lsConMap {
		if hsConMap[k] == "" {
			if result {
				result = false
				logging.CmpInfo(disCode, "T1文件存在问题,请核对")
			}
			logging.CmpInfo(disCode, "T1文件数据缺失", k)
			logging.CmpInfo(disCode, "融先T1文件数据", ConvertToString(v, "GBK", "UTF-8"))
			logging.CmpInfo(disCode, "恒生系统导出T1文件无法找到匹配赎回")
			continue
		}
	}

	if !result {
		logging.CmpInfo(disCode, "T1文件核对完毕")
	}

	return result
}
