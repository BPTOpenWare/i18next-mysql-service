package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	resourceapi "godev.bluffpointtech.com/i18nextservice/proto/resources"
	envUtil "godev.bluffpointtech.com/i18nextservice/util"
)

type Environment struct {
	EnvConfig envUtil.Config
}

var Env *Environment

// logResponseWriter wraps http.ResponseWriter to capture status and bytes written.
type logResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (lw *logResponseWriter) WriteHeader(code int) {
	lw.status = code
	lw.ResponseWriter.WriteHeader(code)
}

func (lw *logResponseWriter) Write(b []byte) (int, error) {
	n, err := lw.ResponseWriter.Write(b)
	lw.bytes += n
	return n, err
}

// LoggingMiddleware logs each request handled by the gateway.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &logResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)
		log.Printf("gateway request method=%s url=%s status=%d bytes=%d remote=%s ua=%q duration=%s",
			r.Method, r.URL.String(), lw.status, lw.bytes, r.RemoteAddr, r.UserAgent(), time.Since(start))
	})
}

func main() {

	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// create environment
	Env = &Environment{envUtil.EnvConfig}

	grpcPort := Env.EnvConfig.GetConfig("GrpcServerPort")

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", "0.0.0.0:"+grpcPort)

	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object and attach telemetry (otelgrpc switched to StatsHandler)
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	// Serve gRPC Servesr
	log.Println("Serving gRPC on 0.0.0.0:" + grpcPort)
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// setup grpc client options for rest gateway
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(otelgrpc.WithTracerProvider(tp)),
		),
	}

	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(CustomMetaDataRequestMatcher),
		runtime.WithIncomingHeaderMatcher(CustomHeaderRequestMatcher),
		runtime.WithForwardResponseOption(HttpResponseModifier),
		runtime.WithOutgoingHeaderMatcher(CustomHeaderResponseMatcher),
		runtime.WithErrorHandler(CustomErrorHandler))

	err = resourceapi.RegisterResourcesAPIHandlerFromEndpoint(context.Background(), gwmux, "0.0.0.0:"+grpcPort, opts)
	if err != nil {
		log.Fatalln("Failed to register gRPC gateway:", err)
	}

	gatewayPort := Env.EnvConfig.GetConfig("GatewayServerPort")

	gwServer := &http.Server{
		Addr: ":" + gatewayPort,
		Handler: LoggingMiddleware(
			otelhttp.NewHandler(gwmux, "server",
				otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			),
		),
		ConnState: func(c net.Conn, state http.ConnState) {
			log.Printf("gateway conn remote=%s state=%s", c.RemoteAddr(), state)
		},
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:" + gatewayPort)
	log.Fatalln(gwServer.ListenAndServe())
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdout.New()
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func CustomMetaDataRequestMatcher(ctx context.Context, r *http.Request) metadata.MD {
	md := make(map[string]string)
	if method, ok := runtime.RPCMethod(ctx); ok {
		md["method"] = method // /grpc.gateway.examples.internal.proto.examplepb.LoginService/Login
	}
	if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
		md["pattern"] = pattern // /v1/example/login
	}
	return metadata.New(md)
}

func CustomHeaderRequestMatcher(key string) (string, bool) {
	switch key {
	// haproxy x forwarded for
	case "X-Forwarded-For":
		return "grpcgateway-" + key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func CustomHeaderResponseMatcher(key string) (string, bool) {
	switch key {
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func HttpResponseModifier(ctx context.Context, w http.ResponseWriter, p protoreflect.ProtoMessage) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status codeS
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}

	return nil
}

func CustomErrorHandler(ctx context.Context, m *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	runtime.DefaultHTTPErrorHandler(ctx, m, marshaler, w, r, err)
}
