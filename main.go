package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
)

var monitoringEnabled int32 = 1

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("URL Replacer")
	systray.SetTooltip("Replaces x.com and twitter.com links")

	mToggle := systray.AddMenuItem("Pause Monitoring", "Pause or Resume clipboard monitoring")
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
			}
		}
	}()
}

func onExit() {

}

func monitorClipboard() {
	reX := regexp.MustCompile(`https?://(?:www\.)?x\.com[^\s]*`)
	reTwitter := regexp.MustCompile(`https?://(?:www\.)?twitter\.com[^\s]*`)

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

			if err := clipboard.WriteAll(text); err != nil {
				log.Printf("Failed to write to clipboard: %v", err)
			} else {
				fmt.Println("Replaced URLs in the clipboard.")
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
