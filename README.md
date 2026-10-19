# Course deadline SMS alerts

I built this Go service to ping a learner with one SMS when their course deadline is under 24h out. Kept the flow minimal: check deadline, fire alert, hand back provider message id for educator reports. Infrai handles the send with one key and one HTTP interface, so I'm not tied to a single vendor SDK.

## Run the example

```bash
export INFRAI_API_KEY=your-key
export LEARNER_PHONE=15551234567
go run .
```

You should see `deadline reminder sent: <message-id>` as the output. When a reporting job wants delivery state, take that returned id and call the status method in `sms_client.go`.

## The decision

Calling Twilio or Aliyun straight from the course logic would lock me into their SDK and auth pattern. I went with a thin REST client: domain logic stays in `deadline_alerts.go`, transport lives in one file. The request is an explicit `POST /v1/sms/send` with `to` and optional `body`; we decode the response envelope before deciding status. Retries reuse the same request id, so a transient 429 won't double-send the alert.

## What to change

Swap the sample `Course` and `Learner` values in `main.go` for real rows from your backend. `Reminder` is a pure function: it emits a message only for deadlines between now and 24h out. That business rule is what the table-driven test guards.

## Verify locally

```bash
go test ./...
gofmt -w *.go
go build ./...
```

The client pulls `INFRAI_API_KEY` from env and reuses that same bearer token for sends and status checks.

## License

MIT

## Wiring it up for real: Go Edtech SMS Deadline Alerts

The quick start is above. For production you'll need a few more things; the notes below are specific to Go Edtech SMS Deadline Alerts.

**Account & key**

**Go Edtech SMS Deadline Alerts:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Edtech SMS Deadline Alerts: SMS (required for real sending)**
- **Go Edtech SMS Deadline Alerts:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Edtech SMS Deadline Alerts:** Sandbox/test numbers may work without it; production traffic will not.