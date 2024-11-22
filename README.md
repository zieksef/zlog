# zlog

## 1. Overview
A logging library designed to provide a chainable API for logging operations, built on top of the Uber [Zap](https://github.com/uber-go/zap) with [Zerolog (https://github.com/rs/zerolog)](https://github.com/rs/zerolog) style.



## 2. Quick Start

```
zlog.Info().Msg("an example")

{"level":"INFO","ts":"2024-11-12T14:16:47.522+0800","caller":"cmd/main.go:20","msg":"an example"}
```



## 3. Usage

### 3.1 Initialization

```go
// package level initialization with default option
zlog.Init(zlog.DefaultOption())
zlog.Info().Msg("an example")


// object level initialization with default option
logger := zlog.New(zlog.DefaultOption())


// initialization with user defined option
logger := zlog.New(zlogOption{
    Dir:            "logs",
    Filename:       "access.log",
    DisableConsole: false,
    Rotation:       true,
    MaxSize:        20,
    MaxAge:         30,
    SkipLevel:      2,
    Writer:         os.Stderr,
})

```

### 3.2 Chainable Option API

```go
// example of production envrionment initialization
opt := zlog.DefaultOption().
        SetDir(config.Dir).
        SetFilename(config.LogFile).
        SetWriter(os.Stdout).
        SetDisableConsole(true)

logger := zlog.New(opt)
```

### 3.3 Chainable Logging API

```go
// logging with message
logger.Info().
    Str("trace_id", "1234").
    Str("account_id", "abc").
    Int("account_status", 1).
    Msg("query account")
    
    
// logging with formatted message
logger.Info().
    Str("trace_id", "1234").
    Str("account_id", "abc").
    Int("account_status", 1).
    Msgf("query account[%d]", "aid")
    
    
// logging without message
logger.Info().
    Str("trace_id", "1234").
    Str("account_id", "abc").
    Int("account_status", 1).
    Send()
    
  
// logging with error
// {"level":"ERROR","ts":"2024-11-12T14:11:50.488+0800","caller":"runtime/proc.go:271","msg":"query account","trace_id":"1234","account_id":"abc","account_status":1,"error":"record not found"}
logger.Error().
    Str("trace_id", "1234").
    Str("account_id", "abc").
    Int("account_status", 1).
    Err(errors.New("record not found")).
    Msg("failed to query account")
```

### 3.4 Log Pre-defined Context

```go
logger := zlog.New(zlog.DefaultOption())

logger = logger.With("machine_id", 3).
    With("service", "gateway").
    With("router", "card")

logger.Info().Str("trace_id", "1234").Str("card_number", "12345").Msg("update card")


// {"level":"INFO","ts":"2024-11-12T14:27:02.551+0800","caller":"cmd/main.go:26","msg":"update card","machine_id":3,"service":"gateway","router":"card","trace_id":"1234","card_number":"12345"}
```

### 3.5 Passed Standard Context

```go
logger := zlog.New(zlog.DefaultOption())

ctx := context.Background()
ctx = context.WithValue(ctx, "trace_id", "t123")
ctx = context.WithValue(ctx, "token", "jwt_abc")
ctx = context.WithValue(ctx, "api_key", "key_vip_acc")

logger.AddCtxHook(func(ctx context.Context) map[string]interface{} {
		return map[string]interface{}{
                    "trace_id": ctx.Value("trace_id"),
                    "token":    ctx.Value("token"),
                    "api_key":  ctx.Value("api_key"),
		}
    }
)

logger.Info(ctx).Str("account_id", "1234").Str("card_number", "12345").Msg("update card")

// {"level":"INFO","ts":"2024-11-12T14:46:59.926+0800","caller":"cmd/main.go:38","msg":"update card","trace_id":"t123","token":"jwt_abc","api_key":"key_vip_acc","account_id":"1234","card_number":"12345"}
```

### 3.6 Passed Gin Context

```go
logger := zlog.New(zlog.DefaultOption())

ctx := &gin.Context{}
ctx.Set("trace_id", "t123")
ctx.Set("client_id", "c123")
ctx.Set("customer_id", "customer_123")
ctx.Set("account_id", "a123")

logger.AddCtxHook(func(ctx context.Context) map[string]interface{} {
		c := ctx.(*gin.Context)
		traceId, _ := c.Get("trace_id")
		clientId, _ := c.Get("client_id")
		customerId, _ := c.Get("customer_id")
		accountId, _ := c.Get("account_id")

		return map[string]interface{}{
            "trace_id":    traceId,
            "client_id":   clientId,
            "customer_id": customerId,
            "account_id":  accountId,
		}
})

logger.Info(ctx).Str("account_id", "1234").Str("card_number", "12345").Msg("update card")

// {"level":"INFO","ts":"2024-11-14T10:16:33.876+0800","caller":"cmd/main.go:44","msg":"update card","trace_id":"t123","client_id":"c123","customer_id":"customer_123","account_id":"a123","account_id":"1234","card_number":"12345"}
```



## 4. Comparation

### 4.1 Versus Zap

```go
// zap
logger.Info("create card, allowed card currency",
    zap.String("direct_id", ca.DirectId),
    zap.String("account_id", ca.AccountId),
    zap.String("allowed_card_currency", allowedCardCurrency),
    zap.String("req_card_currency", reqCardCurrency),
)

logger.Info("create card, allowed card currency", zap.String("direct_id", ca.DirectId), zap.String("account_id", ca.AccountId), zap.String("allowed_card_currency", allowedCardCurrency), zap.String("req_card_currency", reqCardCurrency))


// shorter, less embedded and much more clear
// no worries about commas and brackets
logger.Info().
		Str("direct_id", ca.DirectId).
		Str("account_id", ca.AccountId).
		Str("allowed_card_currency", allowedCardCurrency).
		Str("req_card_currency", reqCardCurrency).
		Msg("create card, allowed card currency")

logger.Info().Str("direct_id", ca.DirectId).Str("account_id", ca.AccountId).Str("allowed_card_currency", allowedCardCurrency).Str("req_card_currency", reqCardCurrency).Msg("create card, allowed card currency")

```

### 4.2 Versus Current Project

```go
// current
var (
	Lg      *zap.Logger
	TraceLg TraceLogger
	PingLg  *zap.Logger
	Clg     *Clog
)

Lg.Info("example log")
TraceLg.InfoWithCtx(ctx, "example log")
Clg.Info(ctx, "example log")


// more compatible and pluggable
logger.Info().Msg("example log")
logger.Info(ctx).Msg("example log with context")
```

also:
- decoupling log library and developer usage
- less initialization steps, less option with more expansible and unified




## 5. Performance

```
[Apple M3 Pro]

goos: darwin
goarch: arm64

# 2 cores
BenchmarkLogEmpty-2                     	 2997499	       391.2 ns/op	     256 B/op	       2 allocs/op
BenchmarkLog-2                          	 2344003	       563.6 ns/op	     256 B/op	       2 allocs/op
BenchmarkLog4Fields-2                   	 2338986	       510.7 ns/op	     256 B/op	       2 allocs/op
BenchmarkLog10Fields-2                  	 1883809	       629.2 ns/op	     256 B/op	       2 allocs/op
BenchmarkLog10FieldsComplex-2           	 1684208	       710.9 ns/op	     280 B/op	       3 allocs/op
BenchmarkLog10FieldsComplexAny-2        	  405367	      2949 ns/op	    1121 B/op	      23 allocs/op
BenchmarkLogPredefinedContextFields-2   	 2649936	       453.6 ns/op	     256 B/op	       2 allocs/op
BenchmarkLogPassedContextFields-2       	 1736319	       689.7 ns/op	     976 B/op	       6 allocs/op
BenchmarkLogPassedGinContextFields-2    	 1257981	       955.6 ns/op	    1361 B/op	       6 allocs/op

# 4 cores
BenchmarkLogEmpty-4                     	 5674430	       196.0 ns/op	     256 B/op	       2 allocs/op
BenchmarkLog-4                          	 4684784	       255.5 ns/op	     256 B/op	       2 allocs/op
BenchmarkLog4Fields-4                   	 4941025	       242.9 ns/op	     256 B/op	       2 allocs/op
BenchmarkLog10Fields-4                  	 4140194	       289.7 ns/op	     256 B/op	       2 allocs/op
BenchmarkLog10FieldsComplex-4           	 3581312	       379.0 ns/op	     280 B/op	       3 allocs/op
BenchmarkLog10FieldsComplexAny-4        	  798502	      1471 ns/op	    1124 B/op	      23 allocs/op
BenchmarkLogPredefinedContextFields-4   	 5403194	       221.4 ns/op	     256 B/op	       2 allocs/op
BenchmarkLogPassedContextFields-4       	 3036115	       395.2 ns/op	     978 B/op	       6 allocs/op
BenchmarkLogPassedGinContextFields-4    	 2135607	       561.6 ns/op	    1363 B/op	       6 allocs/op

# 6 cores
BenchmarkLogEmpty-6                     	 7051023	       160.5 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog-6                          	 5981791	       201.9 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog4Fields-6                   	 6301800	       198.6 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog10Fields-6                  	 5392446	       223.0 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog10FieldsComplex-6           	 4270286	       246.2 ns/op	     281 B/op	       3 allocs/op
BenchmarkLog10FieldsComplexAny-6        	 1000000	      1093 ns/op	    1128 B/op	      23 allocs/op
BenchmarkLogPredefinedContextFields-6   	 6509941	       176.9 ns/op	     257 B/op	       2 allocs/op
BenchmarkLogPassedContextFields-6       	 3420057	       353.7 ns/op	     980 B/op	       6 allocs/op
BenchmarkLogPassedGinContextFields-6    	 2341536	       502.4 ns/op	    1366 B/op	       6 allocs/op

# 8 cores
BenchmarkLogEmpty-8                     	 6564950	       171.6 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog-8                          	 5764022	       207.3 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog4Fields-8                   	 5854122	       201.6 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog10Fields-8                  	 5214798	       233.6 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog10FieldsComplex-8           	 4589498	       263.1 ns/op	     281 B/op	       3 allocs/op
BenchmarkLog10FieldsComplexAny-8        	 1000000	      1090 ns/op	    1133 B/op	      23 allocs/op
BenchmarkLogPredefinedContextFields-8   	 6285962	       189.6 ns/op	     257 B/op	       2 allocs/op
BenchmarkLogPassedContextFields-8       	 3472855	       343.9 ns/op	     980 B/op	       6 allocs/op
BenchmarkLogPassedGinContextFields-8    	 2423517	       491.9 ns/op	    1366 B/op	       6 allocs/op

# 10 cores
BenchmarkLogEmpty-10                      	 6244230	       191.7 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog-10                           	 5582205	       213.8 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog4Fields-10                    	 5682394	       214.8 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog10Fields-10                   	 4737988	       246.5 ns/op	     258 B/op	       2 allocs/op
BenchmarkLog10FieldsComplex-10            	 4265218	       267.5 ns/op	     282 B/op	       3 allocs/op
BenchmarkLog10FieldsComplexAny-10         	 1000000	      1088 ns/op	    1136 B/op	      23 allocs/op
BenchmarkLogPredefinedContextFields-10    	 6228145	       195.6 ns/op	     257 B/op	       2 allocs/op
BenchmarkLogPassedContextFields-10        	 3502688	       338.5 ns/op	     981 B/op	       6 allocs/op
BenchmarkLogPassedGinContextFields-10     	 2381913	       501.8 ns/op	    1367 B/op	       6 allocs/op

# 12 cores
BenchmarkLogEmpty-12                      	 6676146	       170.5 ns/op	     257 B/op	       2 allocs/op
BenchmarkLog-12                           	 5818592	       206.7 ns/op	     258 B/op	       2 allocs/op
BenchmarkLog4Fields-12                    	 5759398	       206.3 ns/op	     258 B/op	       2 allocs/op
BenchmarkLog10Fields-12                   	 4973750	       239.9 ns/op	     258 B/op	       2 allocs/op
BenchmarkLog10FieldsComplex-12            	 4503340	       268.2 ns/op	     282 B/op	       3 allocs/op
BenchmarkLog10FieldsComplexAny-12         	 1000000	      1053 ns/op	    1139 B/op	      23 allocs/op
BenchmarkLogPredefinedContextFields-12    	 6240646	       190.2 ns/op	     258 B/op	       2 allocs/op
BenchmarkLogPassedContextFields-12        	 3544420	       373.8 ns/op	     982 B/op	       6 allocs/op
BenchmarkLogPassedGinContextFields-12     	 2345589	       509.7 ns/op	    1369 B/op	       6 allocs/op

```

