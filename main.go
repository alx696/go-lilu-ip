package main

import (
	"flag"
	"github.com/libp2p/go-nat"
	"log"
	"net/http"
	"net/url"
	"time"
)

var cacheIp string

func main() {
	serverFlag := flag.String("server", "https://ip.dev.lilu.red/", "IP Server")
	idFlag := flag.String("id", "test", "Your ID")
	flag.Parse()

	server := *serverFlag
	id := *idFlag
	log.Println("服务:", server, "ID:", id)
	log.Println("正在检测公网IP")
	natGateway, e := nat.DiscoverGateway()
	for e != nil {
		log.Println("寻找NAT网关错误:", e)

		time.Sleep(time.Minute)
		natGateway, e = nat.DiscoverGateway()
	}

	for {
		natExternalAddress, e := natGateway.GetExternalAddress()
		if e != nil {
			log.Println("获取公网IP错误:", e)

			time.Sleep(time.Minute)
			continue
		}
		ip := natExternalAddress.String()
		log.Println("公网IP:", ip)

		if cacheIp != ip {
			log.Println("更新公网IP")
			formData := url.Values{}
			formData.Set("id", id)
			formData.Set("ip", ip)
			resp, e := http.PostForm(server, formData)
			if e != nil {
				log.Println("更新公网IP错误:", e)

				time.Sleep(time.Minute)
				continue
			}
			log.Println("更新公网IP结果:", resp.StatusCode)

			cacheIp = ip
		}

		time.Sleep(time.Second)
	}
}
