// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DaemonServiceClient is the client API for DaemonService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DaemonServiceClient interface {
	// Login uses setup key to prepare configuration for the daemon.
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	// WaitSSOLogin uses the userCode to validate the TokenInfo and
	// waits for the user to continue with the login on a browser
	WaitSSOLogin(ctx context.Context, in *WaitSSOLoginRequest, opts ...grpc.CallOption) (*WaitSSOLoginResponse, error)
	// Up starts engine work in the daemon.
	Up(ctx context.Context, in *UpRequest, opts ...grpc.CallOption) (*UpResponse, error)
	// Status of the service.
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	// Down engine work in the daemon.
	Down(ctx context.Context, in *DownRequest, opts ...grpc.CallOption) (*DownResponse, error)
	// GetConfig of the daemon.
	GetConfig(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error)
	// List available network routes
	ListRoutes(ctx context.Context, in *ListRoutesRequest, opts ...grpc.CallOption) (*ListRoutesResponse, error)
	// Select specific routes
	SelectRoutes(ctx context.Context, in *SelectRoutesRequest, opts ...grpc.CallOption) (*SelectRoutesResponse, error)
	// Deselect specific routes
	DeselectRoutes(ctx context.Context, in *SelectRoutesRequest, opts ...grpc.CallOption) (*SelectRoutesResponse, error)
	DebugBundle(ctx context.Context, in *DebugBundleRequest, opts ...grpc.CallOption) (*DebugBundleResponse, error)
}

type daemonServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDaemonServiceClient(cc grpc.ClientConnInterface) DaemonServiceClient {
	return &daemonServiceClient{cc}
}

