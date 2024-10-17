package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/emorydu/dbaudit/internal/common"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
)

func RunGRPCServer(registerServer func(server *grpc.Server)) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	RunGRPCServerOnAddr(addr, registerServer)
}

func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	tlsOptions, err := tlsOpt()
	if err != nil {
		return
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			// TODO
			grpc_logrus.UnaryServerInterceptor(nil),
		),
		grpc.ChainStreamInterceptor(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			// TODO
			grpc_logrus.StreamServerInterceptor(nil),
		),
		tlsOptions,
	)

	registerServer(grpcServer)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		// TODO: logs
		return
	}

	// TODO: logs startup
	grpcServer.Serve(l)

}

func tlsOpt() (grpc.ServerOption, error) {
	cert, err := tls.X509KeyPair([]byte(common.ServerCert), []byte(common.ServerKey))
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(common.CaKey)); !ok {
		return nil, errors.New("append certs from pem error")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	return grpc.Creds(creds), nil
}
