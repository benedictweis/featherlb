package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Host     string    `yaml:"host"`
	Port     int       `yaml:"port"`
	Backends []Backend `yaml:"backends"`
}

type Backend struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		slog.Error("Failed to read config file", "error", err)
		return
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		slog.Error("Failed to unmarshal config file", "error", err)
		return
	}

	slog.Info("Config loaded", "config", config)

	wg := sync.WaitGroup{}

	for _, route := range config.Routes {
		wg.Add(1)
		go listenOnRoute(route)
	}

	wg.Wait()
}

func listenOnRoute(route Route) {
	address := net.JoinHostPort(route.Host, strconv.Itoa(route.Port))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("Failed to bind to address", "error", err)
		return
	}
	defer listener.Close()

	slog.Info("featherlb listening", "local_addr", listener.Addr())

	index := uint(0)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "error", err)
			continue
		}

		slog.Info("New connection", "remote_addr", clientConn.RemoteAddr())

		backend := getNextBackend(index, route.Backends)
		index++
		backendConn, err := net.Dial("tcp", net.JoinHostPort(backend.Host, strconv.Itoa(backend.Port)))
		if err != nil {
			slog.Error("Failed to connect to backend", "error", err)
			clientConn.Close()
			continue
		}

		go handleConnection(clientConn, backendConn)
	}
}

func getNextBackend(index uint, backends []Backend) Backend {
	return backends[index%uint(len(backends))]
}

func handleConnection(clientConn, backendConn net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		slog.Info("Starting to copy from client to backend")
		copyFromTo(clientConn, backendConn)
	}()

	go func() {
		defer wg.Done()
		slog.Info("Starting to copy from backend to client")
		copyFromTo(backendConn, clientConn)
	}()

	wg.Wait()
}

func copyFromTo(dst net.Conn, src net.Conn) {
	defer dst.Close()
	defer src.Close()
	if _, err := io.Copy(dst, src); err != nil {
		if errors.Is(err, net.ErrClosed) {
			slog.Info("Graceful shutdown: connection closed")
		} else {
			slog.Error("Error copying data", "error", err)
		}
	}
}
