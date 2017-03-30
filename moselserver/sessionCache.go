/*
 * Copyright 2016 Robin Engel
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
package moselserver

import (
	"time"
	"math/rand"
	"crypto/md5"
	"encoding/hex"
)

// !!! CURRENTLY UNSAVE !!!
// It always produces the same session keys at the moment
type sessionCache struct {
	sessions map[string]session
}

type session struct {
	keyHash []byte
	validTo time.Time
}

func NewSessionCache() *sessionCache {
	cache := new(sessionCache)
	cache.sessions = make(map[string]session)
	return cache
}

func (cache *sessionCache) NewSession() (string, time.Time) {
	s := session{}

	b := make([]byte, 256)
	rand.Read(b)

	key := cache.hash(b)
	keyHash := cache.hash([]byte(key[:]))

	s.keyHash = []byte(keyHash[:])
	s.validTo = time.Now()

	keyHashString := hashToString(s.keyHash[:])
	cache.sessions[keyHashString] = s

	return hashToString(key[:]), time.Now()
}

func (cache *sessionCache) ValidateSession(key string) bool {
	keyBin, _ := hex.DecodeString(key)
	hash := cache.hash(keyBin)
	hashString := hashToString(hash)
	_, ok := cache.sessions[hashString]
	return ok
}

func (cache *sessionCache) hash(b []byte) []byte {
	r := md5.Sum(b)
	return r[:]
}

func hashToString(b []byte) string {
	return hex.EncodeToString(b[:])
}

func stringToHash(s string) []byte {
	r, _ := hex.DecodeString(s)
	return r
}
