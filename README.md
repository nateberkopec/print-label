# print-label

Print labels on a Brother PT-P900W using its raster TCP protocol.

## Install and use with mise

```bash
mise use -g github:nateberkopec/print-label@latest
print-label Birdie
print-label "A very long label"
print-label two labels
print-label "Two very" "long labels"
```

`print-label --usage` outputs a [Usage](https://usage.jdx.dev/) KDL specification for shell completions and other tooling.

## Development

```bash
mise deps
mise run test
mise run lint
mise run build
```

### Releases

We publish GitHub release binaries and version tags. The release workflow uses GoReleaser and attaches tarballs named for mise/ubi-style GitHub release installs.
