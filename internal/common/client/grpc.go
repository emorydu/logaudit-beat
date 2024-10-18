// Copyright 2024 Emory.Du <orangeduxiaocheng@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/emorydu/dbaudit/internal/common"
	"github.com/emorydu/dbaudit/internal/common/genproto/auditbeat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewAuditBeatClient(grpcAddr string) (client auditbeat.AuditBeatServiceClient, close func() error, err error) {
	if grpcAddr == "" {
		return nil, func() error { return nil }, errors.New("empty auditbeat rpc addr")
	}

	opts, err := grpcDialOpts()
	if err != nil {
		return nil, func() error { return nil }, err
	}

	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*50)))

	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}

	return auditbeat.NewAuditBeatServiceClient(conn), conn.Close, nil
}

func grpcDialOpts() ([]grpc.DialOption, error) {
	tlsOptions, err := tlsOpt()
	return []grpc.DialOption{tlsOptions}, err
}

func tlsOpt() (grpc.DialOption, error) {
	cert, err := tls.X509KeyPair([]byte(common.ClientCert), []byte(common.ClientKey))
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(common.CaKey)); !ok {
		return nil, errors.New("append certs from pem error")
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		RootCAs:      certPool,

		InsecureSkipVerify: true,
	})

	return grpc.WithTransportCredentials(creds), nil
}
