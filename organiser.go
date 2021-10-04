package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
)

const detaultDateFormat = "2006-01-02"

func organise(inputFolder, outputFolder, copyOrMove string, w *fyne.Window, p *widget.ProgressBar) error {
	var (
		count            int
		acceptedSuffixes = []string{"jpeg", "jpg", "png", "gif", "raw", "svg", "heif", "heic", "bmp", "MP4"}
	)

	_, f, err := folderCheck(inputFolder)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}
	outputFolder = f + "/output"

	fmt.Printf("Let's start! Input folder: %s Output folder %s \n", inputFolder, outputFolder)
	os.RemoveAll(outputFolder)

	err = filepath.Walk(inputFolder,
		func(path string, file os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			count++

			name := file.Name()
			if !isExtensionAccepted(name, acceptedSuffixes) {
				return nil
			}

			date := file.ModTime().Format(detaultDateFormat)
			year := file.ModTime().Year()

			newFolder := fmt.Sprintf("%s/%d/%s", outputFolder, year, date)

			fmt.Printf("Action: %s, File: %s, Folder: %s \n", copyOrMove, name, newFolder)

			err = os.MkdirAll(newFolder, 0755)
			if err != nil {
				log.Printf("ERROR: %s", err)
				return nil
			}
			if copyOrMove == "copy" {
				err := Copy(path, newFolder+"/"+name)
				if err != nil {
					log.Printf("ERROR: %s", err)
					return nil
				}

			} else {
				err := Move(path, newFolder+"/"+name)
				if err != nil {
					log.Printf("ERROR: %s", err)
					return nil
				}
			}

			v := setBarValue(count)
			if p != nil {
				p.SetValue(v)
			}
			return nil
		})
	if err != nil {
		log.Printf("ERROR: %s", err)
		return nil
	}

	if p != nil {
		p.SetValue(1)
	}

	successMessage := fmt.Sprintf("Managed to %s %d files", copyOrMove, count)
	if w != nil {
		dialog.ShowInformation("Done!", successMessage, *w)
	}

	fmt.Printf("Action: %s, number of files: %d \n", copyOrMove, count)

	return nil
}

func isExtensionAccepted(filename string, suffixes []string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix("."+strings.ToLower(filename), s) {
			return true
		}
	}
	return false
}

// Copy is for copying files
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// Move is for moving files
func Move(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	err = os.Rename(src, dst)
	if err != nil {
		log.Fatal(err)
	}
	return in.Close()
}

func folderCheck(f string) (os.FileInfo, string, error) {
	folderInfo, err := os.Stat(f)
	if os.IsNotExist(err) {
		log.Fatal("Folder does not exist.")
		return nil, f, err
	}
	// if it's not a folder, then it's a file and we get the folder that contains it
	if !folderInfo.IsDir() {
		pathList := strings.Split(f, "/")
		f = strings.Join(pathList[:len(pathList)-1], "/")
		folderInfo, err = os.Stat(f)
		if os.IsNotExist(err) {
			log.Fatal("Folder does not exist.")
			return nil, f, err
		}
	}
	log.Println(folderInfo)
	return folderInfo, f, nil
}
