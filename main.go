package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"os"
	"path/filepath"
	"strings"
)

//go:embed go.mod
var resources embed.FS

var (
	root            = flag.String("r", ".", "Root")
	onlyDirectories = flag.Bool("d", false, "only directories")
	all             = flag.Bool("a", false, "all files")
)

const (
	lastConnector = "`-- "
	connector     = "|-- "
	line          = "|   "
	tab           = "    "
)

func init() {
	common.Init("", "", "", "", "", "", "", "", &resources, nil, nil, run, 0)
}

func run() error {
	startDir := common.CleanPath(*root)

	lines := []string{}
	common.WalkFiles(filepath.Join(startDir, "*"), true, true, func(path string, fi os.FileInfo) error {
		if *onlyDirectories && !fi.IsDir() {
			return nil

		}

		name := filepath.Base(path)

		if !*all && strings.HasPrefix(name, ".") {
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		level := strings.Count(path[len(startDir):], "/") + strings.Count(path, "\\")

		sb := strings.Builder{}
		switch level {
		case 0:
		case 1:
			sb.WriteString(connector)
		default:
			sb.WriteString(strings.Repeat(tab, level-2))
			sb.WriteString(connector)
		}

		sb.WriteString(name)

		lines = append(lines, sb.String())

		return nil
	})

	for i := 0; i < len(lines); i++ {
		p := strings.Index(lines[i], connector)
		if p == -1 {
			continue
		}

		b := i+1 == len(lines)
		if !b {
			pn := strings.Index(lines[i+1], connector)
			b = p != pn
		}

		if b {
			lines[i] = lines[i][:p] + lastConnector + lines[i][p+len(lastConnector):]
		}

		for j := i - 1; j >= 0; j-- {
			if p+len(line) < len(lines[j]) && lines[j][p:p+len(line)] == tab {
				lines[j] = lines[j][:p] + line + lines[j][p+len(line):]
			} else {
				break
			}
		}
	}

	fmt.Printf("%s\n", strings.Join(lines, "\n"))

	return nil
}

func main() {
	common.Run(nil)
}
