/*
 * umoci: Umoci Modifies Open Containers' Images
 * Copyright (C) 2021 SUSE LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// +build gofuzz

package hardening

import (
	"bytes"
	_ "crypto/sha256" // Import is necessary for go-digest
	"github.com/opencontainers/go-digest"
	"io/ioutil"
)

// Fuzz implements the go-fuzz harness
func Fuzz(data []byte) int {
	buffer := bytes.NewBuffer(data)
	size := len(data)
	if !digest.SHA256.Available() {
		return -1
	}
	expectedDigest := digest.SHA256.FromBytes(buffer.Bytes())
	verifiedReader := &VerifiedReadCloser{
		Reader:         ioutil.NopCloser(buffer),
		ExpectedDigest: expectedDigest,
		ExpectedSize:   int64(size),
	}
	_, err := verifiedReader.Read(data)
	if err != nil {
		return 0
	}
	return 1
}
