// SPDX-License-Identifier: MIT

package metamngmt

import (
	"fmt"

	"github.com/ssb-ngi-pointer/go-metafeed/internal/bencodeext"
	"github.com/zeebo/bencode"
	"go.mindeco.de/ssb-refs/tfk"
)

type wrappedTombstone struct {
	Type    bencodeext.String `bencode:"type"`
	SubFeed []byte            `bencode:"subfeed"`

	Tangles map[string]bencodeext.TanglePoint `bencode:"tangles"`
}

// MarshalBencode packs an Tombstone message into bencode extended data.
func (t Tombstone) MarshalBencode() ([]byte, error) {
	var wt wrappedTombstone
	wt.Type = bencodeext.String(t.Type)
	wt.Tangles = tanglesToBencoded(t.Tangles)

	subFeedTFK, err := tfk.FeedFromRef(t.SubFeed)
	if err != nil {
		return nil, fmt.Errorf("metafeed/tombstone: failed to turn subfeed into tfk: %w", err)
	}
	sfBytes, err := subFeedTFK.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("metafeed/tombstone: failed to encode tfk subfeed: %w", err)
	}

	wt.SubFeed = sfBytes

	return bencode.EncodeBytes(wt)
}

// UnmarshalBencode unpacks bencode extended data into an Tombstone message.
func (t *Tombstone) UnmarshalBencode(input []byte) error {
	var wt wrappedTombstone
	err := bencode.DecodeBytes(input, &wt)
	if err != nil {
		return fmt.Errorf("metamgngmt/tombstone: failed to unwrap bencode value: %w", err)
	}

	t.Type = string(wt.Type)
	if t.Type != "metafeed/tombstone" {
		return fmt.Errorf("metafeed/tombstone: invalid message type: %s", t.Type)
	}

	var subFeed tfk.Feed
	err = subFeed.UnmarshalBinary(wt.SubFeed)
	if err != nil {
		return fmt.Errorf("metafeed/tombstone: failed to decode tfk subfeed: %w", err)
	}

	t.SubFeed, err = subFeed.Feed()
	if err != nil {
		return fmt.Errorf("metafeed/tombstone: failed to turn subfeed tfk into feed: %w", err)
	}

	t.Tangles = bencodedToRefTangles(wt.Tangles)

	return nil
}
