package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var monitoringEnabled int32 = 1

const updateExeURL = "https://github.com/wrigglebug/twitter-url-fixer/releases/latest/download/twitter-url-fixer.exe"

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData) // Replace with actual icon data
	systray.SetTitle("URL Replacer")
	systray.SetTooltip("Replaces x.com, twitter.com, and bsky.app links")

	mToggle := systray.AddMenuItem("Pause Monitoring", "Pause or Resume clipboard monitoring")
	mUpdate := systray.AddMenuItem("Check for Updates", "Check for the latest version")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go monitorClipboard()
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mToggle.ClickedCh:
				if atomic.LoadInt32(&monitoringEnabled) == 1 {
					atomic.StoreInt32(&monitoringEnabled, 0)
					mToggle.SetTitle("Resume Monitoring")
				} else {
					atomic.StoreInt32(&monitoringEnabled, 1)
					mToggle.SetTitle("Pause Monitoring")
				}
			case <-mUpdate.ClickedCh:
				checkForUpdates()
			}
		}
	}()
}

func onExit() {}

func monitorClipboard() {
	reX := regexp.MustCompile(`https?://(?:www\.)?x\.com[^\s]*`)
	reTwitter := regexp.MustCompile(`https?://(?:www\.)?twitter\.com[^\s]*`)
	reBsky := regexp.MustCompile(`https?://(?:www\.)?bsky\.app[^\s]*`)

	for {
		if atomic.LoadInt32(&monitoringEnabled) == 1 {
			text, err := clipboard.ReadAll()
			if err != nil {
				log.Printf("Failed to read clipboard: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}

			if reX.MatchString(text) {
				text = reX.ReplaceAllStringFunc(text, func(url string) string {
					return strings.Replace(url, "x.com", "vxtwitter.com", 1)
				})
			}

			if reTwitter.MatchString(text) {
				text = reTwitter.ReplaceAllStringFunc(text, func(url string) string {
					return strings.Replace(url, "twitter.com", "vxtwitter.com", 1)
				})
			}

			if reBsky.MatchString(text) {
				text = reBsky.ReplaceAllStringFunc(text, func(url string) string {
					return strings.Replace(url, "bsky.app", "fxbsky.app", 1)
				})
			}

			if err := clipboard.WriteAll(text); err != nil {
				log.Printf("Failed to write to clipboard: %v", err)
			} else {
				fmt.Println("Replaced URLs in the clipboard.")
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func checkForUpdates() {
	currentExe, err := os.Executable()
	if err != nil {
		log.Printf("Failed to get current executable path: %v", err)
		return
	}

	currentHash, err := computeFileHash(currentExe)
	if err != nil {
		log.Printf("Failed to compute hash of current executable: %v", err)
		return
	}

	tempFile := filepath.Join(os.TempDir(), filepath.Base(updateExeURL))
	out, err := os.Create(tempFile)
	if err != nil {
		log.Printf("Failed to create temporary file: %v", err)
		return
	}
	defer out.Close()

	resp, err := http.Get(updateExeURL)
	if err != nil {
		log.Printf("Failed to download update: %v", err)
		return
	}
	defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		log.Printf("Failed to save update: %v", err)
		return
	}

	newHash, err := computeFileHash(tempFile)
	if err != nil {
		log.Printf("Failed to compute hash of downloaded executable: %v", err)
		return
	}

	if currentHash == newHash {
		log.Println("The downloaded executable is identical to the current version.")
		return
	}

	promptUpdate(tempFile)
}

func computeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func promptUpdate(newExe string) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	Dialog{
		AssignTo: &dlg,
		Title:    "Update Available",
		MinSize:  Size{300, 150},
		Layout:   VBox{},
		Children: []Widget{
			Label{
				Text: "A new version is available. Do you want to update?",
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     "Yes",
						OnClicked: func() {
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "No",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(nil)

	if dlg.Result() == walk.DlgCmdOK {
		applyUpdate(newExe)
	}
}

func applyUpdate(newExe string) {
	currentExe, err := os.Executable()
	if err != nil {
		log.Printf("Failed to get current executable path: %v", err)
		return
	}

	backupExe := currentExe + ".bak"
	if err := os.Rename(currentExe, backupExe); err != nil {
		log.Printf("Failed to backup current executable: %v", err)
		return
	}

	if err := os.Rename(newExe, currentExe); err != nil {
		log.Printf("Failed to replace executable: %v", err)
		os.Rename(backupExe, currentExe) // Restore backup
		return
	}

	restartCmd := exec.Command(currentExe)
	if err := restartCmd.Start(); err != nil {
		log.Printf("Failed to restart application: %v", err)
		return
	}

	os.Exit(0) // Exit current instance
}
