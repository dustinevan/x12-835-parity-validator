# X12 835 Parity Validator

Compares the 835 remittances a candidate claims system produces against the
ones the legacy Medicare systems produced for the same claims, claim by
claim, element by element, and shows the counts, the dollars, and every
difference so subject-matter experts can judge how close the candidate is.
A web app with a Go backend and a PostgreSQL database; the inputs are X12
835 files.

## Set up

You need Go 1.25, Node 22, and PostgreSQL running locally. The defaults
expect Postgres.app on `localhost:5432` with trust authentication (user
`postgres`); set `PG_BIN` if `psql` is elsewhere. Then:

```sh
make db
```

That creates the `validator` database, applies the schema, generates a
synthetic data set (`make gen`: 20,000 legacy claims and two vendors' 835s,
three submission dates each, all invented), loads it, compares it, and
summarizes it. It takes under a minute.

## Start

Three terminals, one process each; Ctrl+C stops it.

```sh
make app     # the Go server on http://localhost:8089, forwarding the UI to Vite
make ui      # the app's UI with hot reload (Vite on 5179)
make docs    # the docs site on http://localhost:5180: the 835 viewer, the compare page, the reference
```

Open http://localhost:8089. `make demo` builds the UI into the binary and
runs it alone on 8089, the shape that ships. `make list` shows every
command; `make check` is the gate (gofmt, vet, tests, the frontend build).

## Loading real data

`bin/vdb` is the batch tool: `vdb migrate`, `vdb load` (a directory or an S3
prefix of 835 files), `vdb compare`, `vdb summarize`. `vdb help` explains
each, and `vdb settings` lists the environment variables the service and the
tool share (database, S3). Only files loaded, compared, and summarized show
in the app.

## License

CC0 1.0 Universal, a public domain dedication: see `LICENSE`. Copy, change,
and use any of it without attribution. The third-party Go modules under
`vendor/` are not covered by it; each keeps its authors' license, included
beside its code.