func (c *daemonServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) WaitSSOLogin(ctx context.Context, in *WaitSSOLoginRequest, opts ...grpc.CallOption) (*WaitSSOLoginResponse, error) {
	out := new(WaitSSOLoginResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/WaitSSOLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) Up(ctx context.Context, in *UpRequest, opts ...grpc.CallOption) (*UpResponse, error) {
	out := new(UpResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/Up", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) Down(ctx context.Context, in *DownRequest, opts ...grpc.CallOption) (*DownResponse, error) {
	out := new(DownResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/Down", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) GetConfig(ctx context.Context, in *GetConfigRequest, opts ...grpc.CallOption) (*GetConfigResponse, error) {
	out := new(GetConfigResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/GetConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) ListRoutes(ctx context.Context, in *ListRoutesRequest, opts ...grpc.CallOption) (*ListRoutesResponse, error) {
	out := new(ListRoutesResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/ListRoutes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) SelectRoutes(ctx context.Context, in *SelectRoutesRequest, opts ...grpc.CallOption) (*SelectRoutesResponse, error) {
	out := new(SelectRoutesResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/SelectRoutes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) DeselectRoutes(ctx context.Context, in *SelectRoutesRequest, opts ...grpc.CallOption) (*SelectRoutesResponse, error) {
	out := new(SelectRoutesResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/DeselectRoutes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *daemonServiceClient) DebugBundle(ctx context.Context, in *DebugBundleRequest, opts ...grpc.CallOption) (*DebugBundleResponse, error) {
	out := new(DebugBundleResponse)
	err := c.cc.Invoke(ctx, "/daemon.DaemonService/DebugBundle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DaemonServiceServer is the server API for DaemonService service.
// All implementations must embed UnimplementedDaemonServiceServer
// for forward compatibility
type DaemonServiceServer interface {
	// Login uses setup key to prepare configuration for the daemon.
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	// WaitSSOLogin uses the userCode to validate the TokenInfo and
	// waits for the user to continue with the login on a browser
	WaitSSOLogin(context.Context, *WaitSSOLoginRequest) (*WaitSSOLoginResponse, error)
	// Up starts engine work in the daemon.
	Up(context.Context, *UpRequest) (*UpResponse, error)
	// Status of the service.
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	// Down engine work in the daemon.
	Down(context.Context, *DownRequest) (*DownResponse, error)
	// GetConfig of the daemon.
	GetConfig(context.Context, *GetConfigRequest) (*GetConfigResponse, error)
	// List available network routes
	ListRoutes(context.Context, *ListRoutesRequest) (*ListRoutesResponse, error)
	// Select specific routes
	SelectRoutes(context.Context, *SelectRoutesRequest) (*SelectRoutesResponse, error)
	// Deselect specific routes
	DeselectRoutes(context.Context, *SelectRoutesRequest) (*SelectRoutesResponse, error)
	DebugBundle(context.Context, *DebugBundleRequest) (*DebugBundleResponse, error)
	mustEmbedUnimplementedDaemonServiceServer()
}

// UnimplementedDaemonServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDaemonServiceServer struct {
}

func (UnimplementedDaemonServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedDaemonServiceServer) WaitSSOLogin(context.Context, *WaitSSOLoginRequest) (*WaitSSOLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WaitSSOLogin not implemented")
}
func (UnimplementedDaemonServiceServer) Up(context.Context, *UpRequest) (*UpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Up not implemented")
}
func (UnimplementedDaemonServiceServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedDaemonServiceServer) Down(context.Context, *DownRequest) (*DownResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Down not implemented")
}
func (UnimplementedDaemonServiceServer) GetConfig(context.Context, *GetConfigRequest) (*GetConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}
func (UnimplementedDaemonServiceServer) ListRoutes(context.Context, *ListRoutesRequest) (*ListRoutesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRoutes not implemented")
}
func (UnimplementedDaemonServiceServer) SelectRoutes(context.Context, *SelectRoutesRequest) (*SelectRoutesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SelectRoutes not implemented")
}
func (UnimplementedDaemonServiceServer) DeselectRoutes(context.Context, *SelectRoutesRequest) (*SelectRoutesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeselectRoutes not implemented")
}
func (UnimplementedDaemonServiceServer) DebugBundle(context.Context, *DebugBundleRequest) (*DebugBundleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DebugBundle not implemented")
}
func (UnimplementedDaemonServiceServer) mustEmbedUnimplementedDaemonServiceServer() {}

// UnsafeDaemonServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DaemonServiceServer will
// result in compilation errors.
type UnsafeDaemonServiceServer interface {
	mustEmbedUnimplementedDaemonServiceServer()
}

func RegisterDaemonServiceServer(s grpc.ServiceRegistrar, srv DaemonServiceServer) {
	s.RegisterService(&DaemonService_ServiceDesc, srv)
}

func _DaemonService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_WaitSSOLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WaitSSOLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).WaitSSOLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/WaitSSOLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).WaitSSOLogin(ctx, req.(*WaitSSOLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_Up_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).Up(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/Up",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).Up(ctx, req.(*UpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_Down_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).Down(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/Down",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).Down(ctx, req.(*DownRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/GetConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).GetConfig(ctx, req.(*GetConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_ListRoutes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoutesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).ListRoutes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/ListRoutes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).ListRoutes(ctx, req.(*ListRoutesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_SelectRoutes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SelectRoutesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).SelectRoutes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/SelectRoutes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).SelectRoutes(ctx, req.(*SelectRoutesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_DeselectRoutes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SelectRoutesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).DeselectRoutes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/DeselectRoutes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).DeselectRoutes(ctx, req.(*SelectRoutesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DaemonService_DebugBundle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DebugBundleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DaemonServiceServer).DebugBundle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/daemon.DaemonService/DebugBundle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DaemonServiceServer).DebugBundle(ctx, req.(*DebugBundleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DaemonService_ServiceDesc is the grpc.ServiceDesc for DaemonService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DaemonService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "daemon.DaemonService",
	HandlerType: (*DaemonServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _DaemonService_Login_Handler,
		},
		{
			MethodName: "WaitSSOLogin",
			Handler:    _DaemonService_WaitSSOLogin_Handler,
		},
		{
			MethodName: "Up",
			Handler:    _DaemonService_Up_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _DaemonService_Status_Handler,
		},
		{
			MethodName: "Down",
			Handler:    _DaemonService_Down_Handler,
		},
		{
			MethodName: "GetConfig",
			Handler:    _DaemonService_GetConfig_Handler,
		},
		{
			MethodName: "ListRoutes",
			Handler:    _DaemonService_ListRoutes_Handler,
		},
		{
			MethodName: "SelectRoutes",
			Handler:    _DaemonService_SelectRoutes_Handler,
		},
		{
			MethodName: "DeselectRoutes",
			Handler:    _DaemonService_DeselectRoutes_Handler,
		},
		{
			MethodName: "DebugBundle",
			Handler:    _DaemonService_DebugBundle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "daemon.proto",
}
