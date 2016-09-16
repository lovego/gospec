package check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/bughou-go/spec/check/names"
	"github.com/bughou-go/spec/check/sizes"
)

func init() {
	parseConfig()
}

// parse config from file
func parseConfig() {
	p := configFilePath()
	if p == `` {
		return
	}
	if content, err := ioutil.ReadFile(p); err == nil {
		parseConfigContent(content)
	} else {
		panic(err)
	}
}

func parseConfigContent(content []byte) {
	var config = &struct {
		Sizes *sizes.TConfig
		Names *names.TConfig
	}{
		Sizes: &sizes.Config,
		Names: &names.Config,
	}
	if err := json.Unmarshal(content, config); err != nil {
		panic(err)
	}
	fmt.Println(config)
}

func configFilePath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for ; dir != `/`; dir = path.Dir(dir) {
		if p := testConfig(dir); p != `` {
			return p
		}
	}
	if p := testConfig(dir); p != `` {
		return p
	}
	return ``
}

func testConfig(dir string) string {
	const file = `gospec.json`
	p := path.Join(dir, file)
	if _, err := os.Stat(p); err == nil {
		return p
	} else if os.IsNotExist(err) {
		return ``
	} else {
		panic(err)
	}
}
