package fixture

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/reflection"
	reflectionv1alpha "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type PathResponseHandler struct {
	Path            string
	PathHandlerFunc func(rw http.ResponseWriter, r *http.Request)
}

func StartGrpcTargetTestServer(callStats *CallStats) (*grpc.Server, int) {
	return startGrpcTargetTestServer(callStats, func(server *grpc.Server) {
		reflection.Register(server)
	})
}

func StartGrpcTargetTestServerReflectionV1(callStats *CallStats) (*grpc.Server, int) {
	return startGrpcTargetTestServer(callStats, func(server *grpc.Server) {
		reflection.RegisterV1(server)
	})
}

func StartGrpcTargetTestServerReflectionV1Alpha(callStats *CallStats) (*grpc.Server, int) {
	return startGrpcTargetTestServer(callStats, func(server *grpc.Server) {
		reflectionServer := reflection.NewServer(reflection.ServerOptions{Services: server})
		reflectionv1alpha.RegisterServerReflectionServer(server, reflectionServer)
	})
}

// It uses the test.proto from grpc-testing: https://github.com/grpc/grpc-go/blob/40a879c23a0dc77234d17e0699d074d5fd151bd0/test/grpc_testing/test.proto
func startGrpcTargetTestServer(callStats *CallStats, reflRegFunc func(*grpc.Server)) (*grpc.Server, int) {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(callStats.UnaryInterceptor()),
	)
	grpc_testing.RegisterTestServiceServer(server, &grpc_testing.UnimplementedTestServiceServer{})
	reflRegFunc(server)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.Fatal("Server failed : ", err)
		}
	}()
	return server, port
}

// StartHttpTargetTestServer starts a HTTP server on the provided port
// Optionally, it receives a list of handler functions
func StartHttpTargetTestServer(pathHandlers []PathResponseHandler) (*http.Server, int) {
	router := http.NewServeMux()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Sleep for half a second to simulate a slow server
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	for _, pathHandler := range pathHandlers {
		router.HandleFunc(pathHandler.Path, pathHandler.PathHandlerFunc)
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	addr := listener.Addr().(*net.TCPAddr).String()
	port := listener.Addr().(*net.TCPAddr).Port

	server := &http.Server{
		Handler: h2c.NewHandler(router, &http2.Server{}),
	}

	go func() {
		err := server.Serve(listener)
		fmt.Printf("Listening on [%s]...\n", addr)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed. Err: %v", err)
		}
	}()

	return server, port
}
