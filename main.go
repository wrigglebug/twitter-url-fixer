package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
	_ "embed"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
)

//go:embed favicon.ico
var iconData []byte

var (
	monitoringEnabled   int32 = 1
	selectedReplacement       = "vxtwitter"
)

const updateExeURL = "https://github.com/wrigglebug/twitter-url-fixer/releases/latest/download/twitter-url-fixer.exe"

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("URL Replacer")
	systray.SetTooltip("Replaces x.com, twitter.com, and bsky.app links")

		mReplacements := systray.AddMenuItem("Set Replacement for x.com", "Choose replacement for x.com URLs")
		mFixVX := mReplacements.AddSubMenuItem("fixvx", "Replace x.com with fixvx")
		mVXTwitter := mReplacements.AddSubMenuItem("vxtwitter", "Replace x.com with vxtwitter") // Default
		mGirlCockX := mReplacements.AddSubMenuItem("girlcockx", "Replace x.com with girlcockx")
		mStupidPenisX := mReplacements.AddSubMenuItem("stupidpenisx", "Replace x.com with stupidpenisx")

	mToggle := systray.AddMenuItem("Pause Monitoring", "Pause or Resume clipboard monitoring")

	mUpdate := systray.AddMenuItem("Check for Updates", "Check for the latest version")

	mVXTwitter.Check()

	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go monitorClipboard()
	go func() {
		checkForUpdates()

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
			case <-mFixVX.ClickedCh:
				updateReplacement("fixvx", mFixVX, mVXTwitter, mGirlCockX, mStupidPenisX)
			case <-mVXTwitter.ClickedCh:
				updateReplacement("vxtwitter", mFixVX, mVXTwitter, mGirlCockX, mStupidPenisX)
			case <-mGirlCockX.ClickedCh:
				updateReplacement("girlcockx", mFixVX, mVXTwitter, mGirlCockX, mStupidPenisX)
			case <-mStupidPenisX.ClickedCh:
				updateReplacement("stupidpenisx", mFixVX, mVXTwitter, mGirlCockX, mStupidPenisX)
			case <-mUpdate.ClickedCh:
				checkForUpdates()
			}
		}
	}()
}

func onExit() {}

func updateReplacement(choice string, mFixVX, mVXTwitter, mGirlCockX, mStupidPenisX *systray.MenuItem) {
	selectedReplacement = choice
	mFixVX.Uncheck()
	mVXTwitter.Uncheck()
	mGirlCockX.Uncheck()
	mStupidPenisX.Uncheck()

	switch choice {
	case "fixvx":
		mFixVX.Check()
	case "vxtwitter":
		mVXTwitter.Check()
	case "girlcockx":
		mGirlCockX.Check()
	case "stupidpenisx":
		mStupidPenisX.Check()
	}
}

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

			originalText := text
			changed := false

			if reX.MatchString(text) {
				text = reX.ReplaceAllStringFunc(text, func(url string) string {
					changed = true
					return strings.Replace(url, "x.com", selectedReplacement+".com", 1)
				})
			}

			if reTwitter.MatchString(text) {
				text = reTwitter.ReplaceAllStringFunc(text, func(url string) string {
					changed = true
					return strings.Replace(url, "twitter.com", selectedReplacement+".com", 1)
				})
			}

			if reBsky.MatchString(text) {
				text = reBsky.ReplaceAllStringFunc(text, func(url string) string {
					changed = true
					return strings.Replace(url, "bsky.app", "fxbsky.app", 1)
				})
			}

			if changed && text != originalText {
				if err := clipboard.WriteAll(text); err != nil {
					log.Printf("Failed to write to clipboard: %v", err)
				} else {
					fmt.Println("Replaced URLs in the clipboard.")
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func checkForUpdates() {
    log.Println("Checking for updates...")

    currentExe, err := os.Executable()
    if err != nil {
        log.Printf("Failed to get current executable path: %v", err)
        return
    }

    log.Printf("Current executable path: %s", currentExe)

    currentHash, err := computeFileHash(currentExe)
    if err != nil {
        log.Printf("Failed to compute hash of current executable: %v", err)
        return
    }

    log.Printf("Current executable hash: %s", currentHash)

    tempFile := filepath.Join(os.TempDir(), filepath.Base(updateExeURL))
    log.Printf("Temporary file path for update: %s", tempFile)

    if _, err := os.Stat(tempFile); err == nil {
        log.Printf("Existing update file found, removing: %s", tempFile)
        if err := os.Remove(tempFile); err != nil {
            log.Printf("Failed to remove existing file: %v", err)
            return
        }
    }

    out, err := os.Create(tempFile)
    if err != nil {
        log.Printf("Failed to create temporary file: %v", err)
        return
    }
    defer out.Close()

    log.Printf("Downloading update from %s", updateExeURL)
    resp, err := http.Get(updateExeURL)
    if err != nil {
        log.Printf("Failed to download update: %v", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("Failed to download update: HTTP Status Code %d", resp.StatusCode)
        return
    }

    if _, err := io.Copy(out, resp.Body); err != nil {
        log.Printf("Failed to save update: %v", err)
        return
    }

    log.Println("Update downloaded successfully.")

    newHash, err := computeFileHash(tempFile)
    if err != nil {
        log.Printf("Failed to compute hash of downloaded executable: %v", err)
        return
    }

    log.Printf("Downloaded executable hash: %s", newHash)

    if currentHash == newHash {
        log.Println("The downloaded executable is identical to the current version. No update needed.")
        return
    }


    log.Println("Applying the update...")
    applyUpdate(tempFile)
}

func computeFileHash(filePath string) (string, error) {
    log.Printf("Computing file hash for: %s", filePath)

    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hash := sha256.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }

    computedHash := fmt.Sprintf("%x", hash.Sum(nil))
    log.Printf("Computed hash: %s", computedHash)

    return computedHash, nil
}

func applyUpdate(newExe string) {
    currentExe, err := os.Executable()
    if err != nil {
        log.Printf("Failed to get current executable path: %v", err)
        return
    }

    log.Printf("Backing up current executable to: %s.bak", currentExe)
    backupExe := currentExe + ".bak"
    if err := os.Rename(currentExe, backupExe); err != nil {
        log.Printf("Failed to backup current executable: %v", err)
        return
    }

    log.Printf("Scheduling the replacement of the executable after the application closes.")
    go func() {
        time.Sleep(2 * time.Second)

        log.Printf("Replacing current executable with the new version.")
        if err := os.Rename(newExe, currentExe); err != nil {
            log.Printf("Failed to replace executable: %v", err)
            os.Rename(backupExe, currentExe)
            return
        }

        log.Printf("Restarting the application.")

        restartCmd := exec.Command(currentExe)
        restartCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

        if err := restartCmd.Start(); err != nil {
            log.Printf("Failed to restart application: %v", err)
            return
        }

        os.Exit(0)
    }()
}