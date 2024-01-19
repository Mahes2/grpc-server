package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime/debug"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/grpc-server/pb"
	"github.com/grpc-server/services/employee"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	rpcLogger log.Logger
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	rpcLogger = log.With(logger, "service", "gRPC/server", "component", "grpc-example")

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		// Order matters e.g. tracing interceptor have to create span first for the later exemplars to work.
		logging.UnaryServerInterceptor(
			interceptorLogger(rpcLogger),
			logging.WithFieldsFromContext(generateLogFields),
			logging.WithLogOnEvents(
				logging.StartCall,
				logging.PayloadReceived,
				logging.PayloadSent,
				logging.FinishCall,
			),
		),
		selector.UnaryServerInterceptor(
			auth.UnaryServerInterceptor(authenticator),
			selector.MatchFunc(authMatcher),
		),
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(recoveryHandler)),
	))

	pb.RegisterEmployeeServer(s, &employee.Server{})
	reflection.Register(s)

	port := "1531"
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		level.Error(logger).Log("failed to listen: %v", err)
	}

	level.Info(logger).Log("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		level.Error(logger).Log("failed to serve: %v", err)
	}
}

func authMatcher(ctx context.Context, callMeta interceptors.CallMeta) bool {
	return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
}

func authenticator(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	// TODO: This is example only, perform proper Oauth/OIDC verification!
	if token != "yolo" {
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	}
	// NOTE: You can also pass the token in the context for further interceptors or gRPC service code.
	return ctx, nil
}

func interceptorLogger(l log.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		largs := append([]any{"msg", msg}, fields...)
		switch lvl {
		case logging.LevelDebug:
			_ = level.Debug(l).Log(largs...)
		case logging.LevelInfo:
			_ = level.Info(l).Log(largs...)
		case logging.LevelWarn:
			_ = level.Warn(l).Log(largs...)
		case logging.LevelError:
			_ = level.Error(l).Log(largs...)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func generateLogFields(ctx context.Context) logging.Fields {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return logging.Fields{"traceID", span.TraceID().String()}
	}
	return nil
}

func recoveryHandler(p any) (err error) {
	level.Error(rpcLogger).Log("msg", "recovered from panic", "panic", p, "stack", debug.Stack())
	return status.Errorf(codes.Internal, "%s", p)
}
