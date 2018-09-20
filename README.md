# vero-go

Basic Golang client for[Vero](https://www.getvero.com/)developed at Uniplaces

## Usage

#### Import
```go
import "github.com/uniplaces/vero-go/vero"
```

#### Setup
```go
client := vero.NewClient("YOUR_AUTH_TOKEN")
```

#### Basic usage
```go
// Identify
data := make(map[string]interface{})
data["First name"] = "Jeff"
data["Last name"] = "Kane"
 
email := "jeff@yourdomain.com"
 
client.Identify("1234567890", data, email)
 
// Unsubscribe
client.Usubscribe("1234567890", data, email)
 
// Tags
add := []string{"Blog reader"}
remove := []string{}
                
client.Tags("1234567890", add, remove)
```
