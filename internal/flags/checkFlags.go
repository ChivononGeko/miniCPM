package flags

import "os"

func ArgsCheck(args []string) bool {
	if len(args) > 4 {
		return false
	}

	if len(args) == 1 && !(args[0] == "-help" || args[0] == "--help") {
		return false
	}

	if len(args) == 2 && !(args[0] == "-port" || args[0] == "--port" || args[0] == "-dir" || args[0] == "--dir") {
		return false
	}

	if len(args) == 3 {
		return false
	}

	if len(args) == 4 {
		if !(args[0] == "-port" || args[0] == "--port" || args[0] == "-dir" || args[0] == "--dir") {
			return false
		}
		if !(args[2] == "-port" || args[2] == "--port" || args[2] == "-dir" || args[2] == "--dir") {
			return false
		}
	}

	return true
}

func PortDirChecks(port int, dir string) bool {
	if !(port > 1024 && port <= 65535) {
		return false
	}

	if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
		return false
	}

	return true
}
