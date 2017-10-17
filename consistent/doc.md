# package consistent
```go
import "stathat.com/c/consistent"
```
Install: `go get stathat.com/c/consistent`

`consistent` is a Go package implementing a consistent hash function.

Consistent hashing is often used to distribute requests to a changing set of servers. For example, say you have some cache servers cacheA, cacheB, and cacheC. You want to decide which cache server to use to look up information on a user.

You could use a typical hash table and hash the user id to one of cacheA, cacheB, or cacheC. But with a typical hash table, if you add or remove a server, almost all keys will get remapped to different results, which basically could bring your service to a grinding halt while the caches get rebuilt.

With a consistent hash, adding or removing a server drastically reduces the number of keys that get remapped.

Read more about consistent hashing on [wikipedia](http://en.wikipedia.org/wiki/Consistent_hashing).

## Usage
You can use `consistent` for a lot of different purposes, but for our example, we'll continue with the task of finding a cache server for a user.

Start off by creating a `Consistent` object:
```go
cons := consistent.New()
```
Add an initial set of servers:
```go
cons.Add("cacheA")
cons.Add("cacheB")
cons.Add("cacheC")
```
Find a cache server for a user:
```go
server, err := cons.Get("user_89138238")
if err != nil {
        return err
}
fmt.Println("server:", server)
```
If you find out that server cacheB is unavailable, remove it:
```go
cons.Remove("cacheB")
```
Subsequential `Get()` calls will efficiently map the user key to cacheA and cacheC.

If a new server comes online, add it:
```go
cons.Add("cacheX")
```