# twitter-url-fixer
Automatically converts twitter / x.com links to vxtwitter / fixvx.com

### Usage
 - [Download the latest release](https://github.com/wrigglebug/twitter-url-fixer/releases/latest/download/twitter-url-fixer.exe)
 - Run twitter-url-fixer.exe (it runs minimized to your system tray)
 - Copy a twitter or X.com link
 - The program will automatically convert the twitter / X links in your clipboard to vxtwitter / fixvx.com for embedding in discord
 - You can right click the icon in the system tray to toggle clipboard replacement on/off
### Building
 ```bash
 # clone the repository
 $ git clone https://github.com/wrigglebug/twitter-url-fixer.git
 # install prerequisites
 $ go get github.com/atotto/clipboard
 $ go get github.com/getlantern/systray
 # build for windows
 $ go build -ldflags="-H windowsgui"
 ```
