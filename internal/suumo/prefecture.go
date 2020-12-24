package suumo

import "fmt"

const (
	PrefectureCodeSaitama  = "11"
	PrefectureCodeChiba    = "12"
	PrefectureCodeTokyo    = "13"
	PrefectureCodeKanagawa = "14"
)

type Prefecture struct {
	Name    string `json:"name"`
	URLPath string `json:"urlPath"`
	Code    string `json:"code"`
}

func PrefectureFromCode(code string) Prefecture {
	switch code {
	case PrefectureCodeSaitama:
		return Prefecture{
			Name:    "埼玉県",
			URLPath: "saitama",
			Code:    "11",
		}

	case PrefectureCodeChiba:
		return Prefecture{
			Name:    "千葉県",
			URLPath: "chiba",
			Code:    "12",
		}

	case PrefectureCodeTokyo:
		return Prefecture{
			Name:    "東京都",
			URLPath: "tokyo",
			Code:    "13",
		}

	case PrefectureCodeKanagawa:
		return Prefecture{
			Name:    "神奈川県",
			URLPath: "kanagawa",
			Code:    "14",
		}

	default:
		panic(fmt.Sprintf("unknown code: %q", code))
	}
}
