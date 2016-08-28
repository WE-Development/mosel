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
package commons

import (
	"os/exec"
	"bytes"
	"log"
	"strings"
)

func ExecuteScript(arg ...string) (map[string]string, error) {
	cmd := exec.Command("/bin/bash", arg...)
	out := &bytes.Buffer{}
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	res := make(map[string]string)
	for _, line := range strings.Split(out.String(), "\n") {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			log.Printf("Invalid grap data '%s'", line)
			continue
		}

		graph := parts[0]
		value := parts[1]
		res[graph] = value
	}
	return res, nil
}
