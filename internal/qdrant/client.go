package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/0ex-d/recipe_combo_qdrant/internal/model"
)

type QdrantClient struct {
	baseURL    string
	collection string
	httpClient *http.Client
	timeout    time.Duration
}

func NewQdrantClient(baseURL, collection string) *QdrantClient {
	const defaultTimeout = 10 * time.Second
	return &QdrantClient{
		baseURL:    baseURL,
		collection: collection,
		httpClient: &http.Client{Timeout: defaultTimeout},
		timeout:    defaultTimeout,
	}
}

func (c *QdrantClient) CreateCollectionIfNotExist(ctx context.Context, vectorSize uint64) error {
	url := fmt.Sprintf("%s/collections/%s", c.baseURL, c.collection)
	reqCtx, cancel := c.requestContext(ctx)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	var createBody CreateBodyReq
	createBody.Vectors.Size = vectorSize
	createBody.Vectors.Distance = "Cosine"

	payload, err := json.Marshal(createBody)
	if err != nil {
		return err
	}
	req, err = http.NewRequestWithContext(reqCtx, http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		if resp.StatusCode == 401 {
			return fmt.Errorf("[CreateCollectionIfNotExist]:qdrant create collection auth error: body: %s", bodyStr)
		}
		return fmt.Errorf("[CreateCollectionIfNotExist]:qdrant create collection failed: body: %s", bodyStr)
	}
	return nil
}

func (c *QdrantClient) UpsertRecipe(ctx context.Context, recipe model.Recipe, vector []float32) error {
	url := fmt.Sprintf("%s/collections/%s/points?wait=true", c.baseURL, c.collection)
	reqCtx, cancel := c.requestContext(ctx)
	defer cancel()
	point := PointQuery{
		ID:     recipe.ID,
		Vector: vector,
		Payload: UpsertRecipePayload{
			Title:        recipe.Title,
			Ingredients:  recipe.Ingredients,
			Diet:         recipe.Diet,
			Cuisine:      recipe.Cuisine,
			CookTime:     recipe.CookTimeMin,
			Calories:     recipe.Calories,
			Source:       recipe.Source,
			Instructions: recipe.Instructions,
		},
	}
	if recipe.CreatedAt != nil {
		point.Payload.CreatedAt = recipe.CreatedAt.Format(time.RFC3339)
	}

	body := UpsertRecipeReq{
		Points: []PointQuery{point},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("UpsertRecipe:qdrant upsert failed: %s", string(data))
	}
	return nil
}

func (c *QdrantClient) requestContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, c.timeout)
}
