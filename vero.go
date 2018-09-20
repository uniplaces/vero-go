package vero_go

type Client interface {
	Identify(userId string, data map[string]interface{}, email *string) ([]byte, error)
	Reidentify(userId string, newUserId string) ([]byte, error)
	Update(userId string, changes map[string]interface{}) ([]byte, error)
	Tags(userId string, add []string, remove []string) ([]byte, error)
	Unsubscribe(userId string) ([]byte, error)
	Resubscribe(userId string) ([]byte, error)
	Track(
		eventName string,
		identity map[string]string,
		data map[string]interface{},
		extras map[string]interface{},
	) (
		[]byte,
		error,
	)
}
