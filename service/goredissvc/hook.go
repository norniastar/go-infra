package goredissvc

//type redisTracingHook struct {
//	tracer opentracing.Tracer
//}
//
//func (hook redisTracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
//	span := opentracing.SpanFromContext(ctx)
//	defer span.Finish()
//
//	if err := cmd.Err(); err != nil {
//		recordError("db.error", span, err)
//	}
//	return nil
//}
//
//func (hook redisTracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
//	span := opentracing.SpanFromContext(ctx)
//	defer span.Finish()
//
//	for i, cmd := range cmds {
//		if err := cmd.Err(); err != nil {
//			recordError("db.error"+strconv.Itoa(i), span, err)
//		}
//	}
//	return nil
//}
//
//func (hook redisTracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
//	span, ctx := hook.createSpan(ctx, cmd.FullName())
//	span.SetTag("db.type", "redis")
//	return ctx, nil
//}
//
//func (hook redisTracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
//	span, ctx := hook.createSpan(ctx, "pipeline")
//	span.SetTag("db.type", "redis")
//	span.SetTag("db.redis.num_cmd", len(cmds))
//	return ctx, nil
//}
//
//func (hook redisTracingHook) createSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
//	span := opentracing.SpanFromContext(ctx)
//	if span != nil {
//		childSpan := hook.tracer.StartSpan(operationName, opentracing.ChildOf(span.Context()))
//		return childSpan, opentracing.ContextWithSpan(ctx, childSpan)
//	}
//
//	return opentracing.StartSpanFromContextWithTracer(ctx, hook.tracer, operationName)
//}
//
//func recordError(errorTag string, span opentracing.Span, err error) {
//	if err != redis.Nil {
//		span.SetTag(string(ext.Error), true)
//		span.SetTag(errorTag, err.Error())
//	}
//}
//
//type iHook interface {
//	AddHook(hook redis.Hook)
//}
//
//// NewHookOption is hook选项
//func NewHookOption(tracer opentracing.Tracer) contract.RedisOption {
//	return func(m contract.IRedis) {
//		m.(*redisAdapter).getClient().(iHook).AddHook(&redisTracingHook{
//			tracer: tracer,
//		})
//	}
//}
