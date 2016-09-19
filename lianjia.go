package spider_lib

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	"github.com/henrylee2cn/pholcus/logs" //信息输出
	// net包
	// "net/http" //设置http.Header
	// "net/url"
	// 编码包
	// "encoding/xml"
	// "encoding/json"
	// 字符串处理包
	// "regexp"
	"strconv"
	// "strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	LianjiaXiaoqu.Register()
}

var LianjiaXiaoqu = &Spider{
	Name:        "lianjia",
	Description: "链家小区",
	// Pausetime: 300,
	// Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{Url: "http://tj.lianjia.com/xiaoqu/", Rule: "获取页码URL"})
		},

		Trunk: map[string]*Rule{

			"获取页码URL": {
				ParseFunc: func(ctx *Context) {
					// district := [7]string{"huangdao", "chengyang", "jiaozhou", "jiaonan", "jimo", "pingdu", "laixi"}
					// haspages := [7]int{17, 10, 2, 3, 2, 1, 1}
					// for j := 0; j <= 6; j++ {
					// 	logs.Log.Alert("第" + strconv.Itoa(j+1) + "个小区-共页数" + strconv.Itoa(haspages[j]))
					// 	for i := 1; i <= haspages[j]; i++ {
					// 		logs.Log.Critical("地区" + district[j] + "页码" + strconv.Itoa(i))
					// 		ctx.AddQueue(&request.Request{
					// 			Url:  "http://qd.lianjia.com/xiaoqu/" + district[j] + "/pg" + strconv.Itoa(i),
					// 			Rule: "小区列表",
					// 		})
					// 	}
					// }

					for i := 1; i <= 66; i++ {
						logs.Log.Critical("页码" + strconv.Itoa(i))
						ctx.AddQueue(&request.Request{
							Url:  "http://xm.lianjia.com/xiaoqu/pg" + strconv.Itoa(i),
							Rule: "小区列表",
						})
					}
				},
			},

			"小区列表": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					query.Find(".leftContent > .listContent > li").Each(func(i int, s *goquery.Selection) {
						logs.Log.Critical(s.Text())
						if url, ok := s.Find("a").Attr("href"); ok {
							// if i > 2 {
							// 	return
							// }
							logs.Log.Critical("小区" + url)
							ctx.AddQueue(&request.Request{Url: url, Rule: "xiaoqu_details"})
						}
					})
				},
			},

			"xiaoqu": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					var arrow string
					query.Find(".content_area > .location > .mod_cont > a").Each(func(n int, t *goquery.Selection) {
						arrow, _ = t.Attr("href")
					})

					logs.Log.Critical("坐标" + arrow)
					temp := ctx.CopyTemps()
					temp[ctx.GetItemField(15, "结果")] = arrow
					ctx.Output(temp)
				},
			},

			"组合temp": {
				ParseFunc: func(ctx *Context) {
					temp := ctx.CopyTemps()
					temp[ctx.GetItemField(15, "结果")] = "detail"
					ctx.Parse("结果")
				},
			},

			"结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"buildYears",
					"buildType",
					"propertyCosts",
					"propertyCompony",
					"business",
					"overview",
					"buildingTotal",
					"houseTotal",
					"nearbyStores",
					"detailTitle",
					"detailDesc",
					"station",
					"city",
					"district",
					"region",
					"arrow",
				},
				ParseFunc: func(ctx *Context) {
					// 结果存入Response中转
					ctx.Output(ctx.CopyTemps())
				},
			},

			"xiaoqu_details": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					var buildYears, buildType, propertyCosts, propertyCompony, business, overview, buildingTotal, houseTotal,
						nearbyStores, station, city, district, region string
					var xiaoqu_id string

					detailTitle := query.Find(".detailHeader > .detailTitle").First().Text()
					detailDesc := query.Find(".detailHeader > .detailDesc").First().Text()
					query.Find(".intro > .fl > a").Each(func(i int, fla *goquery.Selection) {
						switch i {
						case 0:
							station = fla.Text()
						case 1:
							city = fla.Text()
						case 2:
							district = fla.Text()
						case 3:
							region = fla.Text()
						case 4:
							xiaoqu_id, _ = fla.Attr("href")
						}
					})
					query.Find(".xiaoquInfo > .xiaoquInfoItem").Each(func(i int, s *goquery.Selection) {
						xiaoquInfoLabel := s.Find(".xiaoquInfoLabel").First().Text()
						xiaoquInfoContent := s.Find(".xiaoquInfoContent").First().Text()

						switch xiaoquInfoLabel {
						case "建筑年代":
							buildYears = xiaoquInfoContent
						case "建筑类型":
							buildType = xiaoquInfoContent
						case "物业费用":
							propertyCosts = xiaoquInfoContent
						case "物业公司":
							propertyCompony = xiaoquInfoContent
						case "开发商":
							business = xiaoquInfoContent
						case "小区概况":
							overview = xiaoquInfoContent
						case "楼栋总数":
							buildingTotal = xiaoquInfoContent
						case "房屋总数":
							houseTotal = xiaoquInfoContent
						case "附近门店":
							nearbyStores = xiaoquInfoContent
						}

					})
					temp := ctx.CreatItem(map[int]interface{}{
						0:  buildYears,
						1:  buildType,
						2:  propertyCosts,
						3:  propertyCompony,
						4:  business,
						5:  overview,
						6:  buildingTotal,
						7:  houseTotal,
						8:  nearbyStores,
						9:  detailTitle,
						10: detailDesc,
						11: station,
						12: city,
						13: district,
						14: region,
					}, "结果")

					// ctx.Parse("结果")
					// query.Find(".resblockQAAgent > .resblckQAEntrance > form > input").Each(func(i int, resblock *goquery.Selection) {
					// 	var inputName, _ = resblock.Attr("name")
					// 	if inputName == "xiaoqu_id" {
					// 		xiaoqu_id, _ = resblock.Attr("value")
					// 		logs.Log.Alert(xiaoqu_id)
					// 	}
					// })

					ctx.AddQueue(&request.Request{
						Url:  "http://m.lianjia.com/xm" + xiaoqu_id,
						Rule: "xiaoqu",
						Temp: temp,
					})
				},
			},
		},
	},
}
