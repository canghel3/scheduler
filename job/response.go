package job

type Response struct {
	id  string
	err error
}

func NewResponse(id string, err error) Response {
	return Response{
		id:  id,
		err: err,
	}
}

func (r Response) ID() string {
	return r.id
}

func (r Response) Err() error {
	return r.err
}
