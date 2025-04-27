package notifier

import (
	"strings"
	"text/template"
	"time"

	"go.uber.org/zap"

	"github.com/GoSeoTaxi/tg_echo/internal/telegram"
)

type Notifier struct {
	sender *telegram.Sender
	log    *zap.Logger
}

type Message struct {
	Body string
	Time time.Time
	IP   string
}

var tpl = template.Must(template.New("msg").Parse(`
++++++++++++++++++++
{{.Body}}
{{ .Time.UTC.Format "2006-01-02T15:04:05Z07:00" }}
{{- if .IP }}
ip:{{.IP}}
{{- end }}
++++++++++++++++++++`))

func New(token string, chatID int64, log *zap.Logger) (*Notifier, error) {
	s, err := telegram.New(token, chatID)
	if err != nil {
		return nil, err
	}
	return &Notifier{sender: s, log: log}, nil
}

func (n *Notifier) Send(m Message) error {
	var b strings.Builder
	if err := tpl.Execute(&b, m); err != nil {
		n.log.Error("tpl", zap.Error(err))
		return err
	}
	if err := n.sender.Send(b.String()); err != nil {
		n.log.Error("telegram", zap.Error(err))
		return err
	}
	n.log.Info("sent", zap.String("body", m.Body))
	return nil
}
