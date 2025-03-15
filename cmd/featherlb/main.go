package main

import (
	"featherlb/cmd/featherlb/strategies"
	"featherlb/cmd/featherlb/types"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"sync"
	"time"
)

var bufPool = sync.Pool{New: func() interface{} { return make([]byte, 32*1024) }}
var connectionPool = make(map[string]*sync.Pool)

func main() {
	configPath := flag.String("config", "", "Path to the config file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()
	if *configPath == "" {
		slog.Error("Config file path is required")
		return
	}

	configureLogging(*debug)

	config, err := types.ReadConfigFromFile(*configPath)
	if err != nil {
		slog.Error("Failed to read config file", "error", err)
		return
	}

	slog.Debug("Config loaded", "location", *configPath, "config", config)

	// Initialize connection pools for each backend
	for _, route := range config.Routes {
		for _, backend := range route.Backends {
			key := net.JoinHostPort(backend.Host, strconv.Itoa(backend.Port))
			connectionPool[key] = &sync.Pool{
				New: func() interface{} {
					connection, err := net.DialTimeout("tcp", key, 3*time.Second)
					if err != nil {
						slog.Warn("Error creating backend connection", "error", err)
						return nil
					}
					return connection
				},
			}
		}
	}

	wg := sync.WaitGroup{}

	for _, route := range config.Routes {
		wg.Add(1)
		go listenOnRoute(route)
	}

	wg.Wait()
}

func listenOnRoute(route types.Route) {
	address := net.JoinHostPort(route.Host, strconv.Itoa(route.Port))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("Failed to bind to address", "error", err)
		return
	}
	defer listener.Close()

	slog.Info("featherlb listening", "local_addr", listener.Addr())

	strategy := strategies.NewRoundRobinStrategy()
	for _, backend := range route.Backends {
		strategy.AddBackend(backend)
	}

	slog.Info("initialized strategy", "strategy", "round-robin")

	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "error", err)
			continue
		}

		slog.Info("New connection", "remote_addr", clientConnection.RemoteAddr())

		go handleConnection(clientConnection, strategy)
	}
}

func handleConnection(clientConnection net.Conn, strategy strategies.Strategy) {
	defer clientConnection.Close()

	backend, err := strategy.Next()
	if err != nil {
		slog.Error("Failed to get backend", "error", err)
		return
	}
	backendConnection, err := net.Dial("tcp", net.JoinHostPort(backend.Host, strconv.Itoa(backend.Port)))
	if err != nil {
		slog.Error("Failed to connect to backend", "error", err)
		return

	}

	backendHostPort := net.JoinHostPort(backend.Host, strconv.Itoa(backend.Port))

	serverConn := connectionPool[backendHostPort].Get()

	if serverConn == nil {
		var err error
		serverConn, err = net.DialTimeout("tcp", backendHostPort, 3*time.Second)
		if err != nil {
			fmt.Println("Failed to connect to backend:", err)
			return
		}
	}

	serverConnection := serverConn.(net.Conn)

	defer connectionPool[backendHostPort].Put(serverConnection)

	var wg sync.WaitGroup

	wg.Add(2)

	go pipeData(clientConnection, backendConnection, &wg)
	go pipeData(backendConnection, clientConnection, &wg)

	wg.Wait()
}

func pipeData(src, dst net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)
	io.CopyBuffer(dst, src, buf)
}
