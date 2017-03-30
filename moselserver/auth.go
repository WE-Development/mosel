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
	"log"
	"os/exec"
	"bytes"
	"time"

	"github.com/bluedevel/mosel/config"
)

// Interface for validating a user-password combination
type authProvider interface {
	Authenticate(user string, passwd string) bool
}

// Will always authenticate successfully. (for testing purposes!)
type AuthTrue struct {
}

// implement authProvider
func (auth *AuthTrue) Authenticate(user string, passwd string) bool {
	return true
}

// Will check a static set of user-password combinations
type AuthStatic struct {
	Users map[string]*moselconfig.UserConfig
}

// implement authProvider
func (auth *AuthStatic) Authenticate(user string, passwd string) bool {
	conf, ok := auth.Users[user]

	if !ok {
		//no user with that name is registered
		return false
	}

	return passwd == conf.Password
}

//Auth Sys
type AuthSys struct {
	AllowedUsers []string
}

//todo find a clever solution for this
func (auth AuthSys) Authenticate(userName string, passwd string) bool {
	if (len(auth.AllowedUsers) == 0) {
		log.Println("AuthSys is configured but no users are allowed")
		return false
	}

	//cmd := exec.Command("su", userName);
	cmd := exec.Command("true", userName);
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	cmd.Stderr = out
	cmd.Stdin = in

	/*if err != nil {
		log.Printf("Failed to perform AuthSys! Stdin could't be opend: %s", err)
		return false
	}*/

	cmd.Run()
	time.Sleep(1 * time.Second)
	in.WriteString(passwd + "\n")
	//close(in)
	cmd.Wait()

	log.Println(out.String())

	success := cmd.ProcessState.Success()
	if success {
		log.Printf("Authenticated user %s via AuthSys", userName)
	} else {
		log.Printf("Rejected user %s via AuthSys", userName)
	}

	return success;
}
