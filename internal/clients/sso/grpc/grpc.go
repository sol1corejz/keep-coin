package sso

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/sol1corejz/keep-coin/internal/domain/models"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"

	"github.com/sol1corejz/keep-coin/internal/lib/cert"
	ssov1 "github.com/sol1corejz/sso-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"os"
	"time"
)

type Client struct {
	api ssov1.AuthClient
}

func New(
	addr string,
	timeout time.Duration,
	log *slog.Logger,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	_, err := generateTlsConfig()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создание клиента с TLS
	cc, err := grpc.NewClient(
		addr,
		//grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: ssov1.NewAuthClient(cc),
	}, nil
}

func generateTlsConfig() (*tls.Config, error) {
	if !cert.CertExists() {
		slog.Info("Generating new TLS certificate")
		certPEM, keyPEM := cert.GenerateCert()
		if err := cert.SaveCert(certPEM, keyPEM); err != nil {
			panic("failed to save TLS certificate")
		}
	}

	slog.Info("loading TLS certificate")
	// Загрузка сертификата
	pemServerCA, err := os.ReadFile(cert.CertificateFilePath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, err
	}

	// Настройка TLS
	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}

	return tlsConfig, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (c *Client) Register(ctx context.Context, data models.AuthRequest) (*ssov1.RegisterResponse, error) {
	const op = "grpc.Register"
	resp, err := c.api.Register(ctx, &ssov1.RegisterRequest{
		Email:    data.Email,
		Password: data.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}

func (c *Client) Login(ctx context.Context, data models.AuthRequest) (*ssov1.LoginResponse, error) {
	const op = "grpc.Login"
	resp, err := c.api.Login(ctx, &ssov1.LoginRequest{
		Email:    data.Email,
		Password: data.Password,
		AppName:  "coin-keeper",
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}
