// SPDX-FileCopyrightText: 2021 The go-metafeed Authors
//
// SPDX-License-Identifier: MIT

package metafeed_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/bencode"
	refs "go.mindeco.de/ssb-refs"
	"go.mindeco.de/ssb-refs/tfk"

	metafeed "github.com/ssb-ngi-pointer/go-metafeed"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func createTestEntry(author refs.FeedRef) metafeed.Payload {
	// now construct a test entry
	return metafeed.Payload{
		Author:    author,
		Sequence:  1,
		Previous:  nil,
		Timestamp: time.Unix(10, 0), // 10 seconds after midnight
		Content:   bencode.RawMessage("12:hello, world"),
	}
}

func ExamplePayload() {
	pubKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	exAuthor, err := refs.NewFeedRefFromBytes(pubKey, refs.RefAlgoFeedBendyButt)
	check(err)

	exampleFile, err := os.OpenFile("example-feed-entry.bendybutt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	check(err)
	defer exampleFile.Close()

	var (
		buf    bytes.Buffer
		dumper = hex.Dumper(os.Stdout)
		w      = io.MultiWriter(&buf, dumper, exampleFile)

		entry = createTestEntry(exAuthor)
	)

	encoded, err := entry.MarshalBencode()
	check(err)

	_, err = w.Write(encoded)
	check(err)

	var decodedEntry []interface{}
	err = bencode.NewDecoder(&buf).Decode(&decodedEntry)
	check(err)

	authorAsString := decodedEntry[0].(string)
	var decodedAuthor tfk.Feed
	err = decodedAuthor.UnmarshalBinary([]byte(authorAsString))
	check(err)

	decodedAuthorRef, err := decodedAuthor.Feed()
	check(err)

	if !decodedAuthorRef.Equal(exAuthor) {
		fmt.Println("wrong author")
	}

	// Output:
	// 00000000  6c 33 34 3a 00 03 01 02  03 04 05 06 07 08 09 0a  |l34:............|
	// 00000010  0b 0c 0d 0e 0f 10 11 12  13 14 15 16 17 18 19 1a  |................|
	// 00000020  1b 1c 1d 1e 1f 20 69 31  65 32 3a 06 02 69 31 30  |..... i1e2:..i10|
	// 00000030  65 31 32 3a 68 65 6c 6c  6f 2c 20 77 6f 72 6c 64  |e12:hello, world|
	// 00000040  65
}

func TestDecodeEntry(t *testing.T) {
	r := require.New(t)

	pubKey := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	exAuthor, err := refs.NewFeedRefFromBytes(pubKey, refs.RefAlgoFeedBendyButt)
	r.NoError(err)

	var (
		buf bytes.Buffer

		entry = createTestEntry(exAuthor)
	)

	encoded, err := entry.MarshalBencode()
	check(err)

	_, err = buf.Write(encoded)
	check(err)

	var ee metafeed.Payload

	err = bencode.NewDecoder(&buf).Decode(&ee)
	r.NoError(err)

	got := fmt.Sprintln()
	got += fmt.Sprintln("Author:", ee.Author.Sigil())
	got += fmt.Sprintln("Seq:", ee.Sequence)
	got += fmt.Sprintln("Previous:", ee.Previous)
	got += fmt.Sprintln("Timestamp:", ee.Timestamp.UTC().String())
	got += fmt.Sprintf("Content: %q", string(ee.Content))

	want := `
Author: @AQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eHyA=.bendybutt-v1
Seq: 1
Previous: <nil>
Timestamp: 1970-01-01 00:00:10 +0000 UTC
Content: "12:hello, world"`
	r.Equal(want, got)
}

// basic bencode test, showing that 1 and 0 integers are used to represent true and false.
// See internal/bencodeext for how we deal with these ambiguities.
func TestBoolValues(t *testing.T) {
	r := require.New(t)

	trueAsBytes, err := bencode.EncodeBytes(true)
	r.NoError(err)

	oneAsBytes, err := bencode.EncodeBytes(1)
	r.NoError(err)

	zeroAsBytes, err := bencode.EncodeBytes(0)
	r.NoError(err)

	r.Equal(trueAsBytes, oneAsBytes)

	var truethy bool
	err = bencode.DecodeBytes(oneAsBytes, &truethy)
	r.NoError(err)

	r.True(truethy)

	err = bencode.DecodeBytes(zeroAsBytes, &truethy)
	r.NoError(err)

	r.False(truethy)
}
