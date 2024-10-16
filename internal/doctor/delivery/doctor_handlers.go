package delivery

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcLib "github.com/neiasit/grpc-library/core"
	"github.com/neiasit/service-boilerplate/internal/doctor/dto"
	"github.com/neiasit/service-boilerplate/internal/doctor/usecase"
	doctorgen "github.com/neiasit/service-boilerplate/pkg/api/grpc/golang/doctor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DoctorHandlers struct {
	logger        *slog.Logger
	doctorUsecase usecase.DoctorUsecase
	doctorgen.UnimplementedDoctorServiceServer
}

func NewDoctorHandlers(logger *slog.Logger, doctorUsecase usecase.DoctorUsecase) *DoctorHandlers {
	return &DoctorHandlers{logger: logger, doctorUsecase: doctorUsecase}
}

func RegisterDoctorHandlers(
	grpcSrv *grpc.Server,
	grpcGtw *runtime.ServeMux,
	grpcConfig *grpcLib.Config,
	logger *slog.Logger,
	doctorUsecase usecase.DoctorUsecase,
) {

	ctx := context.Background()

	impl := NewDoctorHandlers(logger, doctorUsecase)
	doctorgen.RegisterDoctorServiceServer(grpcSrv, impl)
	err := doctorgen.RegisterDoctorServiceHandlerFromEndpoint(ctx, grpcGtw, grpcConfig.Address(), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		logger.Error("failed to register example service handler", "error", err)
		return
	}
	ctx = context.WithValue(ctx, "parameter", "im fucking parametr")
	slog.InfoContext(ctx, "hello world")
}

func (h *DoctorHandlers) Create(ctx context.Context, rq *doctorgen.CreateRq) (*emptypb.Empty, error) {
	createDto := dto.CreateDoctorRs{
		Name:    rq.Name,
		Surname: rq.Surname,
		Cabinet: rq.Cabinet,
		Type:    rq.Type,
	}
	err := h.doctorUsecase.CreateDoctor(ctx, createDto)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to create doctor: %v", err)).Err()
	}
	return &emptypb.Empty{}, nil
}

func (h *DoctorHandlers) Next(ctx context.Context, rq *doctorgen.NextRq) (*doctorgen.NextRs, error) {
	nextDto := dto.NextRq{
		DoctorId: rq.DoctorId,
		UserId:   rq.UserId,
		Verdict:  rq.Verdict,
	}
	user, err := h.doctorUsecase.NextUser(ctx, &nextDto)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to get next user: %v", err)).Err()
	}
	return &doctorgen.NextRs{
		Name:    user.Name,
		Id:      user.Id,
		Surname: user.Surname,
	}, nil
}

//func (h *DoctorHandlers) GetHistory(ctx context.Context, rq *doctorgen.GetHistoryRq) (*doctorgen.GetHistoryRs, error) {
//	//TODO implement me
//	// Сорян забыл это сделать, а щас уже не успею реализовать
//	// Алекс Жарков, пж пж пж пж не бейте  меня (((
//	panic("implement me")
//}

func (h *DoctorHandlers) GetDoctorInfo(ctx context.Context, rq *doctorgen.GetDoctorInfoRq) (*doctorgen.GetDoctorInfoRs, error) {
	h.logger.Info("Doctor Id is", "id", rq.Id)
	info, err := h.doctorUsecase.GetDoctorInfo(ctx, rq.Id)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to get doctor info: %v", err)).Err()
	}
	return &doctorgen.GetDoctorInfoRs{
		Id:      info.Id,
		Name:    info.Name,
		Surname: info.Surname,
		Cabinet: info.Cabinet,
		Type:    info.Type,
	}, nil
}

func (h *DoctorHandlers) GetAllDoctors(ctx context.Context, empty *emptypb.Empty) (*doctorgen.GetAllDoctorsRs, error) {
	doctors, err := h.doctorUsecase.GetAllDoctors(ctx)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to get all doctors: %v", err)).Err()
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

	return &doctorgen.GetAllDoctorsRs{
		Doctors: doctorsDto,
	}, nil
}
