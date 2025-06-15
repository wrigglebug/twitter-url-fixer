package main

import (
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
)

//go:embed favicon.ico
var iconData []byte

var (
	monitoringEnabled   int32 = 1
	selectedReplacement       = "vxtwitter"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("URL Replacer")
	systray.SetTooltip("Replaces x.com, twitter.com, and bsky.app links")

	mReplacements := systray.AddMenuItem("Set Replacement for x.com", "Choose replacement for x.com URLs")
	mFixVX := mReplacements.AddSubMenuItem("fixvx", "Replace x.com with fixvx")
	mVXTwitter := mReplacements.AddSubMenuItem("vxtwitter", "Replace x.com with vxtwitter")
	mGirlCockX := mReplacements.AddSubMenuItem("girlcockx", "Replace x.com with girlcockx")
	mStupidPenisX := mReplacements.AddSubMenuItem("stupidpenisx", "Replace x.com with stupidpenisx")

	mToggle := systray.AddMenuItem("Pause Monitoring", "Pause or Resume clipboard monitoring")

	mVXTwitter.Check()
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
			case <-mFixVX.ClickedCh:
				updateReplacement("fixvx", mFixVX, mVXTwitter, mGirlCockX, mStupidPenisX)
			case <-mVXTwitter.ClickedCh:
				updateReplacement("vxtwitter", mFixVX, mVXTwitter, mGirlCockX, mStupidPenisX)
			case <-mGirlCockX.ClickedCh:
				updateReplacement("girlcockx", mFixVX, mVXTwitter, mGirlCockX, mStupidPenisX)
			case <-mStupidPenisX.ClickedCh:
				updateReplacement("stupidpenisx", mFixVX, mVXTwitter, mGirlCockX, mStupidPenisX)
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
				text = reX.ReplaceAllString(text, selectedReplacement+".com")
				changed = true
			}

			if reTwitter.MatchString(text) {
				text = reTwitter.ReplaceAllString(text, selectedReplacement+".com")
				changed = true
			}

			if reBsky.MatchString(text) {
				text = reBsky.ReplaceAllString(text, "fxbsky.app")
				changed = true
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
