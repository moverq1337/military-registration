package delivery

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcLib "github.com/neiasit/grpc-library/core"
	"github.com/neiasit/service-boilerplate/internal/user/dto"
	"github.com/neiasit/service-boilerplate/internal/user/usecase"
	doctorgen "github.com/neiasit/service-boilerplate/pkg/api/grpc/golang/doctor"
	usergen "github.com/neiasit/service-boilerplate/pkg/api/grpc/golang/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type UserHandlers struct {
	logger      *slog.Logger
	userUsecase usecase.UserUsecase
	usergen.UnimplementedUserServiceServer
}

func NewUserHandlers(logger *slog.Logger, userUsecase usecase.UserUsecase) *UserHandlers {
	return &UserHandlers{logger: logger, userUsecase: userUsecase}
}

func RegisterHandlers(
	grpcSrv *grpc.Server,
	grpcGtw *runtime.ServeMux,
	grpcConfig *grpcLib.Config,
	logger *slog.Logger,
	userUsecase usecase.UserUsecase,
) {

	ctx := context.Background()

	impl := NewUserHandlers(logger, userUsecase)
	usergen.RegisterUserServiceServer(grpcSrv, impl)
	err := usergen.RegisterUserServiceHandlerFromEndpoint(ctx, grpcGtw, grpcConfig.Address(), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		logger.Error("failed to register example service handler", "error", err)
		return
	}
	ctx = context.WithValue(ctx, "parameter", "im fucking parametr")
	slog.InfoContext(ctx, "hello world")
}

func (h *UserHandlers) CreateUser(ctx context.Context, rq *usergen.CreateUserRq) (*emptypb.Empty, error) {
	h.logger.Info("Name And Surname", "name", rq.Name, "surname", rq.Surname)
	if err := h.userUsecase.Create(ctx, &dto.CreateUserRq{
		Name:    rq.Name,
		Surname: rq.Surname,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *UserHandlers) GetDoctors(ctx context.Context, rq *usergen.GetDoctorsRq) (*usergen.GetDoctorsRs, error) {
	doctors, err := h.userUsecase.GetDoctors(ctx, rq.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var doctorsDto []*doctorgen.DoctorItem
	for _, doctor := range doctors {
		doctorsDto = append(doctorsDto, &doctorgen.DoctorItem{
			Id:      doctor.Id,
			Name:    doctor.Name,
			Surname: doctor.Surname,
			Cabinet: doctor.Cabinet,
			Type:    doctor.Type,
		})
	}

	return &usergen.GetDoctorsRs{
		Doctors: doctorsDto,
	}, nil
}
