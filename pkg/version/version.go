package version

import "runtime/debug"

// version is set at build time by GoReleaser or manual ldflags:
//
//	go build -ldflags "-X github.com/Dicklesworthstone/beads_viewer/pkg/version.version=v1.2.3"
//
// It starts empty so init() can distinguish "ldflags set it" from "no injection".
var version string

// fallback is the hardcoded version kept in sync with the latest release tag.
// Used only when both ldflags and debug.ReadBuildInfo fail to provide a version.
const fallback = "v0.14.4"

// Version is the resolved application version, populated by init().
var Version string

func init() {
	switch {
	case version != "":
		// 1. Build-time ldflags injection (GoReleaser, Nix, manual).
		Version = version
	case versionFromBuildInfo() != "":
		// 2. Module version from "go install ...@vX.Y.Z".
		Version = versionFromBuildInfo()
	default:
		// 3. Hardcoded fallback (always available, manually bumped per release).
		Version = fallback
	}
}

// versionFromBuildInfo extracts the module version stamped by the Go toolchain
// when the binary is built via "go install". Returns empty string if unavailable.
func versionFromBuildInfo() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	v := info.Main.Version
	if v == "" || v == "(devel)" {
		return ""
	}
	if v[0] != 'v' {
		v = "v" + v
	}
	return v
}
