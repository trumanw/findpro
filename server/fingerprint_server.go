package server

import (
    "golang.org/x/net/context"

    pb "github.com/trumanw/findpro/pb"
)

// FingerprintGRPCServer is the object of GRPC server
type FingerprintGRPCServer struct{}

// NewFingerprintGRPCServer returns an instance of GRPC server
func NewFingerprintGRPCServer() pb.FingerprintServiceServer {
    return new(FingerprintGRPCServer)
}

// Learn is the implementation of the api '/learn'
func (s *FingerprintGRPCServer) Learn(ctx context.Context, fingerprint *pb.Fingerprint) (*pb.Result, error) {
    learningResult := &pb.Result{
        Message: "Could not learn.",
        Success: false,
        Location: nil,
    }
    return learningResult, nil
}

// Track is the implementation of the api '/track'
func (s *FingerprintGRPCServer) Track(ctx context.Context, fingerprint *pb.Fingerprint) (*pb.Result, error) {
    trackingResult := &pb.Result{
        Message: "Could not track.",
        Success: false,
        Location: nil,
    }
    return trackingResult, nil
}
