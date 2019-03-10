package stack

import (
	"fmt"
	"github.com/Visteras/webhook-service-updater/docker"
	"github.com/spf13/viper"
	"gopkg.in/macaron.v1"
	"log"

	"net/http"
)

func Deploy(ctx *macaron.Context) {
	if ctx.Req.Method == http.MethodGet {
		service := ctx.Params(":service")
		filename := ctx.Params(":filename")
		args := []string{"stack", "deploy", service, "-c", "/app/files/yml/" + filename + ".yml"}

		log.Println(fmt.Sprintf("\033[1;32m Stack Deploy %s/%s from (%s/%s) \033[0m", service, filename, ctx.RemoteAddr(), ctx.Req.Header.Get(viper.GetString("WSU_USER"))))

		str, err := docker.DockerCmd(args...)
		if err != nil {
			ctx.Resp.Write([]byte(err.Error()))
		} else {
			ctx.Resp.Write([]byte(str))
		}
	}
}
