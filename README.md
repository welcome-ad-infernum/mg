# mg

## How to use:

### With defaults:

 * IOS `mg`
 * Windows `mg.exe`

### Available flags:

 * `-n`: int, number of requests per each target (default 1000000)
 * `-s`: string, url to endpoint or file name (default "ukraine.txt")
 * `-t`: string, source type to use (file or endpoint) (default "file")
 * `-w`: int, number of workers per logical CPU (default 10)

### Target model for endpoint to consume:

```json
{
	"id": "<random-uuid>",
	"url": "https://example.ru",
	"method": "GET",
	"data": null,
	"headers": [["<header-key>", "<header-value>"]],
	"proxy_url": null
}
```
