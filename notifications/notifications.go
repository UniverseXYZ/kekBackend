package notifications

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "notifs")

type Notifications struct {
	db *sql.DB
}

func (n *Notifications) Run(ctx context.Context) {
	// poll for new jobs
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				log.Info("loop")
				// get jobs
				// execute jobs
			case <-ctx.Done():
				log.Info("received exit signal, stopping")
			}
		}
	}()
}

func New(config Config) (*Notifications, error) {
	log.Info("connecting to postgres")
	db, err := sql.Open("postgres", config.PostgresConnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "open postgres connection")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "ping postgres connection")
	}

	return &Notifications{
		db: db,
	}, nil
}

//
// type Notifications struct {
// 	stmt   *sql.Stmt
// 	notifs []Notification
// }
//
// func (n *Notifications) Append(notification ...Notification) {
// 	n.notifs = append(n.notifs, notification...)
// }
//
// func (n *Notifications) Exec() error {
// 	for _, notif := range n.notifs {
// 		_, err := n.stmt.Exec(notif.Target, notif.NotificationType, notif.IncludedInBlock, notif.StartsOn, notif.ExpiresOn, notif.Message, notif.Metadata)
// 		if err != nil {
// 			return errors.Wrap(err, "could not exec statement")
// 		}
// 	}
//
// 	err := n.stmt.Close()
// 	if err != nil {
// 		return errors.Wrap(err, "could not close exec statement")
// 	}
//
// 	return nil
// }
//
// func NewWithTx(tx *sql.Tx) (*Notifications, error) {
// 	var n Notifications
//
// 	stmt, err := tx.Prepare(pq.CopyIn("notifications", "target", "type", "included_in_block", "starts_on", "expires_on", "message", "metadata"))
// 	if err != nil {
// 		return nil, errors.Wrap(err, "prepare notifications statement")
// 	}
//
// 	n.stmt = stmt
//
// 	return &n, nil
// }
//
