package job

type Response struct {
	id   string
	err  error
	data any
}

func NewResponse(id string, err error, data any) Response {
	return Response{
		id:   id,
		err:  err,
		data: data,
	}
}

func (r Response) ID() string {
	return r.id
}

func (r Response) Err() error {
	return r.err
}

func (r Response) Data() any {
	return r.data
}
