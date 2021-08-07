package dto

import "encoding/json"

type RPCRequest struct {
	Method  string          `json:"method,omitempty"`
	Id      string          `json:"id,omitempty"`
	Jsonrpc string          `json:"jsonrpc"`
	Params  json.RawMessage `json:"params"`
}
