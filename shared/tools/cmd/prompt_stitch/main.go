package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	masterPath := flag.String("master", "shared/prompt/experiments/master.poml", "Path to the primary master prompt (set empty to skip).")
	subDir := flag.String("subdir", "shared/prompt/experiments/sub_actions", "Directory containing sub-action prompts.")
	outputPath := flag.String("output", "shared/prompt/experiments/master_compiled.poml", "Destination file for the stitched prompt.")
	flag.Parse()

	filesToStitch := make([]string, 0, 16)
	if *masterPath != "" {
		filesToStitch = append(filesToStitch, *masterPath)
	}

	subFiles, err := collectPomlFiles(*subDir)
	if err != nil {
		fatalErr(fmt.Errorf("collect sub prompts: %w", err))
	}
	filesToStitch = append(filesToStitch, subFiles...)

	if len(filesToStitch) == 0 {
		fatalErr(fmt.Errorf("no prompt files discovered; check master and subdir flags"))
	}

	var buf bytes.Buffer
	for _, path := range filesToStitch {
		content, err := os.ReadFile(path)
		if err != nil {
			fatalErr(fmt.Errorf("read %s: %w", path, err))
		}

		if buf.Len() > 0 && !bytes.HasSuffix(buf.Bytes(), []byte("\n\n")) {
			buf.WriteString("\n")
		}

		buf.WriteString(fmt.Sprintf("<!-- BEGIN %s -->\n", path))
		buf.Write(content)
		if !bytes.HasSuffix(content, []byte("\n")) {
			buf.WriteString("\n")
		}
		buf.WriteString(fmt.Sprintf("<!-- END %s -->\n\n", path))
	}

	if err := os.WriteFile(*outputPath, buf.Bytes(), 0o644); err != nil {
		fatalErr(fmt.Errorf("write output %s: %w", *outputPath, err))
	}

	fmt.Printf("Stitched %d prompt files → %s\n", len(filesToStitch), *outputPath)
}

func collectPomlFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".poml") {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	sort.Strings(files)
	return files, nil
}

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "prompt_stitch: %v\n", err)
	os.Exit(1)
}



