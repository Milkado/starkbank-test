package queue

type (
	Invoice struct {
		Amount float64  `json:"amount" xml:"amount" form:"amount" query:"amount"`
		Name   string `json:"name" xml:"name" form:"name" query:"name"`
		TaxId  string `json:"tax_id" xml:"tax_id" form:"tax_id" query:"tax_id"`
	}

	QueueRequest struct {
		length int
		front *ListRequest
		rear *ListRequest
	}

	ListRequest struct {
		data *Invoice
		next *ListRequest
	}
)
func (ln *QueueRequest) LengthRequest() int {
	return ln.length
}

func (ln *QueueRequest) isRequestEmpty() bool {
	return ln.length == 0
}

func (q *QueueRequest) EnqueueRequest(data Invoice) {
	temp := &ListRequest{data: &data}

	if q.isRequestEmpty() {
		q.front = temp
	} else {
		q.rear.next = temp
	}

	q.rear = temp
	q.length++
}

func (q *QueueRequest) DequeueRequest() *Invoice {
	if q.isRequestEmpty() {
		return &Invoice{}
	}

	result := q.front.data
	q.front = q.front.next

	if q.front == nil {
		q.rear = nil
	}

	q.length--
	return result
}