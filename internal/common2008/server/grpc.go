// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/emorydu/dbaudit/internal/common"
	"github.com/emorydu/dbaudit/internal/common/logs"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
)

func init() {
	logger := logrus.New()
	logs.SetFormatter(logger)
	logger.SetLevel(logrus.WarnLevel)

	grpc_logrus.ReplaceGrpcLogger(logrus.NewEntry(logger))
}

func RunGRPCServer(registerServer func(server *grpc.Server)) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	addr := fmt.Sprintf(":%s", port)
	RunGRPCServerOnAddr(addr, registerServer)
}

func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	tlsOptions, err := tlsOpt()
	if err != nil {
		return
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc.ChainStreamInterceptor(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
		),
		tlsOptions,
	)

	registerServer(grpcServer)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.WithField("grpcEndpoint", addr).Info("Starting: gRPC Listener")
	logrus.Fatal(grpcServer.Serve(l))

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
