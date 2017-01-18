package server

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	pb "github.com/trumanw/findpro/pb"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
)

// Run is the main function of gRPC server
func Run(host string, port int, etcdns []string) error {
    // Bind gRPC server to specific host and port
    addr := host + ":" + strconv.Itoa(port)
    l, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }

    // register the service endpoint to etcd
    cli, cerr := clientv3.New(clientv3.Config{
		Endpoints:   etcdns,
		DialTimeout: 5 * time.Second,
	})
    if cerr != nil {
        return cerr
    }
    defer cli.Close()

    r := &etcdnaming.GRPCResolver{Client: cli}
    r.Update(context.TODO(), "findpro", naming.Update{Op: naming.Add, Addr: addr, Metadata: "..."})

    // Init gRPC interceptors
    s := grpc.NewServer(
    	grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
    	grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)

    // Init gRPC services
    pb.RegisterFingerprintServiceServer(s, NewFingerprintGRPCServer())

    // After all your registrations, make sure all of the Prometheus metrics are initialized.
	grpc_prometheus.Register(s)
    // Register Prometheus metrics handler.
	http.Handle("/metrics", prometheus.Handler())

    // Launch gRPC server
    s.Serve(l)
    return nil
}
