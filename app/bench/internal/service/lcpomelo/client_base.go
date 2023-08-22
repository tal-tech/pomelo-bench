package lcpomelo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"pomelo_bench/pomelosdk"
	"sync/atomic"
	"time"
)

// runAndWaitConnect Connector 初始化握手信息和保持连接
func (c *ClientConnector) runAndWaitConnect(ctx context.Context, connector *pomelosdk.Connector, address string, timeout time.Duration, failed func(message string)) error {
	err := connector.InitReqHandshake("0.6.0", "golang-websocket", nil, map[string]interface{}{"uid": "dude"})
	if err != nil {
		return err
	}
	err = connector.InitHandshakeACK(13)
	if err != nil {
		return err
	}

	// 确保连接成功再返回
	cond := syncx.NewCond()

	connector.Connected(func() {

		cond.Signal()
	})

	go func() {

		// 增加超时时间
		ctx2, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = connector.Run(ctx2, address, 10)
		if err != nil {

			logx.WithContext(ctx).Errorf("[%d] pomelo Connector.Run failed ,err:%s", c.uid, err)

			if failed != nil {
				failed(err.Error())
			}

		}
	}()

	_, ok := cond.WaitWithTimeout(timeout + 5*time.Second)
	if !ok {
		return errors.New("run timeout")
	}

	return nil
}

// syncRequest 同步发送消息
func (c *ClientConnector) syncRequest(ctx context.Context, connector *pomelosdk.Connector, timeout time.Duration, reqId *uint64, route string, request interface{}, response interface{}) error {

	cond := syncx.NewCond()

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	var (
		responseStr string
	)

	err = connector.Request(route, requestBytes, func(data []byte) {

		responseStr = string(data)

		cond.Signal()
	})

	// 增加发送序号
	atomic.AddUint64(reqId, 1)

	_, ok := cond.WaitWithTimeout(timeout)
	if !ok {
		logx.WithContext(ctx).Errorf("[%d] %d - %s -- request timeout, request: %s ,response: %s", c.uid, atomic.LoadUint64(reqId), route, string(requestBytes), responseStr)
		return errors.New("请求消息超时")
	}

	if response != nil {
		err = jsonx.UnmarshalFromString(responseStr, response)

		if err != nil {
			logx.WithContext(ctx).Errorf("[%d] %d - %s -- request failed, err: %s , request: %s ,response: %s", c.uid, atomic.LoadUint64(reqId), route, err, string(requestBytes), responseStr)
			return err
		}
	}

	logx.WithContext(ctx).Infof("[%d] %d - %s -- request success, request: %s ,response: %s", c.uid, atomic.LoadUint64(reqId), route, string(requestBytes), responseStr)

	return nil
}

// asyncRequest 异步简单发送消息
func (c *ClientConnector) asyncRequest(ctx context.Context, connector *pomelosdk.Connector, reqId *uint64, route string, sendData []byte, cb pomelosdk.Callback) error {

	id := atomic.LoadUint64(reqId)

	err := connector.Request(route, sendData, func(data []byte) {

		logx.WithContext(ctx).Infof("[%d] %d - %s -- callback success, response.data: %s ", c.uid, id, route, string(data))

		if cb != nil {
			cb(data)
		}

	})

	// 增加发送序号
	atomic.AddUint64(reqId, 1)

	logx.WithContext(ctx).Infof("[%d] %d - %s -- request success, request: %s ", c.uid, id, route, string(sendData))

	return err
}

// notify 同步发送通知
func (c *ClientConnector) notify(ctx context.Context, connector *pomelosdk.Connector, reqId *uint64, route string, request interface{}) error {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	err = connector.Notify(route, requestBytes)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%d] %d - %s -- notify failed, request: %s ", c.uid, atomic.LoadUint64(reqId), route, string(requestBytes))
		return err
	}

	// 增加发送序号
	atomic.AddUint64(reqId, 1)

	logx.WithContext(ctx).Infof("[%d] %d - %s -- notify success, request: %s ", c.uid, atomic.LoadUint64(reqId), route, string(requestBytes))
	return nil
}
