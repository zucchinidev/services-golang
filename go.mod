module github.com/zucchini/services-golang

go 1.24.3

require (
	github.com/ardanlabs/conf/v3 v3.8.0
	github.com/arl/statsviz v0.6.0
	github.com/go-json-experiment/json v0.0.0-20250517221953-25912455fbc8
	github.com/google/uuid v1.6.0
)

require (
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	github.com/antonholmquist/jason v1.0.0 // indirect
	github.com/bsiegert/ranges v0.0.0-20111221115336-19303dc7aa63 // indirect
	github.com/divan/expvarmon v0.0.0-20230430154648-8e0b3d2778b3 // indirect
	github.com/gizak/termui v0.0.0-20181228210747-b136f68f55f1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/nsf/termbox-go v0.0.0-20180613055208-5c94acc5e6eb // indirect
	github.com/pyk/byten v0.0.0-20140925233358-f847a130bf6d // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/telemetry v0.0.0-20240522233618-39ace7a40ae7 // indirect
	golang.org/x/tools v0.30.0 // indirect
	golang.org/x/vuln v1.1.4 // indirect
	honnef.co/go/tools v0.6.1 // indirect
)

tool (
	github.com/divan/expvarmon
	golang.org/x/vuln/cmd/govulncheck
	honnef.co/go/tools/cmd/staticcheck
)
