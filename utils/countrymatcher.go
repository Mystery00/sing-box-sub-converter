package utils

import (
	"regexp"
	"strings"
)

// CountryPatterns 定义国家代码和对应的正则表达式模式
var CountryPatterns = map[string]*regexp.Regexp{
	"🇭🇰": regexp.MustCompile(`香港|沪港|呼港|中港|HKT|HKBN|HGC|WTT|CMI|穗港|广港|京港|🇭🇰|HK|Hongkong|Hong Kong|HongKong|HONG KONG`),
	"🇹🇼": regexp.MustCompile(`台湾|台灣|臺灣|台北|台中|新北|彰化|台|CHT|HINET|TW|Taiwan|TAIWAN`),
	"🇲🇴": regexp.MustCompile(`澳门|澳門|(\s|-)?MO\d*|CTM|MAC|Macao|Macau`),
	"🇸🇬": regexp.MustCompile(`新加坡|狮城|獅城|沪新|京新|泉新|穗新|深新|杭新|广新|廣新|滬新|SG|Singapore|SINGAPORE`),
	"🇯🇵": regexp.MustCompile(`日本|东京|東京|大阪|埼玉|京日|苏日|沪日|广日|上日|穗日|川日|中日|泉日|杭日|深日|JP|Japan|JAPAN`),
	"🇺🇸": regexp.MustCompile(`美国|美國|京美|硅谷|凤凰城|洛杉矶|西雅图|圣何塞|芝加哥|哥伦布|纽约|广美|(\s|-)?US\d*|USA|America|United States`),
	"🇰🇷": regexp.MustCompile(`韩国|韓國|首尔|首爾|韩|韓|春川|KOR|KR|Kr|(North\s)?Korea`),
	"🇰🇵": regexp.MustCompile(`朝鲜|KP|North Korea`),
	"🇷🇺": regexp.MustCompile(`俄罗斯|俄羅斯|毛子|俄国|RU|RUS|Russia`),
	"🇮🇳": regexp.MustCompile(`印度|孟买|(\s|-)?IN\d*|IND|India|INDIA|Mumbai`),
	"🇮🇩": regexp.MustCompile(`印尼|印度尼西亚|雅加达|ID|IDN|Indonesia`),
	"🇬🇧": regexp.MustCompile(`英国|英國|伦敦|UK|England|United Kingdom|Britain`),
	"🇩🇪": regexp.MustCompile(`德国|德國|法兰克福|(\s|-)?DE\d*|(\s|-)?GER\d*|🇩🇪|German|GERMAN`),
	"🇫🇷": regexp.MustCompile(`法国|法國|巴黎|FR|France`),
	"🇩🇰": regexp.MustCompile(`丹麦|丹麥|DK|DNK|Denmark`),
	"🇳🇴": regexp.MustCompile(`挪威|(\s|-)?NO\d*|Norway`),
	"🇮🇹": regexp.MustCompile(`意大利|義大利|米兰|(\s|-)?IT\d*|Italy|Nachash`),
	"🇻🇦": regexp.MustCompile(`梵蒂冈|梵蒂岡|(\s|-)?VA\d*|Vatican City`),
	"🇧🇪": regexp.MustCompile(`比利时|比利時|(\s|-)?BE\d*|Belgium`),
	"🇦🇺": regexp.MustCompile(`澳大利亚|澳洲|墨尔本|悉尼|(\s|-)?AU\d*|Australia|Sydney`),
	"🇨🇦": regexp.MustCompile(`加拿大|蒙特利尔|温哥华|多伦多|多倫多|滑铁卢|楓葉|枫叶|CA|CAN|Waterloo|Canada|CANADA`),
	"🇲🇾": regexp.MustCompile(`马来西亚|马来|馬來|MY|Malaysia|MALAYSIA`),
	"🇲🇻": regexp.MustCompile(`马尔代夫|馬爾代夫|(\s|-)?MV\d*|Maldives`),
	"🇹🇷": regexp.MustCompile(`土耳其|伊斯坦布尔|(\s|-)?TR\d|TR_|TUR|Turkey`),
	"🇵🇭": regexp.MustCompile(`菲律宾|菲律賓|(\s|-)?PH\d*|Philippines`),
	"🇹🇭": regexp.MustCompile(`泰国|泰國|曼谷|(\s|-)?TH\d*|Thailand`),
	"🇻🇳": regexp.MustCompile(`越南|胡志明市|(\s|-)?VN\d*|Vietnam`),
	"🇰🇭": regexp.MustCompile(`柬埔寨|(\s|-)?KH\d*|Cambodia`),
	"🇱🇦": regexp.MustCompile(`老挝|(\s|-)LA\d*|Laos`),
	"🇧🇩": regexp.MustCompile(`孟加拉|(\s|-)?BD\d*|Bengal`),
	"🇲🇲": regexp.MustCompile(`缅甸|緬甸|(\s|-)?MM\d*|Myanmar`),
	"🇱🇧": regexp.MustCompile(`黎巴嫩|(\s|-)?LB\d*|Lebanon`),
	"🇺🇦": regexp.MustCompile(`乌克兰|烏克蘭|(\s|-)?UA\d*|Ukraine`),
	"🇭🇺": regexp.MustCompile(`匈牙利|(\s|-)?HU\d*|Hungary`),
	"🇨🇭": regexp.MustCompile(`瑞士|苏黎世|(\s|-)?CH\d*|Switzerland`),
	"🇸🇪": regexp.MustCompile(`瑞典|SE|Sweden`),
	"🇱🇺": regexp.MustCompile(`卢森堡|(\s|-)?LU\d*|Luxembourg`),
	"🇦🇹": regexp.MustCompile(`奥地利|奧地利|维也纳|(\s|-)?AT\d*|Austria`),
	"🇨🇿": regexp.MustCompile(`捷克|(\s|-)?CZ\d*|Czechia`),
	"🇬🇷": regexp.MustCompile(`希腊|希臘|(\s|-)?GR\d*|Greece`),
	"🇮🇸": regexp.MustCompile(`冰岛|冰島|(\s|-)?IS\d*|ISL|Iceland`),
	"🇳🇿": regexp.MustCompile(`新西兰|新西蘭|(\s|-)?NZ\d*|New Zealand`),
	"🇮🇪": regexp.MustCompile(`爱尔兰|愛爾蘭|都柏林|(\s|-)?IE\d*|Ireland|IRELAND`),
	"🇮🇲": regexp.MustCompile(`马恩岛|馬恩島|(\s|-)?IM\d*|Mannin|Isle of Man`),
	"🇱🇹": regexp.MustCompile(`立陶宛|(\s|-)?LT\d*|Lithuania`),
	"🇫🇮": regexp.MustCompile(`芬兰|芬蘭|赫尔辛基|(\s|-)?FI\d*|Finland`),
	"🇦🇷": regexp.MustCompile(`阿根廷|(\s|-)AR\d*|Argentina`),
	"🇺🇾": regexp.MustCompile(`乌拉圭|烏拉圭|(\s|-)?UY\d*|Uruguay`),
	"🇵🇾": regexp.MustCompile(`巴拉圭|(\s|-)?PY\d*|Paraguay`),
	"🇯🇲": regexp.MustCompile(`牙买加|牙買加|(\s|-)?JM\d*|Jamaica`),
	"🇸🇷": regexp.MustCompile(`苏里南|蘇里南|(\s|-)?SR\d*|Suriname`),
	"🇨🇼": regexp.MustCompile(`库拉索|庫拉索|(\s|-)?CW\d*|Curaçao`),
	"🇨🇴": regexp.MustCompile(`哥伦比亚|(\s|-)?CO\d*|Colombia`),
	"🇪🇨": regexp.MustCompile(`厄瓜多尔|(\s|-)?EC\d*|Ecuador`),
	"🇪🇸": regexp.MustCompile(`西班牙|\b(\s|-)?ES\d*|Spain`),
	"🇵🇹": regexp.MustCompile(`葡萄牙|Portugal`),
	"🇮🇱": regexp.MustCompile(`以色列|(\s|-)?IL\d*|Israel`),
	"🇸🇦": regexp.MustCompile(`沙特|利雅得|吉达|Saudi|Saudi Arabia`),
	"🇲🇳": regexp.MustCompile(`蒙古|(\s|-)?MN\d*|Mongolia`),
	"🇦🇪": regexp.MustCompile(`阿联酋|迪拜|(\s|-)?AE\d*|Dubai|United Arab Emirates`),
	"🇦🇿": regexp.MustCompile(`阿塞拜疆|(\s|-)?AZ\d*|Azerbaijan`),
	"🇦🇲": regexp.MustCompile(`亚美尼亚|亞美尼亞|(\s|-)?AM\d*|Armenia`),
	"🇰🇿": regexp.MustCompile(`哈萨克斯坦|哈薩克斯坦|(\s|-)?KZ\d*|Kazakhstan`),
	"🇰🇬": regexp.MustCompile(`吉尔吉斯坦|吉尔吉斯斯坦|(\s|-)?KG\d*|Kyrghyzstan`),
	"🇺🇿": regexp.MustCompile(`乌兹别克斯坦|烏茲別克斯坦|(\s|-)?UZ\d*|Uzbekistan`),
	"🇧🇷": regexp.MustCompile(`巴西|圣保罗|维涅杜|BR|Brazil`),
	"🇨🇱": regexp.MustCompile(`智利|(\s|-)?CL\d*|Chile|CHILE`),
	"🇵🇪": regexp.MustCompile(`秘鲁|祕魯|(\s|-)?PE\d*|Peru`),
	"🇨🇺": regexp.MustCompile(`古巴|Cuba`),
	"🇧🇹": regexp.MustCompile(`不丹|Bhutan`),
	"🇦🇩": regexp.MustCompile(`安道尔|(\s|-)?AD\d*|Andorra`),
	"🇲🇹": regexp.MustCompile(`马耳他|(\s|-)?MT\d*|Malta`),
	"🇲🇨": regexp.MustCompile(`摩纳哥|摩納哥|(\s|-)?MC\d*|Monaco`),
	"🇷🇴": regexp.MustCompile(`罗马尼亚|(\s|-)?RO\d*|Rumania`),
	"🇧🇬": regexp.MustCompile(`保加利亚|保加利亞|(\s|-)?BG\d*|Bulgaria`),
	"🇭🇷": regexp.MustCompile(`克罗地亚|克羅地亞|(\s|-)?HR\d*|Croatia`),
	"🇲🇰": regexp.MustCompile(`北马其顿|北馬其頓|(\s|-)?MK\d*|North Macedonia`),
	"🇷🇸": regexp.MustCompile(`塞尔维亚|塞爾維亞|(\s|-)?RS\d*|Seville|Sevilla`),
	"🇨🇾": regexp.MustCompile(`塞浦路斯|(\s|-)?CY\d*|Cyprus`),
	"🇱🇻": regexp.MustCompile(`拉脱维亚|(\s|-)?LV\d*|Latvia|Latvija`),
	"🇲🇩": regexp.MustCompile(`摩尔多瓦|摩爾多瓦|(\s|-)?MD\d*|Moldova`),
	"🇸🇰": regexp.MustCompile(`斯洛伐克|(\s|-)?SK\d*|Slovakia`),
	"🇪🇪": regexp.MustCompile(`爱沙尼亚|(\s|-)?EE\d*|Estonia`),
	"🇧🇾": regexp.MustCompile(`白俄罗斯|白俄羅斯|(\s|-)?BY\d*|White Russia|Republic of Belarus|Belarus`),
	"🇧🇳": regexp.MustCompile(`文莱|汶萊|BRN|Negara Brunei Darussalam`),
	"🇬🇺": regexp.MustCompile(`关岛|關島|(\s|-)?GU\d*|Guam`),
	"🇫🇯": regexp.MustCompile(`斐济|斐濟|(\s|-)?FJ\d*|Fiji`),
	"🇯🇴": regexp.MustCompile(`约旦|約旦|(\s|-)?JO\d*|Jordan`),
	"🇬🇪": regexp.MustCompile(`格鲁吉亚|格魯吉亞|(\s|-)?GE\d*|Georgia`),
	"🇬🇮": regexp.MustCompile(`直布罗陀|直布羅陀|(\s|-)GI\d*|Gibraltar`),
	"🇸🇲": regexp.MustCompile(`圣马力诺|聖馬利諾|(\s|-)?SM\d*|San Marino`),
	"🇳🇵": regexp.MustCompile(`尼泊尔|(\s|-)?NP\d*|Nepal`),
	"🇫🇴": regexp.MustCompile(`法罗群岛|法羅群島|(\s|-)FO\d*|Faroe Islands`),
	"🇦🇽": regexp.MustCompile(`奥兰群岛|奧蘭群島|(\s|-)?AX\d*|Åland`),
	"🇸🇮": regexp.MustCompile(`斯洛文尼亚|斯洛文尼亞|(\s|-)?SI\d*|Slovenia`),
	"🇦🇱": regexp.MustCompile(`阿尔巴尼亚|阿爾巴尼亞|(\s|-)?AL\d*|Albania`),
	"🇹🇱": regexp.MustCompile(`东帝汶|東帝汶|(\s|-)?TL\d*|East Timor`),
	"🇵🇦": regexp.MustCompile(`巴拿马|巴拿馬|(\s|-)?PA\d*|Panama`),
	"🇧🇲": regexp.MustCompile(`百慕大|(\s|-)?BM\d*|Bermuda`),
	"🇬🇱": regexp.MustCompile(`格陵兰|格陵蘭|(\s|-)?GL\d*|Greenland`),
	"🇨🇷": regexp.MustCompile(`哥斯达黎加|(\s|-)?CR\d*|Costa Rica`),
	"🇻🇬": regexp.MustCompile(`英属维尔京|(\s|-)?VG\d*|British Virgin Islands`),
	"🇻🇮": regexp.MustCompile(`美属维尔京|(\s|-)?VI\d*|United States Virgin Islands`),
	"🇲🇽": regexp.MustCompile(`墨西哥|MX|MEX|MEX|MEXICO`),
	"🇲🇪": regexp.MustCompile(`黑山|(\s|-)?ME\d*|Montenegro`),
	"🇳🇱": regexp.MustCompile(`荷兰|荷蘭|尼德蘭|阿姆斯特丹|NL|Netherlands|Amsterdam`),
	"🇵🇱": regexp.MustCompile(`波兰|波蘭|(\s|-)?PL\d*|POL|Poland`),
	"🇩🇿": regexp.MustCompile(`阿尔及利亚|(\s|-)?DZ\d*|Algeria`),
	"🇧🇦": regexp.MustCompile(`波黑共和国|波黑|(\s|-)?BA\d*|Bosnia and Herzegovina`),
	"🇱🇮": regexp.MustCompile(`列支敦士登|(\s|-)?LI\d*|Liechtenstein`),
	"🇷🇪": regexp.MustCompile(`留尼汪|留尼旺|(\s|-)?RE\d*|Réunion|Reunion`),
	"🇿🇦": regexp.MustCompile(`南非|约翰内斯堡|(\s|-)?ZA\d*|South Africa|Johannesburg`),
	"🇪🇬": regexp.MustCompile(`埃及|(\s|-)?EG\d*|Egypt`),
	"🇬🇭": regexp.MustCompile(`加纳|(\s|-)?GH\d*|Ghana`),
	"🇲🇱": regexp.MustCompile(`马里|馬里|(\s|-)?ML\d*|Mali`),
	"🇲🇦": regexp.MustCompile(`摩洛哥|(\s|-)?MA\d*|Morocco`),
	"🇹🇳": regexp.MustCompile(`突尼斯|(\s|-)?TN\d*|Tunisia`),
	"🇱🇾": regexp.MustCompile(`利比亚|(\s|-)?LY\d*|Libya`),
	"🇰🇪": regexp.MustCompile(`肯尼亚|肯尼亞|(\s|-)?KE\d*|Kenya`),
	"🇷🇼": regexp.MustCompile(`卢旺达|盧旺達|(\s|-)?RW\d*|Rwanda`),
	"🇨🇻": regexp.MustCompile(`佛得角|維德角|(\s|-)?CV\d*|Cape Verde`),
	"🇦🇴": regexp.MustCompile(`安哥拉|(\s|-)?AO\d*|Angola`),
	"🇳🇬": regexp.MustCompile(`尼日利亚|尼日利亞|拉各斯|(\s|-)?NG\d*|Nigeria`),
	"🇲🇺": regexp.MustCompile(`毛里求斯|(\s|-)?MU\d*|Mauritius`),
	"🇴🇲": regexp.MustCompile(`阿曼|(\s|-)?OM\d*|Oman`),
	"🇧🇭": regexp.MustCompile(`巴林|(\s|-)?BH\d*|Bahrain`),
	"🇮🇶": regexp.MustCompile(`伊拉克|(\s|-)?IQ\d*|Iraq`),
	"🇮🇷": regexp.MustCompile(`伊朗|(\s|-)?IR\d*|Iran`),
	"🇦🇫": regexp.MustCompile(`阿富汗|(\s|-)?AF\d*|Afghanistan`),
	"🇵🇰": regexp.MustCompile(`巴基斯坦|(\s|-)?PK\d*|Pakistan|PAKISTAN`),
	"🇶🇦": regexp.MustCompile(`卡塔尔|卡塔爾|(\s|-)?QA\d*|Qatar`),
	"🇸🇾": regexp.MustCompile(`叙利亚|敘利亞|(\s|-)?SY\d*|Syria`),
	"🇱🇰": regexp.MustCompile(`斯里兰卡|斯里蘭卡|(\s|-)?LK\d*|Sri Lanka`),
	"🇻🇪": regexp.MustCompile(`委内瑞拉|(\s|-)?VE\d*|Venezuela`),
	"🇬🇹": regexp.MustCompile(`危地马拉|(\s|-)?GT\d*|Guatemala`),
	"🇵🇷": regexp.MustCompile(`波多黎各|(\s|-)?PR\d*|Puerto Rico`),
	"🇰🇾": regexp.MustCompile(`开曼群岛|開曼群島|盖曼群岛|凯门群岛|(\s|-)?KY\d*|Cayman Islands`),
	"🇸🇯": regexp.MustCompile(`斯瓦尔巴|扬马延|(\s|-)?SJ\d*|Svalbard|Mayen`),
	"🇭🇳": regexp.MustCompile(`洪都拉斯|Honduras`),
	"🇳🇮": regexp.MustCompile(`尼加拉瓜|(\s|-)?NI\d*|Nicaragua`),
	"🇦🇶": regexp.MustCompile(`南极|南極|(\s|-)?AQ\d*|Antarctica`),
	"🇨🇳": regexp.MustCompile(`中国|中國|江苏|北京|上海|广州|深圳|杭州|徐州|青岛|宁波|镇江|沈阳|济南|回国|back|(\s|-)?CN\d*|China`),
}

