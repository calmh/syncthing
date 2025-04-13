// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/v2/config_service.proto

package apiv2connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v2 "github.com/syncthing/syncthing/internal/gen/api/v2"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// ConfigurationServiceName is the fully-qualified name of the ConfigurationService service.
	ConfigurationServiceName = "api.v2.ConfigurationService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ConfigurationServiceGetConfigurationProcedure is the fully-qualified name of the
	// ConfigurationService's GetConfiguration RPC.
	ConfigurationServiceGetConfigurationProcedure = "/api.v2.ConfigurationService/GetConfiguration"
	// ConfigurationServiceWatchConfigurationProcedure is the fully-qualified name of the
	// ConfigurationService's WatchConfiguration RPC.
	ConfigurationServiceWatchConfigurationProcedure = "/api.v2.ConfigurationService/WatchConfiguration"
	// ConfigurationServiceUpdateOptionsProcedure is the fully-qualified name of the
	// ConfigurationService's UpdateOptions RPC.
	ConfigurationServiceUpdateOptionsProcedure = "/api.v2.ConfigurationService/UpdateOptions"
	// ConfigurationServiceAddDeviceProcedure is the fully-qualified name of the ConfigurationService's
	// AddDevice RPC.
	ConfigurationServiceAddDeviceProcedure = "/api.v2.ConfigurationService/AddDevice"
	// ConfigurationServiceRemoveDeviceProcedure is the fully-qualified name of the
	// ConfigurationService's RemoveDevice RPC.
	ConfigurationServiceRemoveDeviceProcedure = "/api.v2.ConfigurationService/RemoveDevice"
	// ConfigurationServiceUpdateDeviceProcedure is the fully-qualified name of the
	// ConfigurationService's UpdateDevice RPC.
	ConfigurationServiceUpdateDeviceProcedure = "/api.v2.ConfigurationService/UpdateDevice"
	// ConfigurationServiceAddFolderProcedure is the fully-qualified name of the ConfigurationService's
	// AddFolder RPC.
	ConfigurationServiceAddFolderProcedure = "/api.v2.ConfigurationService/AddFolder"
	// ConfigurationServiceRemoveFolderProcedure is the fully-qualified name of the
	// ConfigurationService's RemoveFolder RPC.
	ConfigurationServiceRemoveFolderProcedure = "/api.v2.ConfigurationService/RemoveFolder"
	// ConfigurationServiceUpdateFolderProcedure is the fully-qualified name of the
	// ConfigurationService's UpdateFolder RPC.
	ConfigurationServiceUpdateFolderProcedure = "/api.v2.ConfigurationService/UpdateFolder"
)

// ConfigurationServiceClient is a client for the api.v2.ConfigurationService service.
type ConfigurationServiceClient interface {
	GetConfiguration(context.Context, *connect.Request[v2.GetConfigurationRequest]) (*connect.Response[v2.GetConfigurationResponse], error)
	WatchConfiguration(context.Context, *connect.Request[v2.WatchConfigurationRequest]) (*connect.ServerStreamForClient[v2.GetConfigurationResponse], error)
	UpdateOptions(context.Context, *connect.Request[v2.UpdateOptionsRequest]) (*connect.Response[v2.UpdateOptionsResponse], error)
	AddDevice(context.Context, *connect.Request[v2.AddDeviceRequest]) (*connect.Response[v2.AddDeviceResponse], error)
	RemoveDevice(context.Context, *connect.Request[v2.RemoveDeviceRequest]) (*connect.Response[v2.RemoveDeviceResponse], error)
	UpdateDevice(context.Context, *connect.Request[v2.UpdateDeviceRequest]) (*connect.Response[v2.UpdateDeviceResponse], error)
	AddFolder(context.Context, *connect.Request[v2.AddFolderRequest]) (*connect.Response[v2.AddFolderResponse], error)
	RemoveFolder(context.Context, *connect.Request[v2.RemoveFolderRequest]) (*connect.Response[v2.RemoveFolderResponse], error)
	UpdateFolder(context.Context, *connect.Request[v2.UpdateFolderRequest]) (*connect.Response[v2.UpdateFolderResponse], error)
}

