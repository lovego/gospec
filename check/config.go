package check

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/lovego/gospec/check/names"
	"github.com/lovego/gospec/check/orders"
	"github.com/lovego/gospec/check/sizes"
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
		Sizes  *sizes.ConfigT
		Names  *names.ConfigT
		Funcs  *names.FuncConfigT
		Orders *orders.OrdersConfigT
	}{
		Sizes:  &sizes.Config,
		Names:  &names.Config,
		Funcs:  &names.FuncConfig,
		Orders: &orders.Config,
	}
	if err := json.Unmarshal(content, config); err != nil {
		panic(err)
	}
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
