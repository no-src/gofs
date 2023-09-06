package about

import (
	"github.com/no-src/gofs/internal/version"
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
func PrintAbout(out func(format string, args ...any)) {
	out(logo)
	out("The gofs is a real-time file synchronization tool out of the box based on Golang")
	out("Open source repository at: <%s>", openSourceUrl)
	out("Download the latest version at: <%s>", releaseUrl)
	out("The docker image repository address at: <%s>", dockerImageUrl)
	out("Full documentation at: <%s>", documentationUrl)
}
