# Grinnell Writing Lab
It is a Golang script that check the availability of spots for you so that you do not need to refresh page every ten second. It will send you a system-level notification when there exist a avaliable spot. To prevent exploitation, the the script will **NOT** reserve the spot for you automatically.

## Getting Started
### Configuration
Change **YOUR_EMAIL** and **PASSWORD** to your own Grinnell College email address and writing lab password. You can also change the **INTERVAL**, which controls the frequency of script checking the spot availability.
```Go
const USER_NAME string = "YOUR_EMAIL"
const PASSWORD string = "PASSWORD"
const INTERVAL time.Duration = 1
```

### Build & Run 
Execute ```go build``` to build an executable binary for your system. 
