package rng

import (
	"time"

	shamir "github.com/republicprotocol/shamir-go"
)

// ShareMap is a convenience type that associates addresses with shares.
type ShareMap map[Address]shamir.Share

// An InputMessage can be passed to the Rnger as an input. It will be processed
// by the Rnger and an error will be output if the message is an unexpected
// type. No types external to this package should implement this interface.
type InputMessage interface {

	// IsInputMessage is a marker used to restrict InputMessages to types that
	// have been explicitly marked. It is never called.
	IsInputMessage()
}

// An OutputMessage can be passed from the Rnger as an output. Depending on the
// type of output message, the user must route the message to the appropriate
// Rnger in the network. See the documentation specific to each message for
// information on how to handle it. No types external to this package should
// implement this interface.
type OutputMessage interface {

	// IsOutputMessage is a marker used to restrict OutputMessages to types that
	// have been explicitly marked. It is never called.
	IsOutputMessage()
}

// A Nominate message is used to nominate a computation leader. The Leader
// Address is the address of the nominated leader.
type Nominate struct {
	Leader Address
}

// NewNominate creates a new Nominate message.
func NewNominate(addr Address) Nominate {
	return Nominate{
		addr,
	}
}

// IsInputMessage implements the InputMessage interface.
func (message Nominate) IsInputMessage() {
}

// A GenerateRn message signals to the Rnger that is should begin a secure
// random number generation. The secure random number that will be generated is
// identified by a nonce. The nonce must be unique and must be agreed on by all
// Rngers in the network. After receiving this message, an Rnger will produce a
// LocalRnShare for all Rngers in the network. The user must route these
// LocalRnShare messages to their respective Rngers.
type GenerateRn struct {
	Nonce
}

// NewGenerateRn creates a new GenerateRn message.
func NewGenerateRn(nonce Nonce) GenerateRn {
	return GenerateRn{
		nonce,
	}
}

// IsInputMessage implements the InputMessage interface.
func (message GenerateRn) IsInputMessage() {
}

// A ProposeRn message is sent by a computation leader to the compute nodes. It
// is a proposal by the computation leader for the compute nodes to generate a
// local random number which will then be shared to eventually produce a global
// shared random number. This message is sent once the nominated leader has
// received a GenerateRn message.
type ProposeRn struct {
	Nonce

	To   Address
	From Address
}

// NewProposeRn creates a new ProposeRn message.
func NewProposeRn(nonce Nonce, to, from Address) ProposeRn {
	return ProposeRn{
		nonce,
		to,
		from,
	}
}

// IsOutputMessage implements the OutputMessage interface.
func (message ProposeRn) IsOutputMessage() {
}

// IsInputMessage implements the InputMessage interface.
func (message ProposeRn) IsInputMessage() {
}

// A LocalRnShares message is produced by an Rnger after receiving a GenerateRn
// message. A LocalRnShares message will be produced for each Rnger in the
// network and it is up to the user to route this message to the appropriate
// Rnger. A LocalRnShares message can also be passed to an Rnger as input,
// representing the LocalRnShares messages sent to it by other Rngers in the
// network.
type LocalRnShares struct {
	Nonce

	To     Address
	From   Address
	Shares ShareMap
}

// NewLocalRnShares creates a new LocalRnShares message.
func NewLocalRnShares(nonce Nonce, to, from Address, shares map[Address]shamir.Share) LocalRnShares {
	return LocalRnShares{
		nonce,

		to,
		from,
		shares,
	}
}

// IsInputMessage implements the InputMessage interface.
func (message LocalRnShares) IsInputMessage() {
}

// IsOutputMessage implements the OutputMessage interface.
func (message LocalRnShares) IsOutputMessage() {
}

// A ProposeGlobalRnShare message is sent by the computation leader and contains
// shares of the local random numbers generated by a selection of the compute
// nodes. This message is sent to each compute node, which can then sum up the
// contained shares to produce a single share of a global random number.
type ProposeGlobalRnShare struct {
	Nonce

	To     Address
	From   Address
	Shares ShareMap
}

// NewProposeGlobalRnShare creates a new ProposeGlobalRnShare message.
func NewProposeGlobalRnShare(nonce Nonce, to, from Address, shares map[Address]shamir.Share) ProposeGlobalRnShare {
	return ProposeGlobalRnShare{
		nonce,
		to,
		from,
		shares,
	}
}

// IsInputMessage implements the InputMessage interface.
func (message ProposeGlobalRnShare) IsInputMessage() {
}

// IsOutputMessage implements the OutputMessage interface.
func (message ProposeGlobalRnShare) IsOutputMessage() {
}

// A GlobalRnShare message is produced by an Rnger at the end of a successful
// secure random number generation. It is the Shamir's secret share of the
// secure random number that has been generated.
type GlobalRnShare struct {
	Nonce
	shamir.Share

	From Address
}

// NewGlobalRnShare creates a new GlobalRnShare message.
func NewGlobalRnShare(nonce Nonce, share shamir.Share, from Address) GlobalRnShare {
	return GlobalRnShare{
		nonce,
		share,
		from,
	}
}

// IsOutputMessage implements the OutputMessage interface.
func (message GlobalRnShare) IsOutputMessage() {
}

// A VoteGlobalRnShare message is produced by an Rnger after receiving a sufficient number of
// LocalRnShares messages, or after a secure random number generation has
// exceeded its deadline. A VoteGlobalRnShare message will be produced for each Rnger in the
// network and it is up to the user to route this message to the appropriate
// Rnger.
type VoteGlobalRnShare struct {
	Nonce

	To      Address
	From    Address
	Players []Address
}

// NewVoteGlobalRnShare creates a new VoteGlobalRnShare message.
func NewVoteGlobalRnShare(nonce Nonce, to, from Address, players []Address) VoteGlobalRnShare {
	return VoteGlobalRnShare{
		nonce,

		to,
		from,
		players,
	}
}

// IsInputMessage implements the InputMessage interface.
func (message VoteGlobalRnShare) IsInputMessage() {
}

// IsOutputMessage implements the OutputMessage interface.
func (message VoteGlobalRnShare) IsOutputMessage() {
}

// A CheckDeadline message signals to the Rnger that it should clean up all
// pending random number generations that have exceeded their deadline. It is up
// to the user to determine the frequency of this message. Higher frequencies
// will result in more accurate clean up times, but slower performance.
type CheckDeadline struct {
	time.Time
}

// NewCheckDeadline creates a new CheckDeadline message.
func NewCheckDeadline(time time.Time) CheckDeadline {
	return CheckDeadline{
		time,
	}
}

// IsInputMessage implements the InputMessage interface.
func (message CheckDeadline) IsInputMessage() {
}

// Err is a message that is sent by a player when they encounter an error during
// a random number generation.
type Err struct {
	Nonce
	error
}

// NewErr creates a new Err message.
func NewErr(nonce Nonce, err error) Err {
	return Err{
		nonce,
		err,
	}
}

// IsOutputMessage implements the OutputMessage interface.
func (message Err) IsOutputMessage() {
}
