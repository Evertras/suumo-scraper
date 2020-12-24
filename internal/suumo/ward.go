package suumo

type Ward struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Prefecture Prefecture `json:"prefecture"`
}
