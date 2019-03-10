package service

import (
	"fmt"
	"github.com/Visteras/webhook-service-updater/docker"
	"github.com/spf13/viper"
	"gopkg.in/macaron.v1"
	"log"
	"net/http"
)

func Update(ctx *macaron.Context) {
	if ctx.Req.Method == http.MethodGet {
		service := ctx.Params(":service")
		args := []string{"service", "update", service}
		if ctx.Req.URL.Query().Get("force") == "true" {
			args = append(args, "--force")
		}
		if ctx.Req.URL.Query().Get("with-registry-auth") == "true" {
			args = append(args, "--with-registry-auth")
		}
		log.Println(fmt.Sprintf("\033[1;32m Update service %s from (%s/%s) \033[0m", service, ctx.RemoteAddr(), ctx.Req.Header.Get(viper.GetString("WSU_USER"))))
		str, err := docker.DockerCmd(args...)
		if err != nil {
			ctx.Resp.Write([]byte(err.Error()))
		} else {
			ctx.Resp.Write([]byte(str))
		}
	}
}

func Log(ctx *macaron.Context) {
	if ctx.Req.Method == http.MethodGet {
		service := ctx.Params(":service")
		args := []string{"service", "logs", service}
		log.Println(fmt.Sprintf("\033[1;32m Logs service %s from (%s/%s) \033[0m", service, ctx.RemoteAddr(), ctx.Req.Header.Get(viper.GetString("WSU_USER"))))
		str, err := docker.DockerCmd(args...)
		if err != nil {
			ctx.Resp.Write([]byte(err.Error()))
		} else {
			ctx.Resp.Write([]byte(str))
		}
	}
}
