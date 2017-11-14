package pholcus_lib

import (
	"log"
	"strconv"
	"strings"

	"github.com/henrylee2cn/pholcus/app/downloader/request"
	. "github.com/henrylee2cn/pholcus/app/spider"
	"github.com/henrylee2cn/pholcus/common/goquery" //DOM解析
)

func init() {
	ACG.Register()
}

var ACG = &Spider{
	Name:         "测试",
	Description:  "狗屋咨询 [Auto Page] [www.acgdoge.net]",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:        "http://www.acgdoge.net/page/1",
				Rule:       "数据列表",
				Temp:       map[string]interface{}{"p": 1},
				Reloadable: true,
			})
		},
		Trunk: map[string]*Rule{
			"数据列表": {
				ParseFunc: func(ctx *Context) {
					var curr = ctx.GetTemp("p", int(0)).(int)
					if ctx.GetDom().Find("#page_nav .current").Text() != strconv.Itoa(curr) || curr > 100 {
						return
					}
					ctx.AddQueue(&request.Request{
						Url:         "http://www.acgdoge.net/page/" + strconv.Itoa(curr+1),
						Rule:        "数据列表",
						Temp:        map[string]interface{}{"p": curr + 1},
						ConnTimeout: -1,
						Reloadable:  true,
					})
					ctx.Parse("获取列表")
				},
			},

			"获取列表": {
				ParseFunc: func(ctx *Context) {
					ctx.GetDom().
						Find(".post_h_l h2 a").
						Each(func(i int, s *goquery.Selection) {
							url, _ := s.Attr("href")
							log.Println("50 line url:", url)
							ctx.AddQueue(&request.Request{
								Url:         url,
								Rule:        "news",
								ConnTimeout: -1,
							})
						})
				},
			},

			"news": {

				ItemFields: []string{
					"title",
					"time",
					"content",
					"source_url",
					"state",
					"SID",
					"img_url",
					"postid",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					var title, time, content, source_url, SID, img_url string
					var state, postid int
					query.Find("#content").Each(func(i int, s *goquery.Selection) {
						title = s.Find(".post_h_l h2").Text()
						time = s.Find(".post_time .post_t_d").Text() + " " + s.Find(".post_time .post_t_u").Text()
						log.Println("time-----------time", time)
						source_url = ctx.GetUrl()
						state = 0
						postid = 0
						SID = "acg"
						s.Find(".post_t").RemoveAttr("href")
						s.Find(".post_t").RemoveAttr("title")
						s.Find(".post_t noscript").Each(func(j int, jo *goquery.Selection) {
							joStr, _ := jo.Html()
							log.Println("---jostr:---", joStr)
							joStr = strings.Replace(joStr, `&lt;`, `<`, -1)
							joStr = strings.Replace(joStr, `&gt;`, `>`, -1)
							jo.ReplaceWithHtml("<span> " + joStr + "</span>")
						})
						s.Find(".post_t span img").Each(func(x int, xo *goquery.Selection) {
							if img, ok := xo.Attr("src"); ok {
								img_url = img_url + img + ","
							}
						})
						s.Find(".post_t span").Remove()
						s.Find(".post_t img").ReplaceWithHtml("#image")
						img_url = strings.Replace(img_url, `"`, ``, -1)
						content, _ = s.Find(".post_t").Html()
						content = strings.Replace(content, `"`, `'`, -1) + "<p> acg 来源</p>"

					})
					ctx.Output(map[int]interface{}{
						0: title,
						1: time,
						2: content,
						3: source_url,
						4: state,
						5: SID,
						6: img_url,
						7: postid,
					})
				},
			},
		},
	},
}
