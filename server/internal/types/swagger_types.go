// Package types define the structure of request parameters and return results in this package
package types

// This file is public struct, only used to generate swagger documents, it is recommended
// to write comments to all of them, if you use postman or yapi, import swagger.json
// into postman or yapi and fill in the notes automatically to avoid repeating the comments.

// Result output data format
type Result struct {
	Code int         `json:"code"` // return code
	Msg  string      `json:"msg"`  // return information description
	Data interface{} `json:"data"` // return data
}

// Params query parameters
type Params struct {
	Page  int    `json:"page"`           // page number, starting from page 0
	Limit int    `json:"limit"`          // lines per page
	Sort  string `json:"sort,omitempty"` // sorted fields, multi-column sorting separated by commas

	Columns []Column `json:"columns,omitempty"` // query conditions
}

// Column information
type Column struct {
	Name  string      `json:"name"`  // column name
	Exp   string      `json:"exp"`   // expressions, which default to = when the value is null, have =, !=, >, >=, <, <=, like
	Value interface{} `json:"value"` // column value
	Logic string      `json:"logic"` // logical type, default value is "and", support &, and, ||, or
}

// Conditions query conditions
type Conditions struct {
	Columns []Column `json:"columns"` // columns info
}

type Options struct {
	Label    string      `json:"label"` // 标签
	Value    interface{} `json:"value"` // 值
	Children []Options   `json:"children"`
}

// OptionsReply only for api docs
type OptionsReply struct {
	Code int       `json:"code"` // return code
	Msg  string    `json:"msg"`  // return information description
	Data []Options `json:"data"` // return data
}
