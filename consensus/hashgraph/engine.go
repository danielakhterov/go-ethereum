package hashgraph

import (
	"math/big"

	"github.com/danielakhterov/go-ethereum/common"
	"github.com/danielakhterov/go-ethereum/consensus"
	"github.com/danielakhterov/go-ethereum/core/state"
	"github.com/danielakhterov/go-ethereum/core/types"
	"github.com/danielakhterov/go-ethereum/rpc"
	// "github.com/hashgraph/hedera-sdk-go"
)

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (api *API) Author(header *types.Header) (common.Address, error)

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (api *API) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (api *API) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (api *API) VerifyUncles(chain consensus.ChainReader, block *types.Block) error

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (api *API) VerifySeal(chain consensus.ChainReader, header *types.Header) error

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (api *API) Prepare(chain consensus.ChainReader, header *types.Header) error

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// but does not assemble the block.
//
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (api *API) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header)

// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
// rewards) and assembles the final block.
//
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (api *API) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

// Seal generates a new sealing request for the given input block and pushes
// the result into the given channel.
//
// Note, the method returns immediately and will send the result async. More
// than one result may also be returned depending on the consensus algorithm.
func (api *API) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

// SealHash returns the hash of a block prior to it being sealed.
func (api *API) SealHash(header *types.Header) common.Hash

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have.
func (api *API) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int

// APIs returns the RPC APIs this consensus engine provides.
func (api *API) APIs(chain consensus.ChainReader) []rpc.API

// Close terminates any background threads maintained by the consensus engine.
func (api *API) Close() error
