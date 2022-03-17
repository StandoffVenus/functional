# Functional

Functional provides a small set of pure functions that are common in functional programming languages, such as `Reduce`, `Map`, `Filter`, etc.

With the release of Go [1.18](https://go.dev/blog/go1.18), these functions are nothing more than a few lines of code. Prior, there was no practical way to do type-safe functional programming in Go.

# Testing

Since `functional` is a small package, there's not much to test. You can test by hand by running

```
go test ./...
```

Or use Make:

```
make
```