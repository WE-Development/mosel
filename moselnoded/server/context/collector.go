package context

import (
	"github.com/WE-Development/mosel/api"
	"os"
	"os/exec"
	"bytes"
	"strings"
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

func (collector *collector) AddScript(name string, src string) error {
	mkdirIfNotExist(collector.scriptFolder, 0664)
	file, err := os.Open(collector.scriptFolder)

	if err != nil {
		return err
	}

	_, err = file.WriteString(src)
	collector.scripts = append(collector.scripts, name)
	return nil
}

func (collector *collector) FillNodeInfo(info *api.NodeInfo) {
	for _, script := range collector.scripts {
		executeScript(collector.scriptFolder + "/" + script, info)
	}
}

func executeScript(script string, info *api.NodeInfo) {
	cmd := exec.Command("/bin/bash", script)
	out := &bytes.Buffer{}
	cmd.Stdout = out

	res := make(map[string]string)
	for _, line := range
		strings.Split(out.String(), "\n") {
		parts := strings.SplitN(line, ":", 2)
		graph := parts[0]
		value := parts[1]
		res[graph] = value
	}
	(*info)[script] = res
}

func mkdirIfNotExist(path string, perm os.FileMode) bool {
	if ok, _ := exists(path); !ok {
		os.Mkdir(path, perm)
	}
	return false
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