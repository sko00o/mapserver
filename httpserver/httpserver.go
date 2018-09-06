package httpserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"mapserver/config"
	"strconv"
	"strings"

	"gopkg.in/macaron.v1"
)

func StartHttpServer(addr string, conf config.Config) {

	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Use(func(ctx *macaron.Context) {
		ctx.Data["server_host"] = addr
		ctx.Data["latitude"] = conf.CenterLat
		ctx.Data["longitude"] = conf.CenterLng
		ctx.Data["zoom_level"] = conf.ZoomLevel
	})
	m.Get("/", index)
	m.Get("/localmap/:z/:x/:y", loalmap) //静态地图请求接口
	m.RunAddr(addr)
}

func index(ctx *macaron.Context) {
	ctx.HTML(200, "mapdemo")
}

// BasePath 地图瓦片的根目录
const (
	BasePath = "public/_alllayers/"
)

//解析前端请求参数，然后解析返回所需要的瓦片
func loalmap(ctx *macaron.Context) {
	zoom := "L" + ctx.Params(":z")
	x, _ := strconv.Atoi(ctx.Params(":x"))
	y, _ := strconv.Atoi(ctx.Params(":y"))
	imgName := "C" + dec2Hex(int2Bytes(x)) + ".png"
	path := "R" + dec2Hex(int2Bytes(y))
	fmt.Println(BasePath + zoom + "/" + path + "/" + imgName)

	ctx.ServeFileContent(BasePath + zoom + "/" + path + "/" + imgName)
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
