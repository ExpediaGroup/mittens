package fixture

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/grpc_testing"
)

type PathResponseHandler struct {
	Path            string
	PathHandlerFunc func(rw http.ResponseWriter, r *http.Request)
}

// StartGrpcTargetTestServer starts a gRPC server on the provided port
// It uses the test.proto from grpc-testing: https://github.com/grpc/grpc-go/blob/40a879c23a0dc77234d17e0699d074d5fd151bd0/test/grpc_testing/test.proto
func StartGrpcTargetTestServer(port int) *grpc.Server {
	server := grpc.NewServer()
	grpc_testing.RegisterTestServiceServer(server, &grpc_testing.UnimplementedTestServiceServer{})
	reflection.Register(server)

	uri := ":" + fmt.Sprint(port)
	l, _ := net.Listen("tcp", uri)

	go func() {
		err := server.Serve(l)
		if err != nil {
			log.Fatal("Server failed : ", err)
		}
	}()
	return server
}

// StartHttpTargetTestServer starts a HTTP server on the provided port
// Optionally, it receives a list of handler functions
func StartHttpTargetTestServer(port int, pathHandlers []PathResponseHandler) *http.Server {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/delay", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.WriteHeader(http.StatusNoContent)
	})

	for _, pathHandler := range pathHandlers {
		http.HandleFunc(pathHandler.Path, pathHandler.PathHandlerFunc)
	}

	baseURL := ":" + fmt.Sprint(port)
	server := &http.Server{Addr: baseURL}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Server failed : ", err)
		}
	}()
	return server
}
