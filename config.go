package bulletin

import (
	"flag"
)

type Config struct {
	// command line args
	MossPath string
}

// call DefineFlags before myflags.Parse()
func (c *Config) DefineFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.MossPath, "moss", "", "path to moss data dir (default .moss_bull_db)")
}

// call c.ValidateConfig() after myflags.Parse()
func (c *Config) ValidateConfig() error {
	if c.MossPath == "" {
		c.MossPath = ".moss_bull_db"
	}
	return nil
}
