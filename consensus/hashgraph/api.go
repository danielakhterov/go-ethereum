// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package hashgraph

import (
	"github.com/danielakhterov/go-ethereum/params"
	"github.com/hashgraph/hedera-sdk-go"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	client    *hedera.Client
	submitKey hedera.Ed25519PublicKey
}

func NewAPI(config params.HashgraphConfig) API {
	client := hedera.ClientForTestnet().
		SetOperator(config.OpeartorId, config.OperatorKey)

	return API{
		client:    client,
		submitKey: config.SubmitKey,
	}
}
