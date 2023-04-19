package category

// ShortCategory is a structure for create new category
type ShortCategory struct {
	Name        string `json:"name" binding:"required" example:"Электротехника"`
	Description string `json:"description" binding:"required" example:"Электротехнические товары для дома"`
	Image       string `json:"image,omitempty"`
}

// CategoryId is a structure for displaying the result of creating a category
type CategoryId struct {
	Value string `json:"id" uri:"categoryID" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

// Category is a structure for updating a category and displaying results containing a list of categories
type Category struct {
	Id          string `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Name        string `json:"name" binding:"required" example:"Электротехника"`
	Description string `json:"description" binding:"required" example:"Электротехнические товары для дома"`
	Image       string `json:"image,omitempty"`
}
