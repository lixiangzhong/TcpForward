package TcpForward

import (
	"io"
	"log"
	"net"
	"time"
)

func New(local, remote string) {
	listen, err := net.Listen("tcp", local)
	if err != nil {
		log.Println(err)
		log.Println("监听", local, " 失败")
		return
	}
	log.Println("开启转发:本机", local, "->", remote)
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go pipe(conn, remote)
	}
}
func pipe(conn net.Conn, port string) {
	ser_conn, err := net.Dial("tcp", port)
	if err != nil {
		log.Println(err)
		defer conn.Close()
		return
	} else {
		log.Printf("%s访问\n%s建立TCP至%s开始转发\n\n", conn.RemoteAddr().String(),
			ser_conn.LocalAddr().String(), ser_conn.RemoteAddr().String())
		t := time.Now().Add(time.Second * 60)
		conn.SetDeadline(t)
		ser_conn.SetDeadline(t)
		go inbound(conn, ser_conn)
		outbound(conn, ser_conn)
	}
}
func inbound(in net.Conn, out net.Conn) {
	_, err := io.Copy(in, out)
	if err != nil {
		in.Close()
		out.Close()
	}
}

func outbound(in net.Conn, out net.Conn) {
	_, err := io.Copy(out, in)
	if err != nil {
		in.Close()
		out.Close()
	}
}
