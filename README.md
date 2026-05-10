# print-label

Print labels on a Brother PT-P900W using its raster TCP protocol.

## Install with mise

After a GitHub release exists:

```bash
mise use -g github:nateberkopec/print-label@latest
```

## Development

```bash
mise deps
mise run test
mise run lint
mise run build
```

## Preview without printing

```bash
mise run build
./bin/print-label --dry-run --out /tmp/birdie.png Birdie
```

## Print

```bash
./bin/print-label Birdie
```

The command reads `~/.config/label/config.yaml` for compatibility with the previous script. CLI flags override config values.

## Release

Push a version tag to build and publish GitHub release binaries:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow uses GoReleaser and attaches tarballs named for mise/ubi-style GitHub release installs.
