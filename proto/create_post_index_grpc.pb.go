package search

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion7

const (
	SearchEngine_CreatePostIndex_FullMethodName = "/SearchEngine/CreatePostIndex"
	SearchEngine_Search_FullMethodName          = "/SearchEngine/Search"
)

type SearchEngineClient interface {
	CreatePostIndex(ctx context.Context, in *CreatePostIndexRequest, opts ...grpc.CallOption) (*CreatePostIndexResponse, error)
	Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error)
}

type searchEngineClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchEngineClient(cc grpc.ClientConnInterface) SearchEngineClient {
	return &searchEngineClient{cc}
}

func (c *searchEngineClient) CreatePostIndex(ctx context.Context, in *CreatePostIndexRequest, opts ...grpc.CallOption) (*CreatePostIndexResponse, error) {
	out := new(CreatePostIndexResponse)
	err := c.cc.Invoke(ctx, SearchEngine_CreatePostIndex_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchEngineClient) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (*SearchResponse, error) {
	out := new(SearchResponse)
	err := c.cc.Invoke(ctx, SearchEngine_Search_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type SearchEngineServer interface {
	CreatePostIndex(context.Context, *CreatePostIndexRequest) (*CreatePostIndexResponse, error)
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
	mustEmbedUnimplementedSearchEngineServer()
}

type UnimplementedSearchEngineServer struct {
}

func (UnimplementedSearchEngineServer) CreatePostIndex(context.Context, *CreatePostIndexRequest) (*CreatePostIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePostIndex not implemented")
}
func (UnimplementedSearchEngineServer) Search(context.Context, *SearchRequest) (*SearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedSearchEngineServer) mustEmbedUnimplementedSearchEngineServer() {}

type UnsafeSearchEngineServer interface {
	mustEmbedUnimplementedSearchEngineServer()
}

func RegisterSearchEngineServer(s grpc.ServiceRegistrar, srv SearchEngineServer) {
	s.RegisterService(&SearchEngine_ServiceDesc, srv)
}

func _SearchEngine_CreatePostIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchEngineServer).CreatePostIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchEngine_CreatePostIndex_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchEngineServer).CreatePostIndex(ctx, req.(*CreatePostIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchEngine_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchEngineServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SearchEngine_Search_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchEngineServer).Search(ctx, req.(*SearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var SearchEngine_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SearchEngine",
	HandlerType: (*SearchEngineServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePostIndex",
			Handler:    _SearchEngine_CreatePostIndex_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _SearchEngine_Search_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "create_post_index.proto",
}
