package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("URL Replacer")
	systray.SetTooltip("Replaces x.com and twitter.com links")

	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go monitorClipboard()

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {

}

func monitorClipboard() {
	// Regular expressions to find x.com and twitter.com URLs
	reX := regexp.MustCompile(`https?://(?:www\.)?x\.com[^\s]*`)
	reTwitter := regexp.MustCompile(`https?://(?:www\.)?twitter\.com[^\s]*`)

	for {
		// Read current clipboard content
		text, err := clipboard.ReadAll()
		if err != nil {
			log.Printf("Failed to read clipboard: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Replace x.com with fixvx.com
		if reX.MatchString(text) {
			text = reX.ReplaceAllStringFunc(text, func(url string) string {
				return strings.Replace(url, "x.com", "fixvx.com", 1)
			})
		}

		// Replace twitter.com with vxtwitter.com
		if reTwitter.MatchString(text) {
			text = reTwitter.ReplaceAllStringFunc(text, func(url string) string {
				return strings.Replace(url, "twitter.com", "vxtwitter.com", 1)
			})
		}

		// Write the modified text back to the clipboard
		if err := clipboard.WriteAll(text); err != nil {
			log.Printf("Failed to write to clipboard: %v", err)
		} else {
			fmt.Println("Replaced URLs in the clipboard.")
		}

		time.Sleep(1 * time.Second)
	}
}
