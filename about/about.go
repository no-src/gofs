package about

import (
	"github.com/no-src/gofs/version"
	"github.com/no-src/log"
)

const Logo = `
   ________  ________  ________ ________      
  |\   ____\|\   __  \|\  _____\\   ____\     
  \ \  \___|\ \  \|\  \ \  \__/\ \  \___|_    
   \ \  \  __\ \  \\\  \ \   __\\ \_____  \   
    \ \  \|\  \ \  \\\  \ \  \_| \|____|\  \  
     \ \_______\ \_______\ \__\    ____\_\  \ 
      \|_______|\|_______|\|__|   |\_________\
                                  \|_________|

`
const OpenSourceUrl = "https://github.com/no-src/gofs"
const DocumentationUrl = "https://pkg.go.dev/github.com/no-src/gofs@" + version.VERSION
const ReleaseUrl = "https://github.com/no-src/gofs/releases"

func PrintAbout() {
	log.Log(Logo)
	log.Log("The gofs is a file synchronization tool out of the box based on golang")
	log.Log("Open source repository at: <%s>", OpenSourceUrl)
	log.Log("Download the latest version at: <%s>", ReleaseUrl)
	log.Log("Full documentation at: <%s>", DocumentationUrl)
}
