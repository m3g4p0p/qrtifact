package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"m3g4p0p/qrmaster/pretty"
)

var console = pretty.NewConsole(os.Stdout)

type PingRequest struct {
	Ping bool `json:"ping"`
}

func handleConn(conn *net.UDPConn) error {
	data := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			return err
		}
		if _, _, err := conn.WriteMsgUDP(data[:n], nil, addr); err != nil {
			return err
		}
	}
}

var cmds = map[string]func() error{
	"server": func() error {
		conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 9000})
		if err != nil {
			return err
		}
		defer conn.Close()
		return handleConn(conn)
	},
	"client": func() error {
		conn, err := net.DialUDP("udp", nil, &net.UDPAddr{Port: 9000})
		if err != nil {
			return err
		}
		defer conn.Close()

		data, err := json.Marshal(PingRequest{Ping: true})
		if err != nil {
			return err
		}
		if _, err := conn.Write(data); err != nil {
			return err
		}

		res := make([]byte, 1024)

		if n, err := conn.Read(res); err != nil {
			return err
		} else {
			fmt.Println(string(res[:n]))
		}

		return nil
	},
}

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(console, nil)))
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalln("missing command")
	}

	if err := cmds[flag.Arg(0)](); err != nil {
		log.Fatalln(err)
	}
}
