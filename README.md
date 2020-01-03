# Superman Detector 

This was an interesting exercise! The US military actually uses similar restrictions (speed/height) to control the export of GPS receivers, classifying them as munitions if they continue to operate outside the parameters: https://en.wikipedia.org/wiki/Global_Positioning_System#Restrictions_on_civilian_use 

## Components/Structure

I used a fairly common golang repository structure here - `cmd/superman/` is where the main package is located, `data/` is where I store the two databases the program interacts with, and `pkg/` is where the majority of the code is stored. `bin/` is a directory created during the build task, and is not committed into the git repository. `test/` is scratchwork - raw cURLs live in `functional.sh` and `unit.sh` is currently a one-line `go test ./...` script

My goal with all my projects is to have all tooling operations occur within the Makefile. In this project, we have tasks for `go fmt`, `golangci-lint`, `docker build`, DB creation, binary creation (for local development/quick iteration), unit tests, and the initial setup. This currently only installs golangci-lint, but ideally it should install golang itself and any other prerequisites to development as well. 

My application uses two middlewares - `pkg/sanitize` to sanitize input and `pkg/geolocation` to grab the lat/lon coordinates - before hitting the SQLite DB in `pkg/haversine`

## External Libraries

- **Gin**: This is a golang microservice library that vastly simplifies the process of routing, using middleware chains, and preserving context up and down the stack
- **envconfig**: This library simplifies retrieving environment variables that hold configuration data, and using them in go programs by way of structs. 
- **geoip2-golang**: A client I found to simplify extraction of coordinates from an IP using the MaxMind DB format
- **go-sqlite3**: I used golang's inbuilt databases/sql package to operate on a sqlite database, but it doesn't come with a sqlite3 driver by default so I used this one. This does necessitate the use of CGO in my build process, but I thought it was a worthwhile tradeoff as other pure-go libraries don't have the same level of adoption/clout surrounding them. 
- **google/uuid**: This is a simple library I added to validate the UUID within the ID input field. 

## TODO

I ran into a couple of time constraints and personal commitments so my schedule for this challenge was shortened. Here's what my next steps look like: 

- A better logger. Right now I'm using fmt.Printf statements and while they're not terrible I'd still much rather be using a leveled logger. 
- Error handling. There are a lot of errors getting passed right back through the JSON from the err.Error(), which is a shortcut I take while developing but should be fixed before committing out of a feature branch 
- Better tests. Unit tests for the middleware are missing, and functional tests (beyond copy/pasting cURL calls) would be great to catch regressions in the code. 
- Better DB queries. The current code is copy/pasted between the prev and next row blocks, I believe they can be made into a single function. Also, by using the Prepare() function instead of Query() (which prepares a new statement every time) I should be able to achieve better performance at scale. 
- CI/CD. I've more or less ignored CI/CD during this challenge. I do have a personal gitlab server with an Auto Devops pipeline running, which takes a lot of the work off of my back, but I would have liked to get a proper gitlab-ci.yml or a .travis.yml file in the repo to run some simple CI on every commit. 


