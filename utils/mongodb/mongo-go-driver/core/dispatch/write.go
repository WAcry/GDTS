// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package dispatch

import (
	"context"

	"GDTS/utils/mongodb/mongo-go-driver/bson"
	"GDTS/utils/mongodb/mongo-go-driver/core/command"
	"GDTS/utils/mongodb/mongo-go-driver/core/description"
	"GDTS/utils/mongodb/mongo-go-driver/core/topology"
	"GDTS/utils/mongodb/mongo-go-driver/core/writeconcern"
)

// Write handles the full cycle dispatch and execution of a write command against the provided
// topology.
func Write(
	ctx context.Context,
	cmd command.Write,
	topo *topology.Topology,
	selector description.ServerSelector,
) (bson.Reader, error) {
	ss, err := topo.SelectServer(ctx, selector)
	if err != nil {
		return nil, err
	}

	desc := ss.Description()
	conn, err := ss.Connection(ctx)
	if err != nil {
		return nil, err
	}

	if !writeconcern.AckWrite(cmd.WriteConcern) {
		go func() {
			defer func() { _ = recover() }()
			defer conn.Close()

			_, _ = cmd.RoundTrip(ctx, desc, conn)
		}()

		return nil, command.ErrUnacknowledgedWrite
	}
	defer conn.Close()

	return cmd.RoundTrip(ctx, desc, conn)
}
