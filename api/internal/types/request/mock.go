package request

type MockRequest struct {
	Field1 string `json:"field1" binding:"required"`
	Field2 string `json:"field2" binding:"required,eqfield=Field1"`
}