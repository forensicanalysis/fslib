<h1 align="center">{{.Name}}</h1>

<p  align="center">
 <a href="https://{{.ModulePath}}/actions"><img src="https://{{.ModulePath}}/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/{{.RelModulePath}}"><img src="https://codecov.io/gh/{{.RelModulePath}}/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://goreportcard.com/report/{{.ModulePath}}"><img src="https://goreportcard.com/badge/{{.ModulePath}}" alt="report" /></a>
 <a href="https://pkg.go.dev/{{.ModulePath}}"><img src="https://img.shields.io/badge/go.dev-documentation-007d9c?logo=go&logoColor=white" alt="doc" /></a>
 <a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib.svg?type=shield"/></a>
</p>

{{.Doc}}

## Installation

``` shell
go get -u {{.ModulePath}}
```


{{if .Examples}}
## Examples
{{ range $key, $value := .Examples }}
{{if $key}}### {{ $key }}{{end}}
``` go
{{ $value }}
```
{{end}}{{end}}

{{if .Bugs}}
## Known Bugs
{{range .Bugs}}* {{.}}{{end}}
{{end}}



## fs command

The fs command line tool that has various subcommands which imitate unix commands
but for nested file system structures.

 - **fs cat**: Print files
 - **fs file**: Determine files types
 - **fs hashsum**: Print hashsums
 - **fs ls**: List directory contents
 - **fs stat**: Display file status
 - **fs strings**: Find the printable strings in an object, or other binary, file
 - **fs tree**: List contents of directories in a tree-like format


#### Download

https://github.com/forensicanalysis/fslib/releases

#### Usage Examples

List all files in a zip file:
```
fs ls test.zip
```

Extract the Amcache.hve file from a NTFS image in a zip file:

```
fs cat case/evidence.zip/ntfs.dd/Windows/AppCompat/Programs/Amcache.hve > Amcache.hve
```

Hash all files in a zip file:
```
fs hashsum case/evidence.zip/*
```


{{if .Subpackages}}
## Subpackages

| Package | Description |
| --- | --- |
{{ range $key, $value := .Subpackages }}| **{{ $value.Name }}** | {{ $value.Synopsis }} |
{{end}}{{end}}

## Contact

For feedback, questions and discussions you can use the [Open Source DFIR Slack](https://github.com/open-source-dfir/slack).

## Acknowledgment

The development of this software was partially sponsored by Siemens CERT, but
is not an official Siemens product.