// NewConfigurationServiceClient constructs a client for the api.v2.ConfigurationService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewConfigurationServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ConfigurationServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	configurationServiceMethods := v2.File_api_v2_config_service_proto.Services().ByName("ConfigurationService").Methods()
	return &configurationServiceClient{
		getConfiguration: connect.NewClient[v2.GetConfigurationRequest, v2.GetConfigurationResponse](
			httpClient,
			baseURL+ConfigurationServiceGetConfigurationProcedure,
			connect.WithSchema(configurationServiceMethods.ByName("GetConfiguration")),
			connect.WithClientOptions(opts...),
		),
		watchConfiguration: connect.NewClient[v2.WatchConfigurationRequest, v2.GetConfigurationResponse](
			httpClient,
			baseURL+ConfigurationServiceWatchConfigurationProcedure,
			connect.WithSchema(configurationServiceMethods.ByName("WatchConfiguration")),
			connect.WithClientOptions(opts...),
		),
		updateOptions: connect.NewClient[v2.UpdateOptionsRequest, v2.UpdateOptionsResponse](
			httpClient,
			baseURL+ConfigurationServiceUpdateOptionsProcedure,
			connect.WithSchema(configurationServiceMethods.ByName("UpdateOptions")),
			connect.WithClientOptions(opts...),
		),
		addDevice: connect.NewClient[v2.AddDeviceRequest, v2.AddDeviceResponse](
			httpClient,
			baseURL+ConfigurationServiceAddDeviceProcedure,
			connect.WithSchema(configurationServiceMethods.ByName("AddDevice")),
			connect.WithClientOptions(opts...),
		),
		removeDevice: connect.NewClient[v2.RemoveDeviceRequest, v2.RemoveDeviceResponse](
			httpClient,
			baseURL+ConfigurationServiceRemoveDeviceProcedure,
			connect.WithSchema(configurationServiceMethods.ByName("RemoveDevice")),
			connect.WithClientOptions(opts...),
		),
		updateDevice: connect.NewClient[v2.UpdateDeviceRequest, v2.UpdateDeviceResponse](
			httpClient,
			baseURL+ConfigurationServiceUpdateDeviceProcedure,
			connect.WithSchema(configurationServiceMethods.ByName("UpdateDevice")),
			connect.WithClientOptions(opts...),
		),
		addFolder: connect.NewClient[v2.AddFolderRequest, v2.AddFolderResponse](
			httpClient,
			baseURL+ConfigurationServiceAddFolderProcedure,
			connect.WithSchema(configurationServiceMethods.ByName("AddFolder")),
			connect.WithClientOptions(opts...),
		),
		removeFolder: connect.NewClient[v2.RemoveFolderRequest, v2.RemoveFolderResponse](
			httpClient,
			baseURL+ConfigurationServiceRemoveFolderProcedure,
			connect.WithSchema(configurationServiceMethods.ByName("RemoveFolder")),
			connect.WithClientOptions(opts...),
		),
		updateFolder: connect.NewClient[v2.UpdateFolderRequest, v2.UpdateFolderResponse](
			httpClient,
			baseURL+ConfigurationServiceUpdateFolderProcedure,
			connect.WithSchema(configurationServiceMethods.ByName("UpdateFolder")),
			connect.WithClientOptions(opts...),
		),
	}
}

