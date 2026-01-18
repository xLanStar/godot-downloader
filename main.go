package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/xlanstar/godot-downloader/internal/args"
	"github.com/xlanstar/godot-downloader/internal/downloader"
	"github.com/xlanstar/godot-downloader/internal/parser"

	"github.com/alexflint/go-arg"
)

func init() {
	arg.MustParse(&args.Args)
}

func main() {
	args.ResolveArgs()
	log.Printf("%v\n", args.Args)
	version, slug := args.ParseGodotVersionAndSlug(args.Args.GodotVersion)
	log.Printf("Version: %s, Slug: %s\n", version, slug)
	if version == "" || slug == "" {
		fmt.Println("Invalid version format")
		return
	}

	url := parser.GetGodotDownloadURL(parser.GodotSearchOptions{
		Version:  version,
		Slug:     slug,
		Platform: string(args.Args.GodotPlatform),
		Mono:     args.Args.GodotMono,
	})
	if len(url) == 0 {
		fmt.Println("Cannot find the download URL you specified")
		return
	}

	file, err := downloader.DownloadURL(url, args.Args.OutputFilePath)
	if err != nil {
		fmt.Printf("Cannot download the file: %v\n", err)
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Failed to close downloaded file: %v", closeErr)
		}
	}()

	if args.Args.Unarchive {
		zipFile, err := zip.OpenReader(args.Args.OutputFilePath)
		if err != nil {
			fmt.Printf("Cannot open the zip file: %v\n", err)
			return
		}
		defer func() {
			if closeErr := zipFile.Close(); closeErr != nil {
				log.Printf("Failed to close zip file: %v", closeErr)
			}
		}()

		for _, file := range zipFile.File {
			if file.FileInfo().IsDir() {
				err = os.MkdirAll(filepath.Dir(file.Name), os.ModePerm)
				if err != nil {
					fmt.Printf("Cannot create the directory: %v\n", err)
					break
				}
				continue
			}

			dstFile, err := os.OpenFile(file.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				fmt.Printf("Cannot create the destination file: %v\n", err)
				break
			}
			defer func() {
				if closeErr := dstFile.Close(); closeErr != nil {
					log.Printf("Failed to close destination file: %v", closeErr)
				}
			}()

			fileInArchive, err := file.Open()
			if err != nil {
				fmt.Printf("Cannot open the file in the zip: %v\n", err)
				break
			}
			defer func() {
				if closeErr := fileInArchive.Close(); closeErr != nil {
					log.Printf("Failed to close file in archive: %v", closeErr)
				}
			}()

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				fmt.Printf("Cannot copy the file: %v\n", err)
				break
			}
		}
	}
}
