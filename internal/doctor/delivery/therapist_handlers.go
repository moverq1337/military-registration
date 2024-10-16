package delivery

import (
	"context"
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
	"log/slog"
)

type TherapistHandlers struct {
	logger        *slog.Logger
	doctorUsecase usecase.DoctorUsecase
	doctorgen.UnimplementedTherapistServiceServer
}

func NewTherapistHandlers(logger *slog.Logger, doctorUsecase usecase.DoctorUsecase) *TherapistHandlers {
	return &TherapistHandlers{logger: logger, doctorUsecase: doctorUsecase}
}

func RegisterTherapistHandlers(
	grpcSrv *grpc.Server,
	grpcGtw *runtime.ServeMux,
	grpcConfig *grpcLib.Config,
	logger *slog.Logger,
	doctorUsecase usecase.DoctorUsecase,
) {

	ctx := context.Background()

	impl := NewTherapistHandlers(logger, doctorUsecase)
	doctorgen.RegisterTherapistServiceServer(grpcSrv, impl)
	err := doctorgen.RegisterTherapistServiceHandlerFromEndpoint(ctx, grpcGtw, grpcConfig.Address(), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		logger.Error("failed to register example service handler", "error", err)
		return
	}
	ctx = context.WithValue(ctx, "parameter", "im fucking parametr")
	slog.InfoContext(ctx, "hello world")
}

func (h *TherapistHandlers) Next(ctx context.Context, rq *doctorgen.NextRq) (*doctorgen.NextRs, error) {
	user, err := h.doctorUsecase.NextUserTherapist(ctx, &dto.TherapistNextRq{
		UserId:  rq.UserId,
		Verdict: rq.Verdict,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &doctorgen.NextRs{
		Id:      user.Id,
		Name:    user.Name,
		Surname: user.Surname,
	}, nil
}

// FinalVerdict Deprecated
// Just use Next
func (h *TherapistHandlers) FinalVerdict(ctx context.Context, rq *doctorgen.FinalVerdictRq) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
