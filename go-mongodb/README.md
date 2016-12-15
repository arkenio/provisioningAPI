# go-mongodb


# Quick start
```go
import "github.com/heimweh/go-mongodb/atlas"
client := atlas.NewClient(username, apiKey)
```

## Listing Group Whitelists
```go
list, _, _ := client.GroupWhiteList.List(groupID)
for _, wl := range list {
  log.Printf("GroupID: %s", wl.GroupID)
  log.Printf("CidrBlock: %s", wl.CidrBlock)
  log.Printf("IPAddress: %s", wl.IPAddress)
}
```

## Creating a Group Whitelist
```go
whiteList := atlas.GroupWhiteList{
  CidrBlock: "0.0.0.0/0",
  GroupID: groupID
}

list, _, _ := client.GroupWhiteList.Create(whiteList)
```

## Retrieving a Group Whitelist
```go
list, _, _ := client.GroupWhiteList.Get(groupID, "<cidr|ipaddress>")
```

## Deleting a Group Whitelist
```go
resp, err := client.GroupWhiteList.Delete(groupID, "<cidr|ipaddress>")
```
