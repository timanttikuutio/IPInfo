A simple GO script used for viewing IP-address information from the command line.

**Installation instructions**

First git clone the repository to your computer and make sure you have GO installed.

You'll then have to install a few go libraries for building the project.

Run the following commands in the newly downloaded project directory:
 - `go get github.com/joho/godotenv`
 - `go get github.com/TwiN/go-color`

After that's done, you can run `go build -o IPInfo main.go` or `GOOS=windows go build -o IPInfo.exe main.go`, if you're on Windows.

That's it!

--

You'll also be able to use prebuilt binaries built by me, found on the releases page.
