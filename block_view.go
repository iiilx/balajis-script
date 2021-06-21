package main

import (
	"github.com/btcsuite/btcd/btcec"
)

// block_view.go is the main work-horse for validating transactions in blocks.
// It generally works by creating an "in-memory view" of the current tip and
// then applying a transaction's operations to the view to see if those operations
// are allowed and consistent with the blockchain's current state. Generally,
// every transaction we define has a corresponding connect() and disconnect()
// function defined here that specifies what operations that transaction applies
// to the view and ultimately to the database. If you want to know how any
// particular transaction impacts the database, you've found the right file. A
// good place to start in this file is ConnectTransaction and DisconnectTransaction.
// ConnectBlock is also good.

const HashSizeBytes = 32

type PKID [33]byte

type UtxoType uint8
type BlockHash [HashSizeBytes]byte

const (
	// UTXOs can come from different sources. We document all of those sources
	// in the UTXOEntry using these types.
	UtxoTypeOutput      UtxoType = 0
	UtxoTypeBlockReward UtxoType = 1
	UtxoTypeBitcoinBurn UtxoType = 2
	// TODO(DELETEME): Remove the StakeReward txn type
	UtxoTypeStakeReward              UtxoType = 3
	UtxoTypeCreatorCoinSale          UtxoType = 4
	UtxoTypeCreatorCoinFounderReward UtxoType = 5

	// NEXT_TAG = 6
)


type PkMapKey [btcec.PubKeyBytesLenCompressed]byte


type MessageKey struct {
	PublicKey   PkMapKey
	BlockHeight uint32
	TstampNanos uint64
}


// MessageEntry stores the essential content of a message transaction.
type MessageEntry struct {
	SenderPublicKey    []byte
	RecipientPublicKey []byte
	EncryptedText      []byte
	// TODO: Right now a sender can fake the timestamp and make it appear to
	// the recipient that she sent messages much earlier than she actually did.
	// This isn't a big deal because there is generally not much to gain from
	// faking a timestamp, and it's still impossible for a user to impersonate
	// another user, which is the important thing. Moreover, it is easy to fix
	// the timestamp spoofing issue: You just need to make it so that the nodes
	// index messages based on block height in addition to on the tstamp. The
	// reason I didn't do it yet is because it adds some complexity around
	// detecting duplicates, particularly if a transaction is allowed to have
	// zero inputs/outputs, which is advantageous for various reasons.
	TstampNanos uint64

	isDeleted bool
}

// Entry for a public key forbidden from signing blocks.
type ForbiddenPubKeyEntry struct {
	PubKey []byte

	// Whether or not this entry is deleted in the view.
	isDeleted bool
}

type LikeKey struct {
	LikerPubKey   PkMapKey
	LikedPostHash BlockHash
}

// LikeEntry stores the content of a like transaction.
type LikeEntry struct {
	LikerPubKey   []byte
	LikedPostHash *BlockHash

	// Whether or not this entry is deleted in the view.
	isDeleted bool
}

type FollowKey struct {
	FollowerPKID PKID
	FollowedPKID PKID
}

// FollowEntry stores the content of a follow transaction.
type FollowEntry struct {
	// Note: It's a little redundant to have these in the entry because they're
	// already used as the key in the DB but it doesn't hurt for now.
	FollowerPKID *PKID
	FollowedPKID *PKID

	// Whether or not this entry is deleted in the view.
	isDeleted bool
}

type DiamondKey struct {
	SenderPKID      PKID
	ReceiverPKID    PKID
	DiamondPostHash BlockHash
}


// DiamondEntry stores the number of diamonds given by a sender to a post.
type DiamondEntry struct {
	SenderPKID      *PKID
	ReceiverPKID    *PKID
	DiamondPostHash *BlockHash
	DiamondLevel    int64

	// Whether or not this entry is deleted in the view.
	isDeleted bool
}

type RecloutKey struct {
	ReclouterPubKey PkMapKey
	// Post Hash of post that was reclouted
	RecloutedPostHash BlockHash
}

// RecloutEntry stores the content of a Reclout transaction.
type RecloutEntry struct {
	ReclouterPubKey []byte

	// BlockHash of the reclout
	RecloutPostHash *BlockHash

	// Post Hash of post that was reclouted
	RecloutedPostHash *BlockHash

	// Whether or not this entry is deleted in the view.
	isDeleted bool
}

type GlobalParamsEntry struct {
	// The new exchange rate to set.
	USDCentsPerBitcoin uint64

	// The new create profile fee
	CreateProfileFeeNanos uint64

	// The new minimum fee the network will accept
	MinimumNetworkFeeNanosPerKB uint64
}


