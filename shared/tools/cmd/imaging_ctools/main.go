package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	op := flag.String("op", "resize", "Operation to run: resize, thumbnail, convert, optimize")
	input := flag.String("input", "", "Path to the source image")
	output := flag.String("output", "", "Path for the processed image")
	size := flag.String("size", "1024x1024", "Geometry to use for resize/thumbnail operations (e.g. 800x600, 50%)")
	quality := flag.Int("quality", 85, "JPEG/WebP quality for optimize operation (1-100)")
	magickBin := flag.String("magick", "", "Override path to the `magick` binary (defaults to env/auto-detect)")
	flag.Parse()

	if *input == "" || *output == "" {
		fatalErr(errors.New("input and output flags are required"))
	}

	magickPath, err := resolveMagickBinary(*magickBin)
	if err != nil {
		fatalErr(err)
	}

	args, err := buildArgs(*op, *input, *output, *size, *quality)
	if err != nil {
		fatalErr(err)
	}

	cmd := exec.Command(magickPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("[imaging_ctools] %s %v\n", magickPath, args)
	if err := cmd.Run(); err != nil {
		fatalErr(fmt.Errorf("magick command failed: %w", err))
	}
}

func buildArgs(op, input, output, size string, quality int) ([]string, error) {
	switch op {
	case "resize":
		return []string{input, "-resize", size, output}, nil
	case "thumbnail":
		return []string{input, "-thumbnail", size, "-strip", output}, nil
	case "convert":
		return []string{input, output}, nil
	case "optimize":
		if quality < 1 || quality > 100 {
			return nil, fmt.Errorf("quality must be between 1 and 100 (got %d)", quality)
		}
		return []string{input, "-strip", "-quality", fmt.Sprintf("%d", quality), output}, nil
	default:
		return nil, fmt.Errorf("unsupported op %q", op)
	}
}

func resolveMagickBinary(explicit string) (string, error) {
	candidates := []string{}
	if explicit != "" {
		candidates = append(candidates, explicit)
	}
	if env := os.Getenv("IMAGEMAGICK_BIN"); env != "" {
		candidates = append(candidates, env)
	}
	candidates = append(candidates,
		filepath.Join("c", "ImageMagick", "utilities", "magick"),
		"magick",
		filepath.Join("c", "ImageMagick", "MagickWand", "magick"),
	)

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		path, err := exec.LookPath(candidate)
		if err == nil {
			return path, nil
		}
		if fileExists(candidate) {
			return candidate, nil
		}
	}

	return "", errors.New("magick binary not found. Set IMAGEMAGICK_BIN or pass -magick pointing to c/ImageMagick build")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func fatalErr(err error) {
	fmt.Fprintf(os.Stderr, "imaging_ctools: %v\n", err)
	os.Exit(1)
}




