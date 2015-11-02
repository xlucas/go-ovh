# go-ovh
A simple helper library around the OVH API for golang developers.

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/xlucas/go-ovh/ovh)


## Requirements
Firsteval, you will need to generate the application credentials in order to use the API. Check [the official guide here](https://api.ovh.com/g934.test) for more details.

## Usage examples

An example using simple string hashmaps as json in/out structures.
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

if jsonIn, err = c.CallSimple("POST", "/cloud/createProject", jsonOut); err != nil {
    log.Fatal("Failed to call the API, reason : ", err)
}

// Access the json we got back
id := jsonIn["project"]
...
```


Another example using struct binding.

```go

import "github.com/xlucas/go-ovh/ovh"

typedef MyRequestStruct struct {
    MyField1    string
    MyField2    string
}

typedef MyResponseStruct struct {
    MyIdField   uint
    MyField1    string
}

c := ovh.Client(ovh.ENDPOINT_EU_OVHCOM, "MyAppKey", "MyAppSecretKey", "MyConsumerKey")

// Check for time lag
if err = c.PollTimeshift(); err != nil {
    log.Fatal("Failed to retrieve timeshift, reason : ", err)
}

// Send our request
out = MyRequestStruct {
    MyField1: "foo",
    MyField2: "bar",
}

var in MyResponseStruct

if in, err = c.Call("POST", "/cloud/createProject", out); err != nil {
    log.Fatal("Failed to call the API, reason : ", err)
}

// Access the response struct
fmt.Printf("Object id is %d", in.MyIdField)
...
```
