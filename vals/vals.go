package vals

import (
	"net/url"

	"github.com/platform-engineering-labs/orbital/opm/security"
	"github.com/platform-engineering-labs/orbital/opm/tree"
	"github.com/platform-engineering-labs/orbital/ops"
	"github.com/platform-engineering-labs/orbital/platform"
)

const (
	ManagedRoot = "/opt/pel"
)

var TreeConfig = &tree.Config{
	OS:       platform.Current().OS,
	Arch:     platform.Current().Arch,
	Security: security.Default,
	Repositories: []ops.Repository{
		{
			Uri: url.URL{
				Scheme:   "https",
				Host:     "hub.platform.engineering",
				Path:     "/repos/platform.engineering/pel",
				Fragment: "stable",
			},
			Priority: 0,
			Enabled:  true,
		},
		{
			Uri: url.URL{
				Scheme:   "https",
				Host:     "hub.platform.engineering",
				Path:     "/repos/platform.engineering/community",
				Fragment: "stable",
			},
			Priority: 1,
			Enabled:  true,
		},
	},
}
