package models

type (
	AstNode struct {
		ID       int      `json:"id"`
		AstType  string   `json:"type"`
		Value    string   `json:"operation"`
		Left     *AstNode `json:"arg1"`
		Right    *AstNode `json:"arg2"`
		Counting bool     `json:"status"`
	}

	Expression struct {
		ID         int     `json:"id"`
		Expression string  `json:"expression"`
		Status     string  `json:"status"`
		Result     float64 `json:"result"`
	}

	User struct {
		ID       int64
		Login    string
		Password string
	}

	Result struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
		Error  string  `json:"error"`
	}
)
