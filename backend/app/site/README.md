# Bundled UI

`make build` builds the frontend into this directory's ignored `dist/`
subdirectory before compiling the Go server. The server embeds those files
and serves them on the same port as `/api/`.

For development, `APP_SERVER_WEB_UI_PROXY=http://127.0.0.1:5179` sends UI
requests and hot-reload connections to Vite instead. API requests remain
in Go. This README keeps the embed directory present in a clean checkout;
starting without a UI build or development proxy reports a configuration error.
