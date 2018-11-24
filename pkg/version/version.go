package version

import (
	"fmt"
	"runtime"
)

// LBaaS version. Update this whenever making a new release.

var (
	gitVersion string = "0.0.1"
	gitCommit  string = ""
	buildDate  string = "1970-01-01T00:00:00Z"
)

type Info struct {
	GitVersion string `json:"gitVersion"`
	GitCommit  string `json:"gitCommit"`
	BuildDate  string `json:"buildDate"`
	Goversion  string `json:"goVersion"`
	Compiler   string `json:"compiler"`
	Platform   string `json:"platform"`
}

func Get() Info {
	return Info{
		GitVersion: gitVersion,
		GitCommit:  gitCommit,
		BuildDate:  buildDate,
		Goversion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func (info Info) String() string {
	return info.GitVersion
}

