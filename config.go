package main

type Config struct {
	//#应用服务地址，省略端口号
	//AppServer = "https://app.wildfirechat.net"
	//ImServerHost = "im.wildfirechat.net"
	//ImServerNodes = ["node1.im.wildfirechat.net", "node2.im.wildfirechat.net", "node3.im.wildfirechat.net"]
	//RoutePort = 80
	//LongLinkPort = 1883
	//# web
	//EnableWeb = false
	//UseWSS = false
	//# UseWSS为false时，默认是80；UseWSS为true时，默认是443
	//WebRoutePort = 80
	//# UseWSS为false时，默认是8083；为true时，默认是8084
	//WsPort = 8083
	//
	//#备选地址配置
	//EnableBackupHost = false
	//BackupImServerHost = "im.wildfirechat.net"
	//BackupImServerNodes = ["bk1.im.wildfirechat.net"]
	//BackupRoutePort = 80
	//BackupLongLinkPort = 1883
	//
	//# 备选地址是否支持web端
	//EnableBackupHostWeb = false
	//BackupUseWSS = false
	//BackupWebRoutePort = 80
	//BackupWsPort = 8084


	AppServer string
	ImServerHost string
	ImServerNodes []string
	RoutePort int
	LongLinkPort int

	EnableWeb bool
	UseWSS bool
	WebRoutePort int
	WsPort int

	EnableBackupHost bool
	BackupImServerHost string
	BackupImServerNodes []string
	BackupRoutePort int
	BackupLongLinkPort int

	EnableBackupHostWeb bool
	BackupUseWSS bool
	BackupWebRoutePort int
	BackupWsPort int

}
