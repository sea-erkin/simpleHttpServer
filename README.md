# simpleHttpServer

Serves files at the current, or whichever directory you provide using the -d flag. 
Also has support for TLS if the certificate chain .pem and key.pem are provided.

```
  Usage of ./simpleHttpServer:
  -c string
        (optional) -c Path to cert chain
  -d string
        (optional) -d Path to directory to serve
  -k string
        (optional) -k Path to cert private key
  -p string
        -p Port to listen on - will set to 80 if not provided
``` 