// configurationServiceClient implements ConfigurationServiceClient.
type configurationServiceClient struct {
	getConfiguration   *connect.Client[v2.GetConfigurationRequest, v2.GetConfigurationResponse]
	watchConfiguration *connect.Client[v2.WatchConfigurationRequest, v2.GetConfigurationResponse]
	updateOptions      *connect.Client[v2.UpdateOptionsRequest, v2.UpdateOptionsResponse]
	addDevice          *connect.Client[v2.AddDeviceRequest, v2.AddDeviceResponse]
	removeDevice       *connect.Client[v2.RemoveDeviceRequest, v2.RemoveDeviceResponse]
	updateDevice       *connect.Client[v2.UpdateDeviceRequest, v2.UpdateDeviceResponse]
	addFolder          *connect.Client[v2.AddFolderRequest, v2.AddFolderResponse]
	removeFolder       *connect.Client[v2.RemoveFolderRequest, v2.RemoveFolderResponse]
	updateFolder       *connect.Client[v2.UpdateFolderRequest, v2.UpdateFolderResponse]
}

// GetConfiguration calls api.v2.ConfigurationService.GetConfiguration.
func (c *configurationServiceClient) GetConfiguration(ctx context.Context, req *connect.Request[v2.GetConfigurationRequest]) (*connect.Response[v2.GetConfigurationResponse], error) {
	return c.getConfiguration.CallUnary(ctx, req)
}

// WatchConfiguration calls api.v2.ConfigurationService.WatchConfiguration.
func (c *configurationServiceClient) WatchConfiguration(ctx context.Context, req *connect.Request[v2.WatchConfigurationRequest]) (*connect.ServerStreamForClient[v2.GetConfigurationResponse], error) {
	return c.watchConfiguration.CallServerStream(ctx, req)
}

// UpdateOptions calls api.v2.ConfigurationService.UpdateOptions.
func (c *configurationServiceClient) UpdateOptions(ctx context.Context, req *connect.Request[v2.UpdateOptionsRequest]) (*connect.Response[v2.UpdateOptionsResponse], error) {
	return c.updateOptions.CallUnary(ctx, req)
}

// AddDevice calls api.v2.ConfigurationService.AddDevice.
func (c *configurationServiceClient) AddDevice(ctx context.Context, req *connect.Request[v2.AddDeviceRequest]) (*connect.Response[v2.AddDeviceResponse], error) {
	return c.addDevice.CallUnary(ctx, req)
}

// RemoveDevice calls api.v2.ConfigurationService.RemoveDevice.
func (c *configurationServiceClient) RemoveDevice(ctx context.Context, req *connect.Request[v2.RemoveDeviceRequest]) (*connect.Response[v2.RemoveDeviceResponse], error) {
	return c.removeDevice.CallUnary(ctx, req)
}

// UpdateDevice calls api.v2.ConfigurationService.UpdateDevice.
func (c *configurationServiceClient) UpdateDevice(ctx context.Context, req *connect.Request[v2.UpdateDeviceRequest]) (*connect.Response[v2.UpdateDeviceResponse], error) {
	return c.updateDevice.CallUnary(ctx, req)
}

// AddFolder calls api.v2.ConfigurationService.AddFolder.
func (c *configurationServiceClient) AddFolder(ctx context.Context, req *connect.Request[v2.AddFolderRequest]) (*connect.Response[v2.AddFolderResponse], error) {
	return c.addFolder.CallUnary(ctx, req)
}

// RemoveFolder calls api.v2.ConfigurationService.RemoveFolder.
func (c *configurationServiceClient) RemoveFolder(ctx context.Context, req *connect.Request[v2.RemoveFolderRequest]) (*connect.Response[v2.RemoveFolderResponse], error) {
	return c.removeFolder.CallUnary(ctx, req)
}

// UpdateFolder calls api.v2.ConfigurationService.UpdateFolder.
func (c *configurationServiceClient) UpdateFolder(ctx context.Context, req *connect.Request[v2.UpdateFolderRequest]) (*connect.Response[v2.UpdateFolderResponse], error) {
	return c.updateFolder.CallUnary(ctx, req)
}

