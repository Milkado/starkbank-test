package queue

type (
	CreatedInvoice struct {
		Id int64 `json:"id" xml:"id" form:"id" query:"id"`
		Name   string `json:"name" xml:"name" form:"name" query:"name"`
		Amount float64  `json:"amount" xml:"amount" form:"amount" query:"amount"`
		Fee float64 `json:"fee" xml:"fee" form:"fee" query:"fee"`

	}

	QueueCreated struct {
		length int
		front *ListCreated
		rear *ListCreated
	}

	ListCreated struct {
		data *CreatedInvoice
		next *ListCreated
	}
)
func (ln *QueueCreated) LengthCreated() int {
	return ln.length
}

func (ln *QueueCreated) isCreatedEmpy() bool {
	return ln.length == 0
}

func (q *QueueCreated) EnqueueCreated(data CreatedInvoice) {
	temp := &ListCreated{data: &data}

	if q.isCreatedEmpy() {
		q.front = temp
	} else {
		q.rear.next = temp
	}

	q.rear = temp
	q.length++
}

func (q *QueueCreated) DequeueCreated() *CreatedInvoice {
	if q.isCreatedEmpy() {
		return &CreatedInvoice{}
	}

	result := q.front.data
	q.front = q.front.next

	if q.front == nil {
		q.rear = nil
	}

	q.length--
	return result
}