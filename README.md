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

Pass multiple labels in one command to avoid paying the printer's leader feed for each label. Each page is fully cut.

```bash
./bin/print-label openclaw sfactory
```

The command reads `~/.config/label/config.yaml` relative to the process's `HOME` for compatibility with the previous script. CLI flags override config values.

Raw TCP confirms only that the print data was written to the printer connection. It does not provide reliable print-completion status, so verify the printer's status LED if a label does not appear.

## Release

Push a version tag to build and publish GitHub release binaries:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The release workflow uses GoReleaser and attaches tarballs named for mise/ubi-style GitHub release installs.
