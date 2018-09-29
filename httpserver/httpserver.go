package httpserver

import (
	"bytes"
	"encoding/binary"
	"mapserver/assets/public"
	"mapserver/assets/templates"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/macaron.v1"
	log "github.com/Sirupsen/logrus"
	"github.com/go-macaron/bindata"
)

//StartServer web server
func StartServer(addr string, debug bool) {
	m := macaron.Classic()

	if debug {
		m.Use(macaron.Renderer())
		m.Use(macaron.Static("public"))
	} else {
		macaron.Env = macaron.PROD
		m.Use(macaron.Renderer(macaron.RenderOptions{
			TemplateFileSystem: bindata.Templates(bindata.Options{
				Asset:      templates.Asset,
				AssetDir:   templates.AssetDir,
				AssetNames: templates.AssetNames,
				AssetInfo:  templates.AssetInfo,
				Prefix:     "templates",
			}),
		}))
		m.Use(macaron.Static("public",
			macaron.StaticOptions{
				FileSystem: bindata.Static(bindata.Options{
					Asset:      public.Asset,
					AssetDir:   public.AssetDir,
					AssetNames: public.AssetNames,
					AssetInfo:  public.AssetInfo,
					Prefix:     "",
				}),
			},
		))
	}

	m.Use(func(ctx *macaron.Context) {
		s := strings.Split(addr, ":")
		ctx.Data["map_server_port"] = s[len(s)-1]
	})

	m.Get("/localmap/:t/:z/:x/:y", localmap)
	m.Get("/", index)

	fmt.Println("listen and run on ", addr)
	log.Fatal(http.ListenAndServe(addr, m))
}

func index(ctx *macaron.Context) {
	ctx.HTML(200, "mapdemo")
}

// 静态地图请求接口
func localmap(ctx *macaron.Context) {
	level := ctx.ParamsInt(":z")
	if level < 1 {
		level = 1
	}
	if level > 18 {
		level = 18
	}
	z := fmt.Sprintf("L%02d", level)
	x := ctx.ParamsInt(":x")
	y := ctx.ParamsInt(":y")
	t := ctx.Params(":t")
	if t == "" {
		t = "normal"
	}
	img := "C" + dec2Hex(int2Bytes(x)) + ".png"
	p := "R" + dec2Hex(int2Bytes(y))
	log.Infoln("localmap/" + t + "/" + z + "/" + p + "/" + img)
	ctx.ServeFileContent("localmap/" + t + "/" + z + "/" + p + "/" + img)
}

func dec2Hex(DecimalSlice []byte) string {
	var sa = make([]string, 0)
	for _, v := range DecimalSlice {
		parseStr := fmt.Sprintf("%02x", v)
		sa = append(sa, parseStr)
	}
	ss := strings.Join(sa, "")
	return ss
}

func int2Bytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}
