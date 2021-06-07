package server

type DescriptionBody struct {
	Description string `json:",omitempty"`
}

type BadRequestResponse struct {
	Error string
	Msg   string
}
