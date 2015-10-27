# go-ovh
A simple helper library around the OVH API for golang developers.

## Requirements
Firsteval, you will need to generate the application credentials in order to use the API. Check [the official guide here](https://api.ovh.com/g934.test) for more details.

## Usage example

```go

import "github.com/xlucas/go-ovh/ovh"

var jsonIn interface{}
var err error

c := ovh.Client(ovh.ENDPOINT_EU_OVHCOM, "MyAppKey", "MyAppSecret", "MyConsumerKey")

// Check for time lag
if err = c.PollTimeshift(); err != nil {
    log.Fatal("Failed to retrieve timeshift, reason : ", err)
}

// Send our request
jsonOut = map[string]interface{}{
    "description": "My New Project",
    "voucher": "My Voucher",
}

if jsonIn, err = c.Call("POST", "/cloud/createProject", jsonOut); err != nil {
    log.Fatal("Failed to call the API, reason : ", err)
}

// Access the json we got back
id := jsonIn["project"]
...
```

