# twitter-url-fixer
Automatically converts x.com links to fixvx.com

### Usage
 - run the executable
 - copy an x.com link
 - the script automatically converts it to fixvx.com for embedding in discord
 - it can be closed from the system tray
### Building
 ```bash
 # clone the repository
 $ https://github.com/wrigglebug/twitter-url-fixer.git
 # install prerequisites
 $ go get github.com/atotto/clipboard
 $ go get github.com/getlantern/systray
 # build for windows
 $ go build -ldflags="-H windowsgui"
 ```
