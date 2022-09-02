package rabbit

import (
	"github.com/streadway/amqp"

	"github.com/h14yhv/golang-lib/clock"
)

const (
	Module            = "RABBIT"
	ScheduleReconnect = 2 * clock.Second
	SchedulePublish   = 3 * clock.Second
	ScheduleConsume   = 3 * clock.Second

	MIMEApplicationJSON = "application/json"
	MIMETextPlain       = "text/plain"

	Transient  = amqp.Transient
	Persistent = amqp.Persistent
)
