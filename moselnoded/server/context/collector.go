package context

import (
	"github.com/WE-Development/mosel/api"
	"os"
	"os/exec"
	"bytes"
	"strings"
	"log"
	"io/ioutil"
)

type collector struct {
	scriptFolder string
	scripts      []string
}

func NewCollector() *collector {
	return &collector{
		scriptFolder: "/tmp/mosel",
		scripts: make([]string, 0),
	}
}

func (collector *collector) AddScript(name string, src []byte) error {
	filePath := collector.scriptFolder + "/" + name

	if _, err := mkdirIfNotExist(collector.scriptFolder, 0764); err != nil {
		return err
	}

	err := ioutil.WriteFile(filePath, src, 0664)

	if err != nil {
		log.Println(err)
		return err
	}

	collector.scripts = append(collector.scripts, name)
	log.Printf("Added script %s", name)
	return nil
}

func (collector *collector) FillNodeInfo(info *api.NodeInfo) {
	for _, script := range collector.scripts {
		collector.executeScript(script, info)
	}
}

func (collector *collector) executeScript(name string, info *api.NodeInfo) {
	script := collector.scriptFolder + "/" + name
	cmd := exec.Command("/bin/bash", script)
	out := &bytes.Buffer{}
	cmd.Stdout = out

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
	(*info)[name] = res
}

func mkdirIfNotExist(path string, perm os.FileMode) (bool, error) {
	if ok, _ := exists(path); !ok {
		err := os.Mkdir(path, perm)
		return err != nil, err
	}
	return false, nil
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}