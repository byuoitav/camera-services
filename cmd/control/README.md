# Control Service
The control service provides the web app for controlling the cameras

## Environment Variables
```
GIN_MODE=debug
PORT=8080
LOG_LEVEL="debug"
DB_ADDRESS=couch_database_address
DB_USERNAME=couch_user
DB_PASSWORD=couch_password
KEY_SERVICE=control_keys_address
CALLBACK_URL=https://cameras-url.av.byu.edu
CLIENT_ID=id_for_tyk
CLIENT_SECRET=secret_for_tyk
GATEWAY_URL=api_gateway_address
OPA_URL=opa_address
OPA_TOKEN=opa_token
SIGNING_SECRET=session_signing_secret
AVER_PROXY=address_for_aver
AXIS_PROXY=address_for_axis
```
## Flags
| Flag               | Shorthand | Default                               | Description                                                                        |
|--------------------|-----------|---------------------------------------|------------------------------------------------------------------------------------|
| `--port`           | `-P`      | `8080`                                | Port to run the server on.                                                         |
| `--log-level`      | `-L`      | `""` (empty)                          | Level to log at. [Refer to zapcore.Level options](https://godoc.org/go.uber.org/zap/zapcore#Level). |
| `--db-address`     |           | `""`                                  | Database address.                                                                  |
| `--db-username`    |           | `""`                                  | Database username.                                                                 |
| `--db-password`    |           | `""`                                  | Database password.                                                                 |
| `--db-insecure`    |           | `false`                               | Don't use SSL in the database connection.                                          |
| `--key-service`    |           | `control-keys.av.byu.edu`            | Address of the control keys service.                                               |
| `--callback-url`   |           | `http://localhost:8080`               | WSO2 callback URL.                                                                  |
| `--client-id`      |           | `""`                                  | WSO2 client ID.                                                                     |
| `--client-secret`  |           | `""`                                  | WSO2 client secret.                                                                 |
| `--gateway-url`    |           | `https://api.byu.edu`                 | WSO2 gateway URL.                                                                   |
| `--opa-url`        |           | `""`                                  | URL of the OPA Authorization server.                                               |
| `--opa-token`      |           | `""`                                  | Token to use for OPA.                                                              |
| `--disable-auth`   |           | `false`                               | Disable all authentication checks.                                                 |
| `--signing-secret` |           | `""`                                  | Secret to sign JWT tokens with.                                                    |
| `--aver-proxy`     |           | `""`                                  | Base URL to proxy camera control requests through.                                 |
| `--axis-proxy`     |           | `""`                                  | Base URL to proxy camera control requests through.                                 |

## Endpoints 
Get Control Information
* <mark>GET</mark> `/api/v1/controlInfo`
* Returns the control information for the cameras

```
GET
	https://cameras-address.byu.edu/api/v1/controlInfo?key=114768

Response

    {"room":"JET-1234","controlGroup":"JET 1234","controlKey":""}

```

Get Cameras
* <mark>GET</mark> `/api/v1/cameras`
* Returns the cameras for the control group
```
GET
    https://cameras-address.byu.edu/api/v1/cameras?room=JET-1234&controlGroup=ITB%201106&controlKey=114768

Response: 

    [{"displayName":"Camera","tiltUp":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/pantilt/up","tiltDown":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/pantilt/down","panLeft":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/pantilt/left","panRight":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/pantilt/right","panTiltStop":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/pantilt/stop","zoomIn":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/zoom/in","zoomOut":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/zoom/out","zoomStop":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/zoom/stop","stream":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/stream","presets":[{"displayName":"Room","savePreset":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/savePreset/0","setPreset":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/preset/0"}],"reboot":"https://cameras-address.byu.edu/proxy/aver/v1/Pro520/JET-1234-CAM1.byu.edu:12345/reboot"}]
	
```
Camera Stream Proxies
* <mark>GET</mark> `/api/v1/proxy/aver/*uri`
* <mark>GET</mark> `/api/v1/proxy/axis/*uri`
