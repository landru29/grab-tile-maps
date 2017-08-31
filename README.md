# grab-tile-maps

Grab all tiles of a map

## Build
```
    go get
    go build
    ./grab-tile-maps --help
```

## Usage

``` bash
  ./grab-tile-maps [flags]
```

### Flags

``` bash
      --out-folder string   output folder (default "./map")
      --out-scheme string   Map tile out files (default "{z}/{x}/{y}.png")
      --url-scheme string   Map tile urls (default "https://b.tile.opentopomap.org/{z}/{x}/{y}.png")
      --zoom-max int        Maximum zoom (default 3)
      --zoom-min int        Minimum zoom
```