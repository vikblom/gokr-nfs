package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-git/go-billy/v5/osfs"
	nfs "github.com/willscott/go-nfs"
	nfshelper "github.com/willscott/go-nfs/helpers"
)

var directory = os.Args[1]

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()

	go func() {
		server := &http.Server{Addr: ":8081", Handler: http.HandlerFunc(handleGetDevices)}
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err)
		}
	}()

	err := run(ctx)
	if err != nil {
		log.Fatal("run", err)
	}
	os.Exit(0)
}

func run(_ context.Context) error {
	listener, err := net.Listen("tcp", ":2049")
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	fmt.Printf("Server running at %s hosting %s\n", listener.Addr(), directory)

	bfs := osfs.New(directory)

	handler := nfshelper.NewNullAuthHandler(bfs)
	cacheHelper := nfshelper.NewCachingHandler(handler, 1024)
	err = nfs.Serve(listener, cacheHelper)
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	fmt.Println("done")
	return nil
}

func handleGetDevices(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("read dir: %s", err)
	} else {
		for _, e := range entries {
			fmt.Fprintf(w, "%s\n", e.Name())
		}
	}
}
