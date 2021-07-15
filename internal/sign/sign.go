// SPDX-License-Identifier: MIT

// Package sign implements the domain seperated signature creation and verification used in bendybutt powered metafeeds.
package sign

import (
	"bytes"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/nacl/auth"
)

var (
	// this gets prepended to the sign()/verify() input and achives domain seperation
	inputPrefix = []byte("metafeeds")

	// these two bytes are TFK/BFE identifiers to clerify that the bytes are a signature
	outputPrefix = []byte{0x04, 0x00}
)

// Create creates the signature over the passed input bytes using the passed secret key.
// To achive domain seperation, the input is prefixed with the string "metafeeds".
// The resulting signature is prefixed with the two bytes 0x0400, which are bendybutt bencode extension to denote it being a signature.
// If hmacSec is not nil, the prefixed input is hashed using nacl's auth.Sum() before the signature is created.
func Create(input []byte, key ed25519.PrivateKey, hmacSec *[32]byte) []byte {
	toSign := append(inputPrefix, input...)
	if hmacSec != nil {
		mac := auth.Sum(toSign, hmacSec)
		toSign = mac[:]
	}

	sig := ed25519.Sign(key, toSign)
	return append(outputPrefix, sig...)
}

// Verify checks if the passed signature was indeed created over the passed data, using the same domain seperation as the Create() function in this package.
// It checks if the signature has the right prefix (0x0400) and then prepends the string "metafeeds" to the data.
// If hmacSec is not nil, the prefixed input is hashed using nacl's auth.Sum() before the signature is verified.
func Verify(data, signature []byte, pubKey ed25519.PublicKey, hmacSec *[32]byte) bool {
	if !bytes.HasPrefix(signature, outputPrefix) {
		return false
	}

	signedMessage := append(inputPrefix, data...)
	if hmacSec != nil {
		mac := auth.Sum(signedMessage, hmacSec)
		signedMessage = mac[:]
	}

	justTheSig := bytes.TrimPrefix(signature, outputPrefix)
	return ed25519.Verify(pubKey, signedMessage, justTheSig)
}
