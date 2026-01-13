package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"zinx/ziface"
)


type GlobalObj struct {
	TcpServer ziface.Iserver //当前zinx全局的srever对象
	Host string //当前服务器主机监听的ip 
	TcpPort int //当前服务器主机监听的端口
	Name string //服务器的名称

	Version string //当前zinx的版本号
	MaxConn int //最大连接数
	MaxPackageSize uint32 //zinx框架数据包的最大值
	WorkPoolSize uint32 //zinx框架工作池的worker数量
	MaxWorkTaskLen uint32
}

var GlobalObject *GlobalObj

//加载文件方法
func (G *GlobalObj) Reload() {
	data,err := os.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("readfile err:",err)
		return
	}

	//将json文件解析到struct中
	err = json.Unmarshal(data,&GlobalObject)
	if err != nil {
		panic(err)
	}
}

//提供init方法 用来初始化Globalobject
func init() {
	GlobalObject = &GlobalObj{
		Name: "zinxApp",
		Version: "v0.3",
		TcpPort: 8887,
		Host: "127.0.0.1",
		MaxConn: 1000,
		MaxPackageSize: 4096,
		WorkPoolSize: 10,
		MaxWorkTaskLen: 1024,
	}

	//尝试从conf/zinx.json去加载用户自定义的参数
	GlobalObject.Reload()
}