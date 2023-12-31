package model

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Qty         int     `json:"qty"`
	Description string  `json:"description"`
}

type ProductsResponse struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Qty         int     `json:"qty"`
	Description string  `json:"description"`
}

type Handler = func([]byte) error

type RabbitMQConusmerParams struct {
	WorkerPoolSize int
	Exchange       string
	QueueName      string
	BindingKey     string
	ConsumerTag    string
}
