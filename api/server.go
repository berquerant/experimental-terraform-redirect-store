package api

import "context"

var (
	_ Server     = NewServerImpl(nil)
	_ Redirector = NewServerImpl(nil)
)

type Server interface {
	Scan(ctx context.Context, r *ScanRequest) (*ScanResponse, error)
	Get(ctx context.Context, r *GetRequest) (*GetResponse, error)
	Put(ctx context.Context, r *PutRequest) (*PutResponse, error)
	Delete(ctx context.Context, r *DeleteRequest) (*DeleteResponse, error)
}

type Redirector interface {
	Redirect(ctx context.Context, r *RedirectRequest) (*RedirectResponse, error)
}

func NewServerImpl(db Database) *ServerImpl {
	return &ServerImpl{
		db: db,
	}
}

type ServerImpl struct {
	db Database
}

func (s *ServerImpl) Scan(ctx context.Context, _ *ScanRequest) (*ScanResponse, error) {
	records, err := s.db.Scan(ctx)
	if err != nil {
		return &ScanResponse{
			Error: err.Error(),
		}, err
	}
	return &ScanResponse{
		Records: records,
	}, nil
}

func (s *ServerImpl) Get(ctx context.Context, r *GetRequest) (*GetResponse, error) {
	record, err := s.db.Get(ctx, r.Name)
	if err != nil {
		return &GetResponse{
			Error: err.Error(),
		}, err
	}
	return &GetResponse{
		Record: record,
	}, nil
}

func (s *ServerImpl) Put(ctx context.Context, r *PutRequest) (*PutResponse, error) {
	if err := s.db.Put(ctx, r.Record); err != nil {
		return &PutResponse{
			Error: err.Error(),
		}, err
	}
	return &PutResponse{
		Record: r.Record,
	}, nil
}

func (s *ServerImpl) Delete(ctx context.Context, r *DeleteRequest) (*DeleteResponse, error) {
	if err := s.db.Delete(ctx, r.Name); err != nil {
		return &DeleteResponse{
			Error: err.Error(),
		}, err
	}
	return nil, nil
}

func (s *ServerImpl) Redirect(ctx context.Context, r *RedirectRequest) (*RedirectResponse, error) {
	record, err := s.db.Get(ctx, r.Name)
	if err != nil {
		return &RedirectResponse{
			Error: err.Error(),
		}, err
	}
	return &RedirectResponse{
		To: record.To,
	}, nil
}