// ConfigurationServiceHandler is an implementation of the api.v2.ConfigurationService service.
type ConfigurationServiceHandler interface {
	GetConfiguration(context.Context, *connect.Request[v2.GetConfigurationRequest]) (*connect.Response[v2.GetConfigurationResponse], error)
	WatchConfiguration(context.Context, *connect.Request[v2.WatchConfigurationRequest], *connect.ServerStream[v2.GetConfigurationResponse]) error
	UpdateOptions(context.Context, *connect.Request[v2.UpdateOptionsRequest]) (*connect.Response[v2.UpdateOptionsResponse], error)
	AddDevice(context.Context, *connect.Request[v2.AddDeviceRequest]) (*connect.Response[v2.AddDeviceResponse], error)
	RemoveDevice(context.Context, *connect.Request[v2.RemoveDeviceRequest]) (*connect.Response[v2.RemoveDeviceResponse], error)
	UpdateDevice(context.Context, *connect.Request[v2.UpdateDeviceRequest]) (*connect.Response[v2.UpdateDeviceResponse], error)
	AddFolder(context.Context, *connect.Request[v2.AddFolderRequest]) (*connect.Response[v2.AddFolderResponse], error)
	RemoveFolder(context.Context, *connect.Request[v2.RemoveFolderRequest]) (*connect.Response[v2.RemoveFolderResponse], error)
	UpdateFolder(context.Context, *connect.Request[v2.UpdateFolderRequest]) (*connect.Response[v2.UpdateFolderResponse], error)
}

// NewConfigurationServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewConfigurationServiceHandler(svc ConfigurationServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	configurationServiceMethods := v2.File_api_v2_config_service_proto.Services().ByName("ConfigurationService").Methods()
	configurationServiceGetConfigurationHandler := connect.NewUnaryHandler(
		ConfigurationServiceGetConfigurationProcedure,
		svc.GetConfiguration,
		connect.WithSchema(configurationServiceMethods.ByName("GetConfiguration")),
		connect.WithHandlerOptions(opts...),
	)
	configurationServiceWatchConfigurationHandler := connect.NewServerStreamHandler(
		ConfigurationServiceWatchConfigurationProcedure,
		svc.WatchConfiguration,
		connect.WithSchema(configurationServiceMethods.ByName("WatchConfiguration")),
		connect.WithHandlerOptions(opts...),
	)
	configurationServiceUpdateOptionsHandler := connect.NewUnaryHandler(
		ConfigurationServiceUpdateOptionsProcedure,
		svc.UpdateOptions,
		connect.WithSchema(configurationServiceMethods.ByName("UpdateOptions")),
		connect.WithHandlerOptions(opts...),
	)
	configurationServiceAddDeviceHandler := connect.NewUnaryHandler(
		ConfigurationServiceAddDeviceProcedure,
		svc.AddDevice,
		connect.WithSchema(configurationServiceMethods.ByName("AddDevice")),
		connect.WithHandlerOptions(opts...),
	)
	configurationServiceRemoveDeviceHandler := connect.NewUnaryHandler(
		ConfigurationServiceRemoveDeviceProcedure,
		svc.RemoveDevice,
		connect.WithSchema(configurationServiceMethods.ByName("RemoveDevice")),
		connect.WithHandlerOptions(opts...),
	)
	configurationServiceUpdateDeviceHandler := connect.NewUnaryHandler(
		ConfigurationServiceUpdateDeviceProcedure,
		svc.UpdateDevice,
		connect.WithSchema(configurationServiceMethods.ByName("UpdateDevice")),
		connect.WithHandlerOptions(opts...),
	)
	configurationServiceAddFolderHandler := connect.NewUnaryHandler(
		ConfigurationServiceAddFolderProcedure,
		svc.AddFolder,
		connect.WithSchema(configurationServiceMethods.ByName("AddFolder")),
		connect.WithHandlerOptions(opts...),
	)
	configurationServiceRemoveFolderHandler := connect.NewUnaryHandler(
		ConfigurationServiceRemoveFolderProcedure,
		svc.RemoveFolder,
		connect.WithSchema(configurationServiceMethods.ByName("RemoveFolder")),
		connect.WithHandlerOptions(opts...),
	)
	configurationServiceUpdateFolderHandler := connect.NewUnaryHandler(
		ConfigurationServiceUpdateFolderProcedure,
		svc.UpdateFolder,
		connect.WithSchema(configurationServiceMethods.ByName("UpdateFolder")),
		connect.WithHandlerOptions(opts...),
	)
	return "/api.v2.ConfigurationService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ConfigurationServiceGetConfigurationProcedure:
			configurationServiceGetConfigurationHandler.ServeHTTP(w, r)
		case ConfigurationServiceWatchConfigurationProcedure:
			configurationServiceWatchConfigurationHandler.ServeHTTP(w, r)
		case ConfigurationServiceUpdateOptionsProcedure:
			configurationServiceUpdateOptionsHandler.ServeHTTP(w, r)
		case ConfigurationServiceAddDeviceProcedure:
			configurationServiceAddDeviceHandler.ServeHTTP(w, r)
		case ConfigurationServiceRemoveDeviceProcedure:
			configurationServiceRemoveDeviceHandler.ServeHTTP(w, r)
		case ConfigurationServiceUpdateDeviceProcedure:
			configurationServiceUpdateDeviceHandler.ServeHTTP(w, r)
		case ConfigurationServiceAddFolderProcedure:
			configurationServiceAddFolderHandler.ServeHTTP(w, r)
		case ConfigurationServiceRemoveFolderProcedure:
			configurationServiceRemoveFolderHandler.ServeHTTP(w, r)
		case ConfigurationServiceUpdateFolderProcedure:
			configurationServiceUpdateFolderHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedConfigurationServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedConfigurationServiceHandler struct{}

