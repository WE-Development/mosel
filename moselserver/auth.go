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
	"io"
)

type authProvider interface {
	Authenticate(user string, passwd string) bool
}

//this provider will always authenticate successfully (for testing purposes)
type AuthTrue struct {

}

//implement authProvider
func (auth *AuthTrue) Authenticate(user string, passwd string) bool {
	return true
}

type AuthSys struct {
	AllowedUsers []string
}

func (auth AuthSys)  Authenticate(user string, passwd string) bool {
	if (len(auth.AllowedUsers) == 0) {
		log.Println("AuthSys is configured but no users are allowed")
		return false
	}

	cmd := exec.Command("su", user);
	in, err := cmd.StdinPipe()

	if err != nil {
		log.Printf("Failed to perform AuthSys! Stdin could't be opend: %s", err)
		return false
	}

	cmd.Run()
	io.WriteString(in, passwd)
	io.WriteString(in, "\n")
	cmd.Wait()

	success := cmd.ProcessState.Success()
	if success {
		log.Printf("Authenticated user %s via AuthSys", user)
	} else {
		log.Printf("Rejected user %s via AuthSys", user)
	}

	return success;
}