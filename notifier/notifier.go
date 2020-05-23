package notifier

type Notifier interface {
	Notify(url string) error
}
