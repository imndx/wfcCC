package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {

	conf, err := loadConfig()

	if err != nil {
		fmt.Println("加载配置文件错误", err.Error())
		return
	}

	// 检查基本配置
	checkAppServer(conf.AppServer)
	checkIMServerVersion("http://" + conf.ImServerHost + "/api/version")
	//checkIMServerTCPPort(conf.ImServerHost, conf.LongLinkPort)
	for _, node := range conf.ImServerNodes {
		checkIMServerVersion("http://" + node + "/api/version")
		checkIMServerTCPPort(node, conf.LongLinkPort)
	}

	// 检查web相关配置
	if conf.EnableWeb {
		var routeUrl string
		if conf.UseWSS {
			routeUrl = fmt.Sprintf("%s%s:%d%s", "https://", conf.ImServerHost, conf.WebRoutePort, "/route")
		} else {
			routeUrl = fmt.Sprintf("%s%s:%d%s", "http://", conf.ImServerHost, conf.WebRoutePort, "/route")
		}
		checkIMServerRouteCors(routeUrl)
		for _, node := range conf.ImServerNodes {
			checkIMServerTCPPort(node, conf.WsPort)
		}	}

	// 检查备选地址相关配置
	if conf.EnableBackupHost {
		checkIMServerVersion("http://" + conf.BackupImServerHost+ "/api/version")
		checkIMServerTCPPort(conf.BackupImServerHost, conf.BackupLongLinkPort)
		for _, node := range conf.BackupImServerNodes {
			checkIMServerVersion("http://" + node + "/api/version")
			checkIMServerTCPPort(node, conf.LongLinkPort)
		}
	}

	if conf.EnableBackupHostWeb {
		var routeUrl string
		if conf.BackupUseWSS {
			routeUrl = fmt.Sprintf("%s%s:%d", "https://", conf.BackupImServerHost, conf.BackupWebRoutePort)
		} else {
			routeUrl = fmt.Sprintf("%s%s:%d", "http://", conf.BackupImServerHost, conf.BackupWebRoutePort)
		}
		checkIMServerRouteCors(routeUrl)
		checkIMServerTCPPort(conf.BackupImServerHost, conf.BackupWsPort)
	}

}

func loadConfig() (*Config, error) {
	conf := &Config{}
	_, err := toml.DecodeFile("conf.toml", conf)
	return conf, err
}

func checkAppServer(addr string) bool {
	resp, err := get(addr)

	if err != nil {
		fmt.Println("check app server error", err.Error())
		return false
	}
	if strings.Compare(resp, "Ok") == 0 {
		fmt.Println("检查App Server成功", addr)
		return true
	}
	fmt.Println("app server response error", resp)
	return false
}

func checkIMServerVersion(addr string) bool {
	resp, err := get(addr)
	if err != nil {
		fmt.Println("检查IM Server版本号失败", err.Error())
		return false
	}
	if strings.Index(resp, "wfmaster") >= 0 {
		fmt.Println("检查IM Server版本号成功", addr)
		return true
	}
	fmt.Sprintln("检查IM Server版本号失败，返回信息未：", resp)
	return false
}

func checkIMServerRouteCors(addr string) bool {
	routeUrl, _ := url.Parse(addr)
	req := &http.Request{
		Method: "OPTIONS",
		URL:    routeUrl,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("检查IM Server route接口跨域失败", addr, err.Error())
		return false
	}
	values := resp.Header.Values("Access-Control-Allow-Origin")
	if len(values) != 1 || strings.Compare("*", values[0]) != 0 {
		fmt.Println("IM Server route接口Access-Control-Allow-Origin配置不对", addr, values)
		return false
	}
	values = resp.Header.Values("Access-Control-Allow-Headers")
	for i := 0; i < len(values); i++ {
		if strings.Index(values[i], "p,uid,cid,appId,appKey") >= 0 {
			fmt.Println("检查IM Server route接口跨域成功")
			return true
		}
	}
	fmt.Println("IM Server route接口Access-Control-Allow-Headers配置不对", values)

	return false
}

func checkIMServerTCPPort(host string, port int) bool {
	remoteAddr, err := tcp(host, port)
	if err != nil {
		fmt.Println("检查IM Server TCP端口错误", host, port, err.Error())
		return false
	}
	fmt.Println("检查IM Server TCP端口成功", host, port, remoteAddr)

	return true
}

func tcp(host string, port int) (string, error) {
	timeout := time.Duration(60) * time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return "", err
	}
	remoteAddr := conn.RemoteAddr()
	defer conn.Close()
	return remoteAddr.String(), nil
}

func get(addr string) (string, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
