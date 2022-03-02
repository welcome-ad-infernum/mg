# mg

## How to use:

### With defaults:

 * MacOS `mg`
 * Windows `mg.exe`

### Available flags:

 * `-n`: int, number of requests per each target (default 1000000)
 * `-s`: string, url to endpoint or file name (default "ukraine.txt")
 * `-t`: string, source type to use (file or endpoint) (default "file")
 * `-w`: int, number of workers per logical CPU (default 10)

### Target model for endpoint to consume:

```json
{
	"id": 0,
	"url": "https://example.ru",
	"method": "GET",
	"data": null,
	"headers": [["<header-key>", "<header-value>"]],
	"proxy_url": null
}
```

### For non-technical users:

 1. Go to [releases page](https://github.com/welcome-ad-infernum/mg/releases).
 2. Find the latest release.
 3. Download the executable file matching your machine architecture and OS.
 4. Open your terminal (`Terminal` for MacOS, `CMD` or `PowerShell` for Windows)
 6. Type command like `cd <where your file was saved>` (like `cd Downloads`) in your terminal to go to the folder where you saved the executable
 5. Run `./mg` or `.\mg.exe` depending on your platform with flags described above
