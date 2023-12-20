# DLS server

DLS server written in Go.

### Environment variables

| Key                  | Description                                                                                                                       | Required |
|----------------------|-----------------------------------------------------------------------------------------------------------------------------------|----------|
| `DATABASE_DBNAME`    | The database to use, defaults to `dls`.                                                                                           | No       |
| `DATABASE_HOST`      | The hostname to connect to the database, defaults to `localhost`.                                                                 | No       |
| `DATABASE_PASSWORD`  | The password used to connect to the database.                                                                                     | No       |
| `DATABASE_USERNAME`  | The username used to connect to the database, defaults to the current user.                                                       | No       |
| `DEBUG`              | Enable debug for requests and responses. **Note:** this is *very* verbose.                                                        | No       |
| `HTTP_HOST`          | The hostname to listen for connections on, defaults to `localhost`. **Note:** this should be the hostname of the proxy/forwarder. | No       |
| `HTTP_PORT`          | The port to listen for connections on, defaults to `80`. **Note:** this should be the port of the proxy/forwarder.                | No       |
| `INSTANCE_REFERENCE` | UUID for the instance, defaults to `00000000-0000-0000-0000-000000000000`.                                                        | No       |
| `KEY_REFERENCE`      | UUID for the key, defaults to `10000000-0000-0000-0000-000000000001`.                                                             | No       |
| `NOTIFICATION_URLS` | URLs to notify via [Shoutrrr](https://github.com/containrrr/shoutrrr).
| `SCOPE_REFERENCE`    | UUID for the scope, defaults to `20000000-0000-0000-0000-000000000002`.                                                           | No       |
| `SIGNING_KEY_PATH`   | Path to the private key used to sign JWTs.                                                                                        | Yes    |

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
