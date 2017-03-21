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
package moselserver

import (
	"gopkg.in/mgo.v2"
)

func (server *MoselServer) initMongo(driverName string, dataSourceName string) (MongoDataSource, error) {
	var session *mgo.Session
	var err error

	if session, err = mgo.Dial(dataSourceName); err != nil {
		return nil, err
	}

	if err = session.Ping(); err != nil {
		return nil, err
	}

	return NewMongoDataSource(driverName, session), nil
}