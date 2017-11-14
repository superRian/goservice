package pholcus_lib

import (
	"log"
	"strconv"
	"strings"

	"github.com/henrylee2cn/pholcus/app/downloader/request"
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"
)

func init() {
	Jiban.Register()
}

var Jiban = &Spider{
	Name:         "羁绊二次元资讯",
	Description:  "羁绊测试 [http://www.005.tv/zx/]",
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{
				Url:         "http://www.005.tv/zx/list_526_1.html",
				Rule:        "请求",
				Temp:        map[string]interface{}{"p": 1},
				ConnTimeout: -1,
				Reloadable:  true,
			})

		},
		Trunk: map[string]*Rule{
			"请求": {
				ParseFunc: func(ctx *Context) {
					var curr = ctx.GetTemp("p", int(0)).(int)
					log.Println("curr--->>", curr)
					log.Println("curr--->>>")
					ctx.GetDom().Find(".pages .dede_pages  .pagelist  .thisclass a").Each(func(ii int, iio *goquery.Selection) {
						url2, _ := iio.Attr("href")
						log.Println("iio--", iio.Text())
						log.Println("url--", url2)
						if url2 != "javascript:void(0);" {
							if curr > 100 {
								return
							}
						}
					})
					//					if ctx.GetDom().Find(".pages .dede_pages  .pagelist  .thisclass  select option:selected").Text() != strconv.Itoa(curr) || curr > 100 {
					//						return
					//					}
					log.Println("-----39-------")
					ctx.AddQueue(&request.Request{
						Url:         "http://www.005.tv/zx/list_526_" + strconv.Itoa(curr+1) + ".html",
						Rule:        "请求",
						Temp:        map[string]interface{}{"p": curr + 1},
						ConnTimeout: -1,
						Reloadable:  true,
					})
					log.Println("-----47-------")
					ctx.Parse("获取列表")
				},
			},

			"获取列表": {
				ParseFunc: func(ctx *Context) {
					log.Println("65 Line ----------", ctx.GetDom().
						Find(".article-list ul li .xs-100 div h3 a").Text())
					ctx.GetDom().
						Find(".article-list ul li .xs-100 div h3 a").
						Each(func(i int, s *goquery.Selection) {
							url, _ := s.Attr("href")
							log.Println("69 Line Url:", url)
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
					"img_url",
					"content",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					var title, time, img_url, content string
					log.Println("--------91")
					query.Find(".article-list-wrap").
						Each(func(j int, jo *goquery.Selection) {
							log.Println("94 line JO")
							title = jo.Find(".articleTitle-name").Text()
							log.Println("Title ---", title)
							time = jo.Find("span.time").Text()
							log.Println("Time---", time)
							jo.Find(".articleContent img").Each(func(x int, xo *goquery.Selection) {
								if img, ok := xo.Attr("src"); ok {
									img_url = img_url + img + ","
									log.Println("img------", img_url)
								}
							})
							jo.Find(".articleContent img").ReplaceWithHtml("#image#")
							jo.Find(".articleContent img").Remove()

							content, _ = jo.Find(".articleContent").Html()
							content = strings.Replace(content, `"`, `'`, -1)
							log.Println("content------", content)
						})
					ctx.Output(map[int]interface{}{
						0: title,
						1: time,
						2: img_url,
						3: content,
					})
				},
			},
		},
	},
}
