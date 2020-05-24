package notifier

type Notifier interface {
	Notify(webhookURL string) error
}
