# print-label

Print labels on a Brother PT-P900W using its raster TCP protocol.

## Install and use with mise

```bash
mise use -g github:nateberkopec/print-label@latest
print-label Birdie
print-label "A very long label"
print-label two labels
print-label "Two very" "long labels"
print-label --image logo.png
```

## Images

`--image` treats every argument as a path to a PNG. Each image prints as is, with no scaling or dithering. Its height must match the tape's printable height exactly: 48px for 3.5mm, 64px for 6mm, 106px for 9mm, 150px for 12mm, 234px for 18mm, 320px for 24mm, and 454px for 36mm. Its width sets the label length. Only fully opaque pure black pixels print; everything else is left blank.

Network printing sends raster data directly to TCP port 9100. It does not wait for a raster status reply: the PT-P900W does not return one on this connection. Check that the loaded tape matches `--tape-width`; successful submission does not confirm physical printing. See Brother's [Raster Command Reference, section 5.6](https://download.brother.com/welcome/docp100407/cv_ptp900_eng_raster_102.pdf) for the network printing flow.

## HTTP server

`print-label --listen :8080` serves print requests instead of printing its arguments. Other flags and the config file set the defaults.

```bash
curl -X POST localhost:8080/print -d '{"labels": ["Birdie", "Sky"], "tape_width": 12}'
```

Send images as base64 PNGs in `images` instead of `labels`:

```bash
curl -X POST localhost:8080/print -d "{\"images\": [\"$(base64 < logo.png | tr -d '\n')\"]}"
```

Send exactly one of `labels` or `images`; `tape_width` is optional. Images follow the same rules as `--image`. Request bodies are limited to 8 MiB. A 200 response means the job was sent, not that it printed. Invalid input returns 400 and printer connection failures return 502, each with an `error` field. `GET /health` returns 204. There is no authentication: listen only on a trusted network.

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
