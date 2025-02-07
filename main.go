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
	files           = flag.String("f", ".", "Root")
	onlyDirectories = flag.Bool("d", false, "Only directories")
	all             = flag.Bool("a", false, "All files")
	comment         = flag.Bool("c", false, "Add a comment prefix on each line")
	indent          = flag.Int("i", 0, "Indent each line")
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
	startDir := common.CleanPath(*files)
	mask := "*"

	if !common.IsDirectory(startDir) {
		mask = filepath.Base(startDir)
		startDir = filepath.Dir(startDir)
	}

	lines := []string{}
	err := common.WalkFiles(filepath.Join(startDir, mask), true, true, func(path string, fi os.FileInfo) error {
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

		sb.WriteString(strings.Repeat(" ", *indent))

		switch level {
		case 0:
		case 1:
			sb.WriteString(connector)
		default:
			sb.WriteString(strings.Repeat(tab, level-1))
			sb.WriteString(connector)
		}

		sb.WriteString(name)

		lines = append(lines, sb.String())

		return nil
	})

	if common.Error(err) {
		return err
	}

	mx := 0
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

		mx = common.Max(mx, len(lines[i]))
	}

	format := fmt.Sprintf("%%-%ds//\n", mx+3)

	for i := range len(lines) {
		if *comment {
			fmt.Printf(format, lines[i])
		} else {
			fmt.Printf("%s\n", lines[i])
		}

	}
	return nil
}

func main() {
	common.Run(nil)
}
