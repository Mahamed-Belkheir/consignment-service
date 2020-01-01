package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/MohamedBelkheirRBK/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

type service struct {
	repo Repository
}

func (s *service) CreateConsignment(ctx *context.Context, req *pb.Consignment) (*pb.Reponse, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return nil, &pb.Reponse{Created: true, Consignment: consignment}
}

func main() {

	repo := &Repository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	pb.RegisterShippingServiceServer(server, &server{repo})

	reflection.Register(server)

	log.Println("Running on port:", port)

	if err := server.serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
