package models

type AstNode struct {
	ID       int      `json:"id"`
	AstType  string   `json:"type"`
	Value    string   `json:"operation"`
	Left     *AstNode `json:"arg1"`
	Right    *AstNode `json:"arg2"`
	Counting bool     `json:"status"`
}

type Result struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
	Error  string  `json:"error"`
}
