package qdrant

type CreateBodyReq struct {
	Vectors struct {
		Size     uint64
		Distance string
	}
}

type UpsertRecipeReq struct {
	Points []PointQuery
}

type PointQuery struct {
	ID      any                 `json:"id"`
	Vector  any                 `json:"vector"`
	Payload UpsertRecipePayload `json:"payload"`
}

type UpsertRecipePayload struct {
	Title        string   `json:"title"`
	Ingredients  []string `json:"ingredients"`
	Diet         string   `json:"diet"`
	Cuisine      string   `json:"cuisine"`
	CookTime     int      `json:"cook_time"`
	Calories     int      `json:"calories"`
	Source       string   `json:"source"`
	Instructions string   `json:"instructions"`
	CreatedAt    string   `json:"created_at"`
}
