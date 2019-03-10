package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Visteras/webhook-service-updater/docker"
	"github.com/Visteras/webhook-service-updater/service"
	"github.com/Visteras/webhook-service-updater/stack"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/macaron.v1"
	"log"
	"net/http"
	"strings"
)

var Prefix string
var WsuUser string
var WsuToken string

type User struct {
	Tokens   []string `json:"tokens"`
	Services []string `json:"services"`
	Stacks   []struct {
		Service  string `json:"service"`
		Filename string `json:"filename"`
	} `json:"stacks"`
	IsAdmin bool     `json:"admin"`
	LockIP  bool     `json:"lock_ip"`
	Ips     []string `json:"ips"`
}

func main() {
	//----- Start binding
	WsuUser = GetEnvWithDefault("WSU_USER", "WSU_USER")
	WsuToken = GetEnvWithDefault("WSU_TOKEN", "WSU_TOKEN")
	Prefix = GetEnvWithDefault("WSU_PREFIX", "")
	docker.DockerExec = GetEnvWithDefault("DOCKER_EXEC", "docker")
	//----- End binding

	viper.SetConfigName("config")
	viper.AddConfigPath("./files/")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.SetConfigType("json")

	m := macaron.Classic()

	m.Group("", func() {
		m.Group(Prefix, func() {
			m.Group("/service", func() {
				m.Get("/update/:service", service.Update)
				m.Get("/logs/:service", service.Log)
			}, CheckServiceAccess)
			m.Group("/stack", func() {
				m.Get("/deploy/:service/:filename", stack.Deploy)
			}, CheckStackAccess)
		})
	}, CheckAccess, CheckIPAccess)

	m.Run()
}

func CheckAccess(ctx *macaron.Context) {
	myErr := false
	user := User{}
	auth := false

	username, err := GetUser(ctx.Req.Header.Get(WsuUser), &user)
	if err != nil {
		log.Printf(err.Error())
		myErr = true
	} else {
		token := ctx.Req.Header.Get(WsuToken)
		for _, v := range user.Tokens {
			if v == token {
				log.Println(fmt.Sprintf("\033[1;32m Auth %s from %s \033[0m", username, ctx.Req.RemoteAddr))
				auth = true
			}
		}
	}

	if myErr || !auth {
		http.Error(ctx.Resp, "Access Denied", http.StatusForbidden)
		return
	}
}

func CheckIPAccess(ctx *macaron.Context) {
	user := User{}
	ip := strings.Split(ctx.Req.RemoteAddr, ":")[0]
	_, _ = GetUser(ctx.Req.Header.Get(WsuUser), &user)
	myErr := true

	if user.LockIP {
		for key, value := range user.Ips {
			if value == ip {
				log.Println(fmt.Sprintf("\033[1;32m%s\033[0m", user.Ips[key]))
				myErr = false
			}
		}

		if myErr {
			http.Error(ctx.Resp, "Access Denied from this IP", http.StatusForbidden)
			return
		}
	}
}

func CheckServiceAccess(ctx *macaron.Context) {
	user := User{}
	s := ctx.Params(":service")
	_, _ = GetUser(ctx.Req.Header.Get(WsuUser), &user)
	myErr := true

	for key, value := range user.Services {
		if value == s {
			log.Println(fmt.Sprintf("\033[1;32m%s\033[0m", user.Services[key]))
			myErr = false

		}
	}

	if myErr {
		http.Error(ctx.Resp, "Access Denied", http.StatusForbidden)
		return
	}

}
func CheckStackAccess(ctx *macaron.Context) {
	user := User{}
	s := ctx.Params(":service")
	fn := ctx.Params(":filename")
	_, _ = GetUser(ctx.Req.Header.Get(WsuUser), &user)
	myErr := true

	for _, value := range user.Stacks {
		if value.Filename == fn && value.Service == s {
			log.Println(fmt.Sprintf("\033[1;32m%s/%s\033[0m", s, fn))
			myErr = false
		}
	}
	if myErr && user.IsAdmin {
		myErr = false
		log.Println(fmt.Sprintf("\033[1;32m Admin(%s) Access for %s/%s \033[0m", user, s, fn))
	}

	if myErr {
		http.Error(ctx.Resp, "Access Denied", http.StatusForbidden)
		return
	}
}

func GetUser(username string, user *User) (string, error) {
	users := viper.Sub("users")
	if users.IsSet(username) {
		u := users.Get(username)

		//FIXME Исправить это извращение
		u2, err := json.Marshal(u)
		if err != nil {
			return username, err
		}

		err = json.Unmarshal(u2, &user)
		if err != nil {
			return username, err
		}
	} else {
		return username, errors.New(fmt.Sprintf("Пользователь не сущесвует: %v", username))
	}

	return username, nil
}

func GetEnvWithDefault(key string, value string) string {
	err := viper.BindEnv(key)
	if err != nil {
		panic(fmt.Errorf("Fatal error binding env(%s): %s \n", key, err))
	}
	viper.SetDefault(key, value)
	return viper.GetString(key)
}
