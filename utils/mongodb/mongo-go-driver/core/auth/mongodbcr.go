// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"crypto/md5"
	"fmt"

	"io"

	"GDTS/utils/mongodb/mongo-go-driver/bson"
	"GDTS/utils/mongodb/mongo-go-driver/core/command"
	"GDTS/utils/mongodb/mongo-go-driver/core/description"
	"GDTS/utils/mongodb/mongo-go-driver/core/wiremessage"
)

// MONGODBCR is the mechanism name for MONGODB-CR.
//
// The MONGODB-CR authentication mechanism is deprecated in MongoDB 4.0.
const MONGODBCR = "MONGODB-CR"

func newMongoDBCRAuthenticator(cred *Cred) (Authenticator, error) {
	return &MongoDBCRAuthenticator{
		DB:       cred.Source,
		Username: cred.Username,
		Password: cred.Password,
	}, nil
}

// MongoDBCRAuthenticator uses the MONGODB-CR algorithm to authenticate a connection.
//
// The MONGODB-CR authentication mechanism is deprecated in MongoDB 4.0.
type MongoDBCRAuthenticator struct {
	DB       string
	Username string
	Password string
}

// Auth authenticates the connection.
//
// The MONGODB-CR authentication mechanism is deprecated in MongoDB 4.0.
func (a *MongoDBCRAuthenticator) Auth(ctx context.Context, desc description.Server, rw wiremessage.ReadWriter) error {

	// Arbiters cannot be authenticated
	if desc.Kind == description.RSArbiter {
		return nil
	}

	db := a.DB
	if db == "" {
		db = defaultAuthDB
	}

	cmd := command.Read{DB: db, Command: bson.NewDocument(bson.EC.Int32("getnonce", 1))}
	ssdesc := description.SelectedServer{Server: desc}
	rdr, err := cmd.RoundTrip(ctx, ssdesc, rw)
	if err != nil {
		return newError(err, MONGODBCR)
	}

	var getNonceResult struct {
		Nonce string `bson:"nonce"`
	}

	err = bson.Unmarshal(rdr, &getNonceResult)
	if err != nil {
		return newAuthError("unmarshal error", err)
	}

	cmd = command.Read{
		DB: db,
		Command: bson.NewDocument(
			bson.EC.Int32("authenticate", 1),
			bson.EC.String("user", a.Username),
			bson.EC.String("nonce", getNonceResult.Nonce),
			bson.EC.String("key", a.createKey(getNonceResult.Nonce)),
		),
	}
	_, err = cmd.RoundTrip(ctx, ssdesc, rw)
	if err != nil {
		return newError(err, MONGODBCR)
	}

	return nil
}

func (a *MongoDBCRAuthenticator) createKey(nonce string) string {
	h := md5.New()

	_, _ = io.WriteString(h, nonce)
	_, _ = io.WriteString(h, a.Username)
	_, _ = io.WriteString(h, mongoPasswordDigest(a.Username, a.Password))
	return fmt.Sprintf("%x", h.Sum(nil))
}
