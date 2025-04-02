package server

import (
	"featherlb/internal/pkg/strategies"
	"featherlb/internal/pkg/types"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"sync"
	"time"
)

type FeatherLBServer struct {
	bufPool        sync.Pool
	connectionPool map[string]*sync.Pool
}

// NewFeatherLBServer creates a new instance of FeatherLBServer.
func NewFeatherLBServer() *FeatherLBServer {
	return &FeatherLBServer{
		bufPool:        sync.Pool{New: func() interface{} { return make([]byte, 32*1024) }},
		connectionPool: make(map[string]*sync.Pool),
	}
}

// StartWithConfig initializes the server with the given configuration and starts listening for incoming connections.
func (s *FeatherLBServer) StartWithConfig(config types.Config) {
	// Initialize connection pools for each endpoint
	for _, route := range config.Routes {
		for _, endpoint := range route.Endpoints {
			key := net.JoinHostPort(endpoint.Host, strconv.Itoa(endpoint.Port))
			s.connectionPool[key] = &sync.Pool{
				New: func() interface{} {
					connection, err := net.DialTimeout("tcp", key, 3*time.Second)
					if err != nil {
						slog.Warn("Error creating endpoint connection", "error", err)
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
		go s.listenOnRoute(route)
	}

	wg.Wait()
}

// listenOnRoute listens for incoming connections on the specified route and handles them.
func (s *FeatherLBServer) listenOnRoute(route types.Route) {
	address := net.JoinHostPort(route.Host, strconv.Itoa(route.Port))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("Failed to bind to address", "error", err)
		return
	}
	defer listener.Close()

	slog.Info("featherlb listening", "local_addr", listener.Addr())

	strategy := strategies.MatchStrategy(route.Strategy)
	for _, endpoint := range route.Endpoints {
		strategy.AddEndpoint(endpoint)
	}

	slog.Info("initialized strategy", "strategy", route.Strategy)

	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", "error", err)
			continue
		}

		slog.Info("New connection", "remote_addr", clientConnection.RemoteAddr())

		go s.handleConnection(clientConnection, strategy)
	}
}

// handleConnection handles the data transfer between the client and the selected endpoint.
func (s *FeatherLBServer) handleConnection(clientConnection net.Conn, strategy strategies.Strategy) {
	defer clientConnection.Close()

	endpoint, err := strategy.Next(*clientConnection.RemoteAddr().(*net.TCPAddr))
	if err != nil {
		slog.Error("Failed to get endpoint", "error", err)
		return
	}
	endpointConnection, err := net.Dial("tcp", net.JoinHostPort(endpoint.Host, strconv.Itoa(endpoint.Port)))
	if err != nil {
		slog.Error("Failed to connect to endpoint", "error", err)
		return

	}

	endpointHostPort := net.JoinHostPort(endpoint.Host, strconv.Itoa(endpoint.Port))

	serverConn := s.connectionPool[endpointHostPort].Get()

	if serverConn == nil {
		var err error
		serverConn, err = net.DialTimeout("tcp", endpointHostPort, 3*time.Second)
		if err != nil {
			fmt.Println("Failed to connect to endpoint:", err)
			return
		}
	}

	serverConnection := serverConn.(net.Conn)

	defer s.connectionPool[endpointHostPort].Put(serverConnection)

	var wg sync.WaitGroup

	wg.Add(2)

	go s.pipeData(clientConnection, endpointConnection, &wg)
	go s.pipeData(endpointConnection, clientConnection, &wg)

	wg.Wait()
}

// pipeData copies data from src to dst using a buffer from the pool.
func (s *FeatherLBServer) pipeData(src, dst net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := s.bufPool.Get().([]byte)
	defer s.bufPool.Put(buf)
	io.CopyBuffer(dst, src, buf)
}
