# Spyglass
The Spyglass service provides the web app for gaining access to any camera control interface in the control service without the control key

## Environment Variables
```
GIN_MODE=debug
PORT=8080
LOG_LEVEL=info
DB_ADDRESS=couch_DB_address
DB_USERNAME=couch_user
DB_PASSWORD=couch_password
KEY_SERVICE=address_for_control-keys_service
CALLBACK_URL=https://spyglass-address.byu.edu
CLIENT_ID=tyk_id
CLIENT_SECRET=tyk_secret
GATEWAY_URL=api_gateway_address
OPA_URL=opa_address
OPA_TOKEN=opa_token
CONTROL_URL=address_for_control_service
```

## Flags
### Server Configuration Flags  

| Flag             | Shorthand | Default               | Description                                                                        |
|------------------|------------|-----------------------|------------------------------------------------------------------------------------|
| `--port`         | `-P`       | `8080`                | Port to run the server on.                                                          |
| `--log-level`    | `-L`       | `""` (empty)         | Level to log at. [Refer to zapcore.Level options](https://godoc.org/go.uber.org/zap/zapcore#Level). |

### Database Configuration Flags  

| Flag              | Default   | Description                                                      |
|-------------------|------------|------------------------------------------------------------------|
| `--db-address`    | `""`       | Database address.                                                |
| `--db-username`   | `""`       | Database username.                                               |
| `--db-password`   | `""`       | Database password.                                               |
| `--db-insecure`   | `false`    | Don't use SSL in the database connection.                         |

### External Service Flags  

| Flag                | Default                            | Description                                                      |
|---------------------|------------------------------------|------------------------------------------------------------------|
| `--key-service`     | `control-keys.av.byu.edu`          | Address of the control keys service.                              |
| `--callback-url`    | `http://localhost:8080`            | WSO2 callback URL.                                               |
| `--client-id`       | `""`                               | WSO2 client ID.                                                  |
| `--client-secret`   | `""`                               | WSO2 client secret.                                              |
| `--gateway-url`     | `https://api.byu.edu`              | WS02 gateway URL.                                                |
| `--opa-url`         | `""`                               | URL of the OPA Authorization server.                              |
| `--opa-token`       | `""`                               | Token to use for OPA.                                             |
| `--disable-auth`    | `false`                            | Disable all authorization and authentication checks.              |
| `--control-url`     | `https://cameras.av.byu.edu/key-login?key=%s` | URL format string of the camera control service.                  |


## Endpoints 
Get Rooms
* <mark>GET</mark> `/api/v1/rooms`
* Returns the rooms available for control

```
GET
	https://spyglass-address.byu.edu/api/v1/rooms

Response

    ["JET-1234","JET-1235","JET-1236","JET-1237","JET-1238"...]

```

Get Rooms
* <mark>GET</mark> `/api/v1/rooms/:room/controlGroups`
* Returns the control groups for the selected room
```
GET
    https://spyglass-address.avdev.byu.edu/api/v1/rooms/JET-1234/controlGroups

Response: 

   ["JET 1234"]
	
```
Camera Stream Proxies
* <mark>GET</mark> `/api/v1/rooms/:room/controlGroups/:controlGroup`
* Redirects to the corresponding page in the control service