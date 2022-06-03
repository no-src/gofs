package about

import (
	"github.com/no-src/gofs/version"
	"github.com/no-src/log"
)

const logo = `
   ________  ________  ________ ________      
  |\   ____\|\   __  \|\  _____\\   ____\     
  \ \  \___|\ \  \|\  \ \  \__/\ \  \___|_    
   \ \  \  __\ \  \\\  \ \   __\\ \_____  \   
    \ \  \|\  \ \  \\\  \ \  \_| \|____|\  \  
     \ \_______\ \_______\ \__\    ____\_\  \ 
      \|_______|\|_______|\|__|   |\_________\
                                  \|_________|

`
const (
	openSourceUrl    = "https://github.com/no-src/gofs"
	documentationUrl = "https://pkg.go.dev/github.com/no-src/gofs@" + version.VERSION
	releaseUrl       = "https://github.com/no-src/gofs/releases"
	dockerImageUrl   = "https://hub.docker.com/r/nosrc/gofs"
)

// PrintAbout print the program logo and basic info
func PrintAbout() {
	log.Log(logo)
	log.Log("The gofs is a file synchronization tool out of the box based on golang")
	log.Log("Open source repository at: <%s>", openSourceUrl)
	log.Log("Download the latest version at: <%s>", releaseUrl)
	log.Log("The docker image repository address at: <%s>", dockerImageUrl)
	log.Log("Full documentation at: <%s>", documentationUrl)
}
