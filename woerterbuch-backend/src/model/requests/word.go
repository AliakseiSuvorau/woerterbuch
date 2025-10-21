package requests

type AddWordRequest struct {
	Article     string `json:"article" validate:"Required"`
	Word        string `json:"word" validate:"Required"`
	Translation string `json:"translation"`
}

type EditWordRequest struct {
	ID          uint64 `json:"id" validate:"Required"`
	Article     string `json:"article"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
}

type DeleteWordRequest struct {
	ID uint64 `json:"id" validate:"Required"`
}