// renameNodeTagWithEmoji 根据字符串内容匹配国家/地区代码，并重命名字符串
func renameNodeTagWithEmoji(originTagName string) string {
	// 排除特殊关键词的情况
	if strings.Contains(originTagName, "INFO") ||
		strings.Contains(originTagName, "FREE") ||
		strings.Contains(originTagName, "GRPC") ||
		strings.Contains(originTagName, "IEPL") ||
		strings.Contains(originTagName, "ARP") ||
		strings.Contains(originTagName, "WARP") ||
		strings.Contains(originTagName, "JMS") ||
		strings.Contains(originTagName, "GER") ||
		strings.Contains(originTagName, "GIA") ||
		strings.Contains(originTagName, "CN2GI") ||
		strings.Contains(originTagName, "TLS") ||
		strings.Contains(originTagName, "RELAY") ||
		strings.Contains(originTagName, "BGP") {
		// 如果包含这些特殊关键词，则继续检查其他匹配模式，而不是直接返回
	}

	for countryCode, pattern := range CountryPatterns {
		// 跳过那些带有特殊关键词标记的匹配器
		if strings.Contains(countryCode, "INFO") ||
			strings.Contains(countryCode, "FREE") ||
			strings.Contains(countryCode, "GRPC") ||
			strings.Contains(countryCode, "IEPL") ||
			strings.Contains(countryCode, "ARP") ||
			strings.Contains(countryCode, "WARP") ||
			strings.Contains(countryCode, "JMS") ||
			strings.Contains(countryCode, "GER") ||
			strings.Contains(countryCode, "GIA") ||
			strings.Contains(countryCode, "CN2GI") ||
			strings.Contains(countryCode, "TLS") ||
			strings.Contains(countryCode, "RELAY") ||
			strings.Contains(countryCode, "BGP") {
			continue
		}

		// 检查字符串是否已经以国家代码开头
		if strings.HasPrefix(originTagName, countryCode) {
			return countryCode + " " + strings.TrimSpace(originTagName[len(countryCode):])
		}

		// 使用正则表达式搜索匹配
		if pattern.MatchString(originTagName) {
			// 特殊处理 🇺🇲 开头的情况
			if strings.HasPrefix(originTagName, "🇺🇲") {
				return countryCode + " " + strings.TrimSpace(originTagName[len("🇺🇲"):])
			} else {
				return countryCode + " " + originTagName
			}
		}
	}

	// 如果没有找到匹配，返回原始字符串
	return originTagName
}
