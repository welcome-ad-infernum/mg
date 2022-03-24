# mg

## Supported platforms

 * Windows binary
 * macOS binary
 * Linux binary and systemd script
 * Docker container
 * Kubernetes helm chart
 * AWS and Azure terraform scripts

## How to use

### Desktop binaries

1. Go to [releases page](https://github.com/welcome-ad-infernum/mg/releases).
2. Find the latest release.
3. Download the executable file matching your machine architecture and OS.
4. Open your terminal (`Terminal` for MacOS, `CMD` or `PowerShell` for Windows)
5. Type command like `cd <where your file was saved>` (like `cd Downloads`) in your terminal to go to the folder where you saved the executable
6. Run `./mg` or `.\mg.exe` depending on your platform with flags described above

### Linux systemd script

* Open the terminal on your Linux PC or server.
* Install it with:

`$ curl -s https://raw.githubusercontent.com/welcome-ad-infernum/mg/main/examples/linux/install.sh | sudo sh -`
* Uninstall it with:

`$ curl -s https://raw.githubusercontent.com/welcome-ad-infernum/mg/main/examples/linux/uninstall.sh | sudo sh -`

### Docker image

* You can run it with

`$ docker run --rm -it --name mg-agent vladstarr/mg-agent:latest`
* Alternatively, you can use `examples/docker-compose.yaml` file from our repo and configure for your needs. 

## For advanced users

### Terraform and Kubernetes deployment

Please navigate to `examples` folder and view respectable README file.

### Available flags for binaries

* `-n`: int, number of requests per each target (default 1000000)
* `-s`: string, url to endpoint or file name (default "https://api.itemstolist.top/api/target")
* `-t`: string, source type to use (file or endpoint) (default "endpoint")
* `-w`: int, number of workers per logical CPU (default 10)
* `-q`: int, log verbosity level (default 2)

### Target model for endpoint to consume

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