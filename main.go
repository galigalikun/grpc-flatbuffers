package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	data "github.com/galigalikun/grpc-flatbuffers/Data"
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

type TestStorage_RetrieveServer interface {
	Send(*flatbuffers.Builder) error
	grpc.ServerStream
}

type TestStorageServer interface {
	Store(context.Context, *Example.Monster) (*flatbuffers.Builder, error)
	Retrieve(*data.User, TestStorage_RetrieveServer) error
}

func (s *server) Store(context context.Context, in *Example.Monster) (*flatbuffers.Builder, error) {
	b := flatbuffers.NewBuilder(0)
	i := b.CreateString(test)
	Example.StatStart(b)
	Example.StatAddId(b, i)
	b.Finish(Example.StatEnd(b))
	return b, nil

}

func (s *server) Retrieve(in *data.User, stream TestStorage_RetrieveServer) error {
	/*
		b := flatbuffers.NewBuilder(0)
		i := b.CreateString(test)
		Example.MonsterStart(b)
		Example.MonsterAddName(b, i)
		b.Finish(Example.MonsterEnd(b))
	*/
	return nil
}

type monsterStorageRetrieveServer struct {
	grpc.ServerStream
}

func (x *monsterStorageRetrieveServer) Send(m *flatbuffers.Builder) error {
	return x.ServerStream.SendMsg(m)
}

var _MonsterStorage_serviceDesc = grpc.ServiceDesc{
	ServiceName: "IMyFirstService",
	HandlerType: (*TestStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SumAsync",
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				fmt.Println("SumAsync!!!")
				in := new(Example.Monster)
				if err := dec(in); err != nil {
					fmt.Println("in error")
					return nil, err
				}
				if interceptor == nil {
					fmt.Printf("interceptor nil:%v\n", srv)
					return srv.(TestStorageServer).Store(ctx, in)
				}
				info := &grpc.UnaryServerInfo{
					Server:     srv,
					FullMethod: "/Example.MonsterStorage/Store",
				}

				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					return srv.(TestStorageServer).Store(ctx, req.(*Example.Monster))
				}
				return interceptor(ctx, in, info, handler)
			},
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName: "StreamingUser",
			Handler: func(srv interface{}, stream grpc.ServerStream) error {
				m := new(data.User)
				fmt.Printf("StreamingUser:%v!!!\n", stream)
				if err := stream.RecvMsg(m); err != nil {
					fmt.Printf("Streaming RecvMsg error%v\n", err)
					return err
				}
				fmt.Printf("Streaming end:%s\n", m.Name())
				return srv.(TestStorageServer).Retrieve(m, &monsterStorageRetrieveServer{stream})
			},
			ServerStreams: true,
		},
	},
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.CustomCodec(flatbuffers.FlatbuffersCodec{}))
	s.RegisterService(&_MonsterStorage_serviceDesc, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
