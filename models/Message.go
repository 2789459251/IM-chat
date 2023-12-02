package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FromId   int64  //发送者
	TargetId int64  //接受者
	Type     int    //类型：群聊 私聊 广播
	Media    int    //消息类型：文字，图片，音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn       *websocket.Conn
	DataQueue  chan []byte
	GroupeSets set.Interface
}

// 初始化node,放映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(w http.ResponseWriter, r *http.Request) {
	//Tudo 校验合法性
	//token := query.Get("token")
	query := r.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//msgtype := query.Get("type")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	isvalida := true //checktoken 待补充
	conn, err := (&websocket.Upgrader{
		//token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2.获取连接
	node := &Node{
		Conn:       conn,
		DataQueue:  make(chan []byte, 50),
		GroupeSets: set.New(set.ThreadSafe),
	}
	//3.用户关系
	//4.userid 与node绑定枷锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送
	go sendProc(node)
	//6.接收逻辑
	go recvProc(node)

	sendMsg(userId, []byte("欢迎进入chat"))
}
func sendProc(node *Node) {
	for {

		select {

		case data := <-node.DataQueue:
			//fmt.Println("[ws]sendMsg >>>>", "msg:", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
func recvProc(node *Node) {
	for {
		_, date, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(date)
		//fmt.Println("[ws] recvProc<<<<<", string(date))
	}
}

var upsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(date []byte) {
	upsendChan <- date
}
func init() {
	go udpSendProc()
	go udpRecvProc()
	//fmt.Println("go inited")
}

// 完成udp数据发送协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 135),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case data := <-upsendChan:
			//fmt.Println("udpSendProc:", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 接受
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])

		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println("udpRecvProc data: ", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

// 后端调度
func dispatch(date []byte) {
	msg := Message{}
	err := json.Unmarshal(date, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: //私信
		//fmt.Println("udpRecvMsg:dispatch:", string(date))
		sendMsg(msg.TargetId, date)
		//case 2://群发
		//sendGroupMsg()
		//case 3://广播
		//sendAllMsg()
		//case 4:
		//
	}
}
func sendMsg(targetId int64, msg []byte) {
	//fmt.Println("sendMsg>>>targetId", targetId, string(msg))
	rwLocker.RLock()
	node, OK := clientMap[targetId]
	rwLocker.RUnlock()
	if OK {
		node.DataQueue <- msg
	}
}
