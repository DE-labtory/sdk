package sdk

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/it-chain/sdk/pb"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	addr       net.TCPAddr
	ctx        context.Context
	cancel     context.CancelFunc
	handler    func(request *pb.Request) (*pb.Response, error)
}

func NewServer(port int) *Server {
	addr := net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: port,
		Zone: "",
	}
	ctx, cf := context.WithCancel(context.Background())
	return &Server{
		addr:   addr,
		ctx:    ctx,
		cancel: cf,
	}
}
func (s *Server) SetHandler(handler func(request *pb.Request) (*pb.Response, error)) {
	s.handler = handler
}

func (s *Server) Listen(waitTimeOutSec int) error {
	l, err := net.ListenTCP("tcp", &s.addr)
	if err != nil {
		return err
	}
	l.SetDeadline(time.Now().Add(time.Duration(waitTimeOutSec) * time.Second))
	if err != nil {
		return err
	}
	s.grpcServer = grpc.NewServer()
	pb.RegisterBistreamServiceServer(s.grpcServer, s)
	return s.grpcServer.Serve(l)
}

func (s *Server) RunICode(streamServer pb.BistreamService_RunICodeServer) error {
	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case <-streamServer.Context().Done():
			return streamServer.Context().Err()
		default:
		}
		req, err := streamServer.Recv()
		if err == io.EOF {
			// return will close stream from server side
			return nil
		}
		if err != nil {
			return err
		}
		res, err := s.handler(req)
		if err != nil {
			return err
		}
		streamServer.Send(res)
		if err != nil {
			return err
		}
	}
}

func (s *Server) Ping(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Server) Close() {
	s.cancel()
}
