# DLS server

DLS server written in Go.

### Environment variables

| Key                     | Description                                                                                                                       | Required |
|-------------------------|-----------------------------------------------------------------------------------------------------------------------------------|----------|
| `CONTAINER_PORT`        | The port for the container to listen for connections on, defaults to `HTTP_PORT`.                                                 | No       |
| `DEBUG`                 | Enable debug for requests and responses. **Note:** this is *very* verbose.                                                        | No       |
| `FORWARDER_HOST`        | The external hostname to reach the server, defaults to `localhost`. **Note:** this should be the hostname of the proxy/forwarder. | No       |
| `FORWARDER_PORT`        | The external port to reach the server, defaults to `80`. **Note:** this should be the port of the proxy/forwarder.                | No       |
| `INSTANCE_REFERENCE`    | UUID for the instance, defaults to `00000000-0000-0000-0000-000000000000`.                                                        | No       |
| `KEY_REFERENCE`         | UUID for the key, defaults to `10000000-0000-0000-0000-000000000001`.                                                             | No       |
| `LEASE_DURATION`        | Duration of lease in format 1w2d3h4m5ns. Max 90d, defaults to `90d`.                                                              | No       |
| `LEASE_RENEWAL_PERCENT` | Percentage of licence time before a renewal should be triggered. Max 100(%), defaults to `15`(%).                                 | No       |
| `NOTIFICATION_URLS`     | URLs to notify via [Shoutrrr](https://github.com/containrrr/shoutrrr).                                                            | No       |
| `SCOPE_REFERENCE`       | UUID for the scope, defaults to `20000000-0000-0000-0000-000000000002`.                                                           | No       |
| `SIGNING_KEY`           | Text for the private key used to sign JWTs.                                                                                       | Yes*     |
| `SIGNING_KEY_PATH`      | Path to the private key used to sign JWTs.                                                                                        | Yes*     |

<sub>* One of `SIGNING_KEY` or `SIGNING_KEY_PATH` must be specified. Providing a text signing key allows DLS to run completely stateless.</sub>

### Creating a signing key

Signing keys can be generated using openssl.

```sh
openssl genrsa -out signingkey.pem 2048 
```

### Serve via HTTPS

Connections must be served via TLS connections, this should be run behind a proxy or other forwarder which provides SSL termination.

### Download tok to instance

The `.tok` file can be downloaded using the following command.

```sh
curl -L -X GET https://<forwarder_host>:<forwarder_port>/token/download -o /etc/nvidia/ClientConfigToken/client_configuration_token_$(date '+%d-%m-%Y-%H-%M-%S').tok
```
