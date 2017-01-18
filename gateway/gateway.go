package gateway

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
	idem "github.com/trumanw/findpro/gateway/handler/idempotent"
	log "github.com/trumanw/findpro/gateway/handler/logrus"
	pb "github.com/trumanw/findpro/pb"
	ng "github.com/urfave/negroni"

	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
)

// Run is the main function to launch a gRPC gateway
func Run(etcdns []string) error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()
    // Add middlewares
    n := ng.New()
    n.Use(log.NewDefultMiddleware())
    n.Use(idem.NewDefaultMiddeleware())
	n.Use(cors.New(cors.Options{}))
	n.Use(gzip.Gzip(gzip.DefaultCompression))
    n.UseHandler(mux)

    // Resolve gRPC servers' connections through etcd
    cli, cerr := clientv3.NewFromURL("http://localhost:2379")
    if cerr != nil {
        return cerr
    }
    defer cli.Close()
    r := &etcdnaming.GRPCResolver{Client: cli}
    b := grpc.RoundRobin(r)
    conn, gerr := grpc.Dial("findpro", grpc.WithInsecure(), grpc.WithBalancer(b))
    if gerr != nil {
        return gerr
    }

    // Register client with the gRPC servers' connections
    err := newGW(ctx, mux, conn)
    if err != nil {
        return err
    }

    http.ListenAndServe(":8080", n)
    return nil
}

// init retrieves the connections through
func newGW(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
    err := pb.RegisterFingerprintServiceHandler(ctx, mux, conn)
    if err != nil {
        return err
    }
    return nil
}
