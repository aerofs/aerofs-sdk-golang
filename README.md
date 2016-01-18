# aerofs-sdk-golang
An AeroFS Private Cloud API SDK written in Golang. The AeroFS Golang SDK is
composed of two packages: 
* **aerofsapi** -  Map the AeroFS API spec to individual calls
  * Supports all routes documented by the AeroFS API v1.3 Specification
* **aerofssdk** - Higher-level interface to the API
  * Supports the creation of File, Folder, Group, GroupMember, SharedFolderMember and User objects

### Installation
```sh
$ go get github.com/aerofs/aerofs-sdk-golang/aerofsapi
$ go get github.com/aerofs/aerofs-sdk-golang/aerofssdk
```

### Melkor
Melkor is a test app that uses the API,SDK to enumerate lists of files, folders
and number of users on an AeroFS deployment
