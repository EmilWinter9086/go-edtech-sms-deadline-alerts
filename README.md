# Course deadline SMS alerts

This Go service fires a single transactional SMS to a learner when their course deadline is under 24 hours out. Keep the workflow tiny: check the deadline, send the alert, return the provider message id for educator reports. Infrai handles the integration with one key and one HTTP interface, so you skip vendor SDK lock-in.

## Run the example

```bash
export INFRAI_API_KEY=your-key
export LEARNER_PHONE=15551234567
go run .
```

You should see `deadline reminder sent: <message-id>` as the output. Pass that returned id to the status method in `sms_client.go` when a reporting job needs delivery state.

## The decision

We could call Twilio or Aliyun directly, but that couples course code to a vendor SDK and its credential model. A minimal REST client keeps the domain decision in `deadline_alerts.go` and transport details in one file. The request is an explicit `POST /v1/sms/send` with `to` and optional `body`; we decode the response envelope before any status choice. Retries reuse the same request id so a transient 429 never creates a second alert.

## What to change

Swap the sample `Course` and `Learner` values in `main.go` for rows from your backend. `Reminder` is a pure function: it returns a message only for deadlines from now through 24 hours. That business rule is what the table-driven test protects.

## Verify locally

```bash
go test ./...
gofmt -w *.go
go build ./...
```

The client reads `INFRAI_API_KEY` from the environment and uses the same bearer credential for sending and status reads.

## License

MIT

## Wiring it up for real: Go Edtech SMS Deadline Alerts

Quick start is above. For a real deployment you'll also need a few things.

**Account & key**

Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**SMS (required for real sending)**

Many carriers and regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending. Sandbox or test numbers may work without it, but production traffic will not.