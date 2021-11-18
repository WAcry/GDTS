// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package dispatch

import (
	"context"

	"GDTS/utils/mongodb/mongo-go-driver/core/command"
	"GDTS/utils/mongodb/mongo-go-driver/core/description"
	"GDTS/utils/mongodb/mongo-go-driver/core/topology"
)

// ListIndexes handles the full cycle dispatch and execution of a listIndexes command against the provided
// topology.
func ListIndexes(
	ctx context.Context,
	cmd command.ListIndexes,
	topo *topology.Topology,
	selector description.ServerSelector,
) (command.Cursor, error) {

	ss, err := topo.SelectServer(ctx, selector)
	if err != nil {
		return nil, err
	}

	conn, err := ss.Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return cmd.RoundTrip(ctx, ss.Description(), ss, conn)
}
