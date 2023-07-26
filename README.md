# snowflake-go

### A full-featured Snowflake ID generator for go

Can set nodeBits, sequenceBits, epoch, nodeId.

### Install

```shell
go get github.com/JYeop/snowflake-go
```

### Example
```go
func main() {
    config := utils.SnowflakeConfig{
        NodeId:       1, // default: 0
        NodeBits:     16,  // default: 10
        SequenceBits: 6,  // default: 12
        Epoch:        1688169600000, // default: 1288834974657
    }
    new, err := utils.NewSnowflake(config)
    if err != nil {
        fmt.Println(err)
        return
    }
    id, err := new.Generate()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(id)
    tm := new.GetTimeFromId(id)
    fmt.Println(tm)
    node := new.GetNodeFromId(id)
    fmt.Println(node)
    sequence := new.GetSequenceFromId(id)
    fmt.Println(sequence)
    fmt.Println(new.Epoch)
}
```
