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
	"errors"
	"math/big"
	"strconv"

	"github.com/danielakhterov/go-ethereum/common"
	"github.com/danielakhterov/go-ethereum/consensus"
	"github.com/danielakhterov/go-ethereum/core/state"
	"github.com/danielakhterov/go-ethereum/core/types"
	"github.com/danielakhterov/go-ethereum/params"
	"github.com/danielakhterov/go-ethereum/rpc"
	"github.com/hashgraph/hedera-sdk-go"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	client  *hedera.Client
	topicID hedera.ConsensusTopicID
	maxFee  hedera.Hbar
}

func NewAPI(config params.HashgraphConfig) (API, error) {
	id, err := hedera.AccountIDFromString(config.OpeartorID)
	if err != nil {
		return API{}, err
	}

	key, err := hedera.Ed25519PrivateKeyFromString(config.OperatorKey)
	if err != nil {
		return API{}, err
	}

	topicID, err := hedera.TopicIDFromString(config.TopicID)
	if err != nil {
		return API{}, err
	}

	client := hedera.ClientForTestnet().
		SetOperator(id, key)

	return API{
		client:  client,
		topicID: topicID,
		maxFee:  hedera.HbarFromTinybar(config.MaxFee),
	}, nil
}

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (api API) Author(header *types.Header) (common.Address, error) {
	return common.Address{}, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (api API) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (api API) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	return nil, nil
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (api API) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (api API) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return nil
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (api API) Prepare(chain consensus.ChainReader, header *types.Header) error {
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// but does not assemble the block.
//
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (api API) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {
	return
}

// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
// rewards) and assembles the final block.
//
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (api API) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	return nil, nil
}

// Seal generates a new sealing request for the given input block and pushes
// the result into the given channel.
//
// Note, the method returns immediately and will send the result async. More
// than one result may also be returned depending on the consensus algorithm.
func (api API) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	count := len(block.Transactions())
	for i, tx := range block.Transactions() {
		json, err := tx.MarshalJSON()
		if err != nil {
			return err
		}

		if len(json) > 4000 {
			return errors.New("unimplemetned")
		}

		memo := strconv.Itoa(i) + " " + strconv.Itoa(count) + " " + block.Hash().Hex()

		txID, err := hedera.NewConsensusMessageSubmitTransaction().
			SetMessage(json).
			SetTopicID(api.topicID).
			SetTransactionMemo(memo).
			SetMaxTransactionFee(api.maxFee).
			Execute(api.client)
		if err != nil {
			return err
		}

		receipt, err := txID.GetReceipt(api.client)
		if err != nil {
			return err
		}

		print(receipt.GetConsensusTopicID().String())
	}

	return nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (api API) SealHash(header *types.Header) common.Hash {
	return common.Hash{}
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have.
func (api API) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (api API) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{}
}

// Close terminates any background threads maintained by the consensus engine.
func (api API) Close() error {
	return api.client.Close()
}
