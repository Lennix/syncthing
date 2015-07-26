// Copyright (C) 2014 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package protocol

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/rand"
	"crypto/sha256"
)

func Encrypt(buf []byte, label []byte, cert tls.Certificate) (out []byte, err error) {
	var ret []byte

	l.Debugln("Before encryption: ", buf)

	// Certificate stuff
	pub, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		l.Debugln("error:", err)
		return nil, err
	}

	pubkey := pub.PublicKey.(*rsa.PublicKey)


	// now to encrypting
	// each encrypted chunk may only be ((pubkey.N.BitLen() + 7) / 8) - 11 byte big, so we may have to cut here


	k := ((pubkey.N.BitLen() + 7) / 8) - 11

	var offset int
	
	for i := 0; i < len(buf); i += k {
		out, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubkey, buf[i:k], label)
		if err != nil {
			l.Debugln("error:", err)
			return nil, err
		}
		ret = append(ret, out...)

		offset += len(out)
	}

	l.Debugln("After encryption: ", ret)

	return ret, nil
}