// This struct holds info on a readers interactions (e.g. likes) with a post.
// It is added to a post entry response in the frontend server api.
type PostEntryReaderState struct {
	// This is true if the reader has liked the associated post.
	LikedByReader bool

	// The number of diamonds that the reader has given this post.
	DiamondLevelBestowed int64

	// This is true if the reader has reclouted the associated post.
	RecloutedByReader bool

	// This is the post hash hex of the reclout
	RecloutPostHashHex string
}


type SingleStake struct {
	// Just save the data from the initial stake for posterity.
	InitialStakeNanos               uint64
	BlockHeight                     uint64
	InitialStakeMultipleBasisPoints uint64
	// The amount distributed to previous users can be computed by
	// adding the creator percentage and the burn fee and then
	// subtracting that total percentage off of the InitialStakeNanos.
	// Example:
	// - InitialStakeNanos = 100
	// - CreatorPercentage = 15%
	// - BurnFeePercentage = 10%
	// - Amount to pay to previous users = 100 - 15 - 10 = 75
	InitialCreatorPercentageBasisPoints uint64

	// These fields are what we actually use to pay out the user who staked.
	//
	// The initial RemainingAmountOwedNanos is computed by simply multiplying
	// the InitialStakeNanos by the InitialStakeMultipleBasisPoints.
	RemainingStakeOwedNanos uint64
	PublicKey               []byte
}

type StakeEntry struct {
	StakeList []*SingleStake

	// Computed for profiles to cache how much has been staked to
	// their posts in total. When a post is staked to, this value
	// gets incremented on the profile. It gets reverted on the
	// profile when the post stake is reverted.
	TotalPostStake uint64
}

type StakeEntryStats struct {
	TotalStakeNanos           uint64
	TotalStakeOwedNanos       uint64
	TotalCreatorEarningsNanos uint64
	TotalFeesBurnedNanos      uint64
	TotalPostStakeNanos       uint64
}

type StakeIDType uint8

const (
	StakeIDTypePost    StakeIDType = 0
	StakeIDTypeProfile StakeIDType = 1
)

type PostEntry struct {
	// The hash of this post entry. Used as the ID for the entry.
	PostHash *BlockHash

	// The public key of the user who made the post.
	PosterPublicKey []byte

	// The parent post. This is used for comments.
	ParentStakeID []byte

	// The body of this post.
	Body []byte

	// The PostHash of the post this post reclouts
	RecloutedPostHash *BlockHash

	// Indicator if this PostEntry is a quoted reclout or not
	IsQuotedReclout bool

	// The amount the creator of the post gets when someone stakes
	// to the post.
	CreatorBasisPoints uint64

	// The multiple of the payout when a user stakes to a post.
	// 2x multiple = 200% = 20,000bps
	StakeMultipleBasisPoints uint64

	// The block height when the post was confirmed.
	ConfirmationBlockHeight uint32

	// A timestamp used for ordering messages when displaying them to
	// users. The timestamp must be unique. Note that we use a nanosecond
	// timestamp because it makes it easier to deal with the uniqueness
	// constraint technically (e.g. If one second spacing is required
	// as would be the case with a standard Unix timestamp then any code
	// that generates these transactions will need to potentially wait
	// or else risk a timestamp collision. This complexity is avoided
	// by just using a nanosecond timestamp). Note that the timestamp is
	// an unsigned int as opposed to a signed int, which means times
	// before the zero time are not represented which doesn't matter
	// for our purposes. Restricting the timestamp in this way makes
	// lexicographic sorting based on bytes easier in our database which
	// is one of the reasons we do it.
	TimestampNanos uint64

	// Users can "delete" posts, but right now we just implement this as
	// setting a flag on the post to hide it rather than actually deleting
	// it. This simplifies the implementation and makes it easier to "undelete"
	// posts in certain situations.
	IsHidden bool

	// Every post has a StakeEntry that keeps track of all the stakes that
	// have been applied to this post.
	StakeEntry *StakeEntry

	// Counter of users that have liked this post.
	LikeCount uint64

	// Counter of users that have reclouted this post.
	RecloutCount uint64

	// Counter of quote reclouts for this post.
	QuoteRecloutCount uint64

	// Counter of diamonds that the post has received.
	DiamondCount uint64

	// The private fields below aren't serialized or hashed. They are only kept
	// around for in-memory bookkeeping purposes.

	// Used to sort posts by their stake. Generally not set.
	stakeStats *StakeEntryStats

	// Whether or not this entry is deleted in the view.
	isDeleted bool

	// How many comments this post has
	CommentCount uint64

	// Indicator if a post is pinned or not.
	IsPinned bool

	// ExtraData map to hold arbitrary attributes of a post. Holds non-consensus related information about a post.
	PostExtraData map[string][]byte
}

