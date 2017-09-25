package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/flatbuffers/tests/MyGame/Example"
)

const (
	addr = "0.0.0.0:50051"
)

var (
	test = "Flatbuffers"
)

type server struct{}

/*
	Store(context.Context, *Monster) (*flatbuffers.Builder, error)
	Retrieve(*Stat, MonsterStorage_RetrieveServer) error
*/

func (s *server) Store(context context.Context, in *Example.Monster) (*flatbuffers.Builder, error) {
	b := flatbuffers.NewBuilder(0)
	i := b.CreateString(test)
	Example.StatStart(b)
	Example.StatAddId(b, i)
	b.Finish(Example.StatEnd(b))
	return b, nil

}

func (s *server) Retrieve(in *Example.Stat, stream Example.MonsterStorage_RetrieveServer) error {
	/*
		b := flatbuffers.NewBuilder(0)
		i := b.CreateString(test)
		Example.MonsterStart(b)
		Example.MonsterAddName(b, i)
		b.Finish(Example.MonsterEnd(b))
	*/
	return nil
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.CustomCodec(flatbuffers.FlatbuffersCodec{}))
	Example.RegisterMonsterStorageServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
