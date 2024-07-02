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
	systray.SetTooltip("Replaces x.com with fixvx.com")

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
	// Regular expression to find x.com URLs
	re := regexp.MustCompile(`https?://(?:www\.)?x\.com[^\s]*`)

	for {
		// Read current clipboard content
		text, err := clipboard.ReadAll()
		if err != nil {
			log.Printf("Failed to read clipboard: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Check if the clipboard contains an x.com URL
		if re.MatchString(text) {
			// Replace x.com with fixvx.com
			newText := re.ReplaceAllStringFunc(text, func(url string) string {
				return strings.Replace(url, "x.com", "fixvx.com", 1)
			})

			// Write the modified text back to the clipboard
			if err := clipboard.WriteAll(newText); err != nil {
				log.Printf("Failed to write to clipboard: %v", err)
			} else {
				fmt.Println("Replaced x.com URLs with fixvx.com in the clipboard.")
			}
		}

		// Sleep for a short duration before checking the clipboard again
		time.Sleep(2 * time.Second)
	}
}
