package handlers

import pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"

type Server struct {
	pb.UnimplementedExecsServiceServer
	pb.UnimplementedStudentsServiceServer
	pb.UnimplementedTeachersServiceServer
}