func (UnimplementedConfigurationServiceHandler) GetConfiguration(context.Context, *connect.Request[v2.GetConfigurationRequest]) (*connect.Response[v2.GetConfigurationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ConfigurationService.GetConfiguration is not implemented"))
}

func (UnimplementedConfigurationServiceHandler) WatchConfiguration(context.Context, *connect.Request[v2.WatchConfigurationRequest], *connect.ServerStream[v2.GetConfigurationResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ConfigurationService.WatchConfiguration is not implemented"))
}

func (UnimplementedConfigurationServiceHandler) UpdateOptions(context.Context, *connect.Request[v2.UpdateOptionsRequest]) (*connect.Response[v2.UpdateOptionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ConfigurationService.UpdateOptions is not implemented"))
}

func (UnimplementedConfigurationServiceHandler) AddDevice(context.Context, *connect.Request[v2.AddDeviceRequest]) (*connect.Response[v2.AddDeviceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ConfigurationService.AddDevice is not implemented"))
}

func (UnimplementedConfigurationServiceHandler) RemoveDevice(context.Context, *connect.Request[v2.RemoveDeviceRequest]) (*connect.Response[v2.RemoveDeviceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ConfigurationService.RemoveDevice is not implemented"))
}

func (UnimplementedConfigurationServiceHandler) UpdateDevice(context.Context, *connect.Request[v2.UpdateDeviceRequest]) (*connect.Response[v2.UpdateDeviceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ConfigurationService.UpdateDevice is not implemented"))
}

func (UnimplementedConfigurationServiceHandler) AddFolder(context.Context, *connect.Request[v2.AddFolderRequest]) (*connect.Response[v2.AddFolderResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ConfigurationService.AddFolder is not implemented"))
}

func (UnimplementedConfigurationServiceHandler) RemoveFolder(context.Context, *connect.Request[v2.RemoveFolderRequest]) (*connect.Response[v2.RemoveFolderResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ConfigurationService.RemoveFolder is not implemented"))
}

func (UnimplementedConfigurationServiceHandler) UpdateFolder(context.Context, *connect.Request[v2.UpdateFolderRequest]) (*connect.Response[v2.UpdateFolderResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.v2.ConfigurationService.UpdateFolder is not implemented"))
}
