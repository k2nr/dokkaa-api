package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func getopt(name, def string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return def
}

type Service struct {
	Port int    `json:"port"`
	Role string `json:"role"`
}

type Container struct {
	Image    string             `json:"image"`
	Name     string             `json:"name"`
	Scale    int                `json:"scale"`
	Env      map[string]string  `json:"env"`
	Links    []string           `json:"links"`
	Command  []string           `json:"command"`
	Services map[string]Service `json:"services"`
}

type Spec struct {
	AppName    string      `json:"app_name" binding:"required"`
	Containers []Container `json:"containers"`
}

func main() {
	router := gin.Default()
	etcdAddr := "http://" + getopt("ETCD_ADDR", "127.0.0.1:4001")

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("/apps/:app", func(c *gin.Context) {
		app := c.Params.ByName("app")
		ec := etcd.NewClient([]string{etcdAddr})
		v, err := ec.Get("/apps/"+app, true, false)
		if err == nil {
			log.Println(v)
		} else {
			log.Println(err)
		}
		c.JSON(200, v.Node)
	})

	router.POST("/apps", func(c *gin.Context) {
		ec := etcd.NewClient([]string{etcdAddr})
		var spec Spec
		// TODO: use Bind
		c.BindWith(&spec, binding.JSON)
		keyBase := "/apps/" + spec.AppName + "/"
		for _, container := range spec.Containers {
			if container.Name == "" {
				continue
			}
			key := keyBase + container.Name + "/manifest"
			v, err := json.Marshal(&container)
			if err != nil {
				log.Println(err)
				continue
			}

			_, err = ec.Set(key, string(v), 0)
			if err != nil {
				log.Println(err)
			}
		}
	})

	port := getopt("PORT", "8080")
	router.Run(":" + port)
}
