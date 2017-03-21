/*
 * Copyright 2017 Robin Engel
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
package context

import (
	"gopkg.in/mgo.v2"
	"time"
	"github.com/bluedevel/mosel/api"
)

type mongoDataPersistence struct {
	session       *mgo.Session
	serverContext *MoseldServerContext
}

func NewMongoDataPersistence(session *mgo.Session) *mongoDataPersistence {
	return &mongoDataPersistence{
		session:session,
	}
}

func (pers *mongoDataPersistence) Init() error {
	return nil
}

func (pers *mongoDataPersistence) Add(node string, t time.Time, info api.NodeInfo) {

}

func (pers *mongoDataPersistence) GetAll() (DataCacheStorage, error) {
	return nil, nil
}

func (pers *mongoDataPersistence) GetAllSince(since time.Duration) (DataCacheStorage, error) {
	return nil, nil
}
