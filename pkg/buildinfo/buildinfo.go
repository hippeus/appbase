package buildinfo

import (
	"fmt"
	"strings"
	"sync"
)

var (
	// version must be baked by extra linker flag via compilation:
	// -ldflags "-X '<path>.buildinfo.version=$(git describe --tags --always --dirty)'"
	version string
	// gitCommit must be baked by extra linker flag via compilation:
	// -ldflags "-X '<path>.buildinfo.commit=$(git rev-list -1 HEAD)'"
	commit string
	// buildTime must be baked by extra linker flag via compilation:
	// -ldflags "-X '<path>.buildinfo.buildTime=$(date -u '+%Y-%m-%d %H:%M:%S')'"
	buildTime string
)

// Build keeps meta information about the build typically linked at compilation.
// example:
//   go build -ldflags "-X 'github.com/hippeus/appbase/internal/buildinfo.version="0.1.0"'"
// If no variable is set via linker, Build prints "Version: unknown"
type Build struct{}

func (b Build) String() string {
	return impl.String()
}
func (b Build) VersionAsMap() map[string]interface{} {
	return impl.VersionAsMap()
}

var (
	impl = &buildInfo{
		vMap: map[string]interface{}{},
	}
)

type buildInfo struct {
	once sync.Once
	v    string
	vMap map[string]interface{}
}

func (b *buildInfo) VersionAsMap() map[string]interface{} {
	b.once.Do(lazyInit(b))
	return b.vMap
}

func (b *buildInfo) String() string {
	b.once.Do(lazyInit(b))
	return b.v
}

func lazyInit(b *buildInfo) func() {
	return func() {
		var sb strings.Builder

		if version != "" {
			b.vMap["version"] = version
			sb.WriteString(fmt.Sprintf("Version: %s", version))
		}
		if commit != "" {
			b.vMap["commit"] = commit
			if sb.Len() != 0 {
				sb.WriteRune('\n')
			}
			sb.WriteString(fmt.Sprintf("Commit: %s", commit))
		}
		if buildTime != "" {
			b.vMap["build_time"] = buildTime
			if sb.Len() != 0 {
				sb.WriteRune('\n')
			}
			sb.WriteString(fmt.Sprintf("Build Time: %s", buildTime))
		}
		if sb.Len() == 0 {
			sb.WriteString("Version: unknown")
		}
		b.v = sb.String()
	}
}
