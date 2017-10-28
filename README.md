# Servers Survey

A small Go program to perform parallel queries on a number of web-sites and
grab the (advertised) server software version they may be running.  This was
inspired by the long-running Netcraft "What's that site running?" monthly
surveys. The list of sites is based on the largest public companies in the
U.S., based on Fortune 1000 rankings.

This was written as a first exercise in Go programming, to learn Go language
basics, running goroutines, using channels, and standard libraries for using
HTTP and parsing textfiles.

An earlier implementation written in Ruby is also included, for language
comparison purposes.

## Building and Running

```
$ go build servers-survey.go

$ ./servers-survey
```
