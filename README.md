[![](https://img.shields.io/docker/pulls/itzg/web-debug-server.svg)](https://hub.docker.com/r/itzg/web-debug-server)

# web-debug-server

A very minimal web server that responds with a page containing the request headers and content

## Usage

```
  -bind host:port
    	The host:port to bind, but using port flag is preferred (env BIND)
  -port int
    	The port to bind (env PORT) (default 8080)
  -response-fixed-body string
    	When set, specifies a fixed body to write to the response (env RESPONSE_FIXED_BODY)
  -response-fixed-content-type string
    	When FixedBody is set, specifies the content type to set (env RESPONSE_FIXED_CONTENT_TYPE) (default "text/plain")
  -response-status int
    	When set, specifies the status code to use in responses (env RESPONSE_STATUS) (default 200)
```