package notifier

type Notifier interface {
	Notify(webhookURL, message string) error
}
