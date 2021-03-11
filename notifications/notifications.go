package notifications

import (
	"database/sql"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Notification struct {
	types.Notification
	TriggeredByBlock int64
}

type Notifications struct {
	stmt   *sql.Stmt
	notifs []Notification
}

func (n *Notifications) Append(notification ...Notification) {
	n.notifs = append(n.notifs, notification...)
}

func (n *Notifications) Exec() error {
	for _, notif := range n.notifs {
		_, err := n.stmt.Exec(notif.Target, notif.NotificationType, notif.TriggeredByBlock, notif.StartsOn, notif.ExpiresOn, notif.Message, notif.Metadata)
		if err != nil {
			return errors.Wrap(err, "could not exec statement")
		}
	}

	err := n.stmt.Close()
	if err != nil {
		return errors.Wrap(err, "could not close exec statement")
	}

	return nil
}

func NewWithTx(tx *sql.Tx) (*Notifications, error) {
	var n Notifications

	stmt, err := tx.Prepare(pq.CopyIn("notifications", "target", "type", "triggered_by_block", "starts_on", "expires_on", "message", "metadata"))
	if err != nil {
		return nil, errors.Wrap(err, "prepare notifications statement")
	}

	n.stmt = stmt

	return &n, nil
}

func NewNotification(target string, typ string, block int64, starts int64, expires int64, msg string, metadata map[string]interface{}) Notification {
	return Notification{
		Notification: types.Notification{
			Target:           target,
			NotificationType: typ,
			StartsOn:         starts,
			ExpiresOn:        expires,
			Message:          msg,
			Metadata:         metadata,
		},
		TriggeredByBlock: block,
	}
}
