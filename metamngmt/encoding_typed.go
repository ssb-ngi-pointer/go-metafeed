// SPDX-FileCopyrightText: 2021 The go-metafeed Authors
//
// SPDX-License-Identifier: MIT

package metamngmt

import (
	"github.com/ssbc/go-metafeed/internal/bencodeext"
	"github.com/zeebo/bencode"
)

// UnmarshalBencode unpacks bencode extended data into an Typed message.
func (t *Typed) UnmarshalBencode(input []byte) error {
	var wt wrappedTyped
	err := bencode.DecodeBytes(input, &wt)
	if err != nil {
		return err
	}
	t.Type = string(wt.Type)
	return nil
}

type wrappedTyped struct {
	Type bencodeext.String `bencode:"type"`
}