type BalanceEntryMapKey struct {
	HODLerPKID  PKID
	CreatorPKID PKID
}


// This struct is mainly used to track a user's balance of a particular
// creator coin. In the database, we store it as the value in a mapping
// that looks as follows:
// <HodlerPKID, CreatorPKID> -> HODLerEntry
type BalanceEntry struct {
	// The PKID of the HODLer. This should never change after it's set initially.
	HODLerPKID *PKID
	// The PKID of the creator. This should never change after it's set initially.
	CreatorPKID *PKID

	// How much this HODLer owns of a particular creator coin.
	BalanceNanos uint64

	// Has the hodler purchased any amount of this user's coin
	HasPurchased bool

	// Whether or not this entry is deleted in the view.
	isDeleted bool
}

// This struct contains all the information required to support coin
// buy/sell transactions on profiles.
type CoinEntry struct {
	// The amount the owner of this profile receives when there is a
	// "net new" purchase of their coin.
	CreatorBasisPoints uint64

	// The amount of BitClout backing the coin. Whenever a user buys a coin
	// from the protocol this amount increases, and whenever a user sells a
	// coin to the protocol this decreases.
	BitCloutLockedNanos uint64

	// The number of public keys who have holdings in this creator coin.
	// Due to floating point truncation, it can be difficult to simultaneously
	// reset CoinsInCirculationNanos and BitCloutLockedNanos to zero after
	// everyone has sold all their creator coins. Initially NumberOfHolders
	// is set to zero. Once it returns to zero after a series of buys & sells
	// we reset the BitCloutLockedNanos and CoinsInCirculationNanos to prevent
	// abnormal bancor curve behavior.
	NumberOfHolders uint64

	// The number of coins currently in circulation. Whenever a user buys a
	// coin from the protocol this increases, and whenever a user sells a
	// coin to the protocol this decreases.
	CoinsInCirculationNanos uint64

	// This field keeps track of the highest number of coins that has ever
	// been in circulation. It is used to determine when a creator should
	// receive a "founder reward." In particular, whenever the number of
	// coins being minted would push the number of coins in circulation
	// beyond the watermark, we allocate a percentage of the coins being
	// minted to the creator as a "founder reward."
	CoinWatermarkNanos uint64
}

type PKIDEntry struct {
	PKID *PKID
	// We add the public key only so we can reuse this struct to store the reverse
	// mapping of pkid -> public key.
	PublicKey []byte

	isDeleted bool
}

type ProfileEntry struct {
	// PublicKey is the key used by the user to sign for things and generally
	// verify her identity.
	PublicKey []byte

	// Username is a unique human-readable identifier associated with a profile.
	Username []byte

	// Some text describing the profile.
	Description []byte

	// The profile pic string encoded as a link e.g.
	// data:image/png;base64,<data in base64>
	ProfilePic []byte

	// Users can "delete" profiles, but right now we just implement this as
	// setting a flag on the post to hide it rather than actually deleting
	// it. This simplifies the implementation and makes it easier to "undelete"
	// profiles in certain situations.
	IsHidden bool

	// CoinEntry tracks the information required to buy/sell coins on a user's
	// profile. We "embed" it here for convenience so we can access the fields
	// directly on the ProfileEntry object. Embedding also makes it so that we
	// don't need to initialize it explicitly.
	CoinEntry

	// Whether or not this entry should be deleted when the view is flushed
	// to the db. This is initially set to false, but can become true if for
	// example we update a user entry and need to delete the data associated
	// with the old entry.
	isDeleted bool

	// TODO(DELETEME): This field is deprecated. It was relevant back when
	// we wanted to allow people to stake to profiles, which isn't something
	// we want to support going forward.
	//
	// The multiple of the payout when a user stakes to this profile. If
	// unset, a sane default is set when the first person stakes to this
	// profile.
	// 2x multiple = 200% = 20,000bps
	StakeMultipleBasisPoints uint64

	// TODO(DELETEME): This field is deprecated. It was relevant back when
	// we wanted to allow people to stake to profiles, which isn't something
	// we want to support going forward.
	//
	// Every provile has a StakeEntry that keeps track of all the stakes that
	// have been applied to it.
	StakeEntry *StakeEntry

	// The private fields below aren't serialized or hashed. They are only kept
	// around for in-memory bookkeeping purposes.

	// TODO(DELETEME): This field is deprecated. It was relevant back when
	// we wanted to allow people to stake to profiles, which isn't something
	// we want to support going forward.
	//
	// Used to sort profiles by their stake. Generally not set.
	stakeStats *StakeEntryStats
}

type OperationType uint