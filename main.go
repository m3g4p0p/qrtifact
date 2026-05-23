package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"log/slog"
	"os"

	"m3g4p0p/qrtifact/pretty"

	"github.com/joho/godotenv"
	"golang.org/x/exp/jsonrpc2"
)

var console = pretty.NewConsole(os.Stdout)

var binder = &jsonrpc2.ConnectionOptions{
	Handler: jsonrpc2.HandlerFunc(Handle),
}

type PingParams struct {
	Ping bool `json:"ping"`
}

type PingResult struct {
	Pong bool `json:"pong"`
}

func Handle(ctx context.Context, req *jsonrpc2.Request) (any, error) {
	if req.Method != "ping" {
		return nil, jsonrpc2.ErrMethodNotFound
	}
	var params PingParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return nil, jsonrpc2.ErrInvalidParams
	}
	return PingResult{Pong: params.Ping}, nil
}

var cmds = map[string]func() error{
	"server": func() error {
		return nil
	},
	"client": func() error {
		ctx := context.Background()

		listener, err := jsonrpc2.NetListener(ctx, "tcp", ":9090", jsonrpc2.NetListenOptions{})
		if err != nil {
			return err
		}
		defer listener.Close()

		server, err := jsonrpc2.Serve(ctx, listener, binder)
		if err != nil {
			return err
		}

		conn, err := jsonrpc2.Dial(ctx, listener.Dialer(), binder)
		if err != nil {
			return err
		}
		defer conn.Close()

		params := PingParams{Ping: true}
		call := conn.Call(ctx, "ping", params)

		var result PingResult
		if err := call.Await(ctx, &result); err != nil {
			return err
		}
		console.Pretty(result)

		return server.Wait()
	},
}

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(console, nil)))

	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
	key, err := hex.DecodeString(os.Getenv("SECRET_KEY"))
	if err != nil {
		log.Fatalln(err)
	}
	_ = key
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalln("msiing command")
	}

	if err := cmds[flag.Arg(0)](); err != nil {
		log.Fatalln(err)
	}
}
