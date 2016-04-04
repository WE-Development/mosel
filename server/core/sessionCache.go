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
package core

import (
	"time"
	"math/rand"
	"crypto/md5"
	"encoding/hex"
)

type sessionCache struct {

}

type session struct {
	keyHash []byte
	validTo time.Time
}

func (cache sessionCache) NewSession(ctx MoselServerContext, user string) (string, time.Time) {
	s := session{}
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	b := make([]byte, 256)
	rand.Read(b)

	key := cache.hash(
		[]byte(string(millis) + user + string(b)))

	keyHash := cache.hash([]byte(key[:]))

	s.keyHash = []byte(keyHash[:])
	s.validTo = time.Now()
	return hex.EncodeToString(key[:]), time.Now()
}

func (cache sessionCache) ValidateSession(key string) bool {
	return true
}

func (cache sessionCache) hash(b []byte) []byte {
	r := md5.Sum(b)
	return r[:]
}