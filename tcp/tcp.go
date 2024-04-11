package tcp

import (
    "io"
    "log"
    "net"
    "fmt"
    "sync"
    "strings"
)

type tcp struct {
    lock sync.Mutex
    src  string
    dsts []string
}

func New(src, dst string) *tcp {
    tp := &tcp{
        src:  src,
        dsts: strings.Split(dst, ","),
    }

    return tp
}

func (p *tcp) Server() {
    listen, err := net.Listen("tcp", p.src)
    if err != nil {
        fmt.Println(err)
        return
    }

    defer listen.Close()

    for {
        conn, err := listen.Accept()
        if err != nil {
            log.Printf("client: %s, local: %s, listen accept error: %v\n", conn.RemoteAddr(), conn.LocalAddr(), err)
            continue
        }

        go p.handle(conn)
    }
}

func (p *tcp) handle(sconn net.Conn) {
    defer sconn.Close()
    dst, ok := p.selectDst()
    if !ok {
        return
    }

    dconn, err := net.Dial("tcp", dst)
    if err != nil {
        log.Printf("dial %v fail: %v\n", dst, err)
        return
    }
    defer dconn.Close()

    ExitChan := make(chan bool, 1)

    // 转发到目标服务器
    go func(sconn net.Conn, dconn net.Conn, exit chan bool) {
        _, err := io.Copy(dconn, sconn)
        if err != nil {
            log.Printf("give message to %v fail: %v\n", dst, err)
            ExitChan <- true
        }
    }(sconn, dconn, ExitChan)

    // 从目标服务器返回数据到客户端
    go func(sconn net.Conn, dconn net.Conn, exit chan bool) {
        _, err := io.Copy(sconn, dconn)
        if err != nil {
            log.Printf("get message fail from %v: %v\n", dst, err)
            ExitChan <- true
        }
    }(sconn, dconn, ExitChan)

    <-ExitChan
}

// ip 轮询
func (p *tcp) selectDst() (string, bool) {
    p.lock.Lock()
    defer p.lock.Unlock()

    if len(p.dsts) < 1 {
        return "", false
    }

    dst := p.dsts[0]
    p.dsts = append(p.dsts[1:], dst)
    return dst, true
}
