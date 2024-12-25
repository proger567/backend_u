package internal

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/sirupsen/logrus"
	"testgenerate_backend_user/internal/app"
	"time"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next   Service
	logger *logrus.Logger
}

func LoggingMiddleware(logger *logrus.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (mw loggingMiddleware) GetUser(ctx context.Context, userName, userRole string) (user app.User, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == GetUser")
	}(time.Now())
	return mw.next.GetUser(ctx, userName, userRole)
}

func (mw loggingMiddleware) GetUsersRole(ctx context.Context) (users []app.User, err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == GetUsersRole")
	}(time.Now())
	return mw.next.GetUsersRole(ctx)
}

func (mw loggingMiddleware) AddUser(ctx context.Context, userAdd app.User) (err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == AddUser")
	}(time.Now())
	return mw.next.AddUser(ctx, userAdd)
}

func (mw loggingMiddleware) UpdateUser(ctx context.Context, user app.User) (err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == UpdateUser")
	}(time.Now())
	return mw.next.UpdateUser(ctx, user)
}

func (mw loggingMiddleware) DeleteUser(ctx context.Context, userName string) (err error) {
	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{
			"took":  time.Since(begin).Milliseconds(),
			"error": err,
		}).Info("method == DeleteUser")
	}(time.Now())
	return mw.next.DeleteUser(ctx, userName)
}

// ----------------------------------------------------------------------------------------------------------------------
type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func InstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram) Middleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			requestCount:   requestCount,
			requestLatency: requestLatency,
			next:           next,
		}
	}
}

func (im instrumentingMiddleware) GetUser(ctx context.Context, userName, userRole string) (user app.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getUser", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	user, err = im.next.GetUser(ctx, userName, userRole)
	return
}

func (im instrumentingMiddleware) GetUsersRole(ctx context.Context) (users []app.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getUserRole", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	users, err = im.next.GetUsersRole(ctx)
	return
}

func (im instrumentingMiddleware) AddUser(ctx context.Context, userAdd app.User) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "addUser", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = im.next.AddUser(ctx, userAdd)
	return
}

func (im instrumentingMiddleware) UpdateUser(ctx context.Context, user app.User) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "updateUser", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = im.next.UpdateUser(ctx, user)
	return
}

func (im instrumentingMiddleware) DeleteUser(ctx context.Context, userName string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "deleteUser", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = im.next.DeleteUser(ctx, userName)
	return
}
