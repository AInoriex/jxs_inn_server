package handler

import (
	"eshop_server/src/common/cache"
	router_dao "eshop_server/src/router/dao"
	router_handler "eshop_server/src/router/handler"
	router_model "eshop_server/src/router/model"
	"eshop_server/src/utils/log"
	"time"

	"gorm.io/gorm"
)

// @Title		定时任务处理支付并更新订单状态
// @Description	轮询获取支付中的订单，查询YLT订单状态，更新订单状态
// @Attention	轮询间隔时间为5秒，该func所有操作需要5s内全部完成
// @Attention	支付超时时间为3分钟，超过5分钟未支付则更新订单状态为超时，订单状态为超时后不再查询YLT订单状态
func UpdateOrderCronjob() {
	var err error
	var paymentTimeOutLimitMins int32 = 3 // 超时限制:3分钟
	// 查询是否有支付中的订单
	paymentList, err := router_dao.GetPaymentsByStatus(router_model.PaymentStatusPaying)
	if err != nil && gorm.ErrRecordNotFound == err {
		log.Infof("UpdateOrderCronjob 无支付中订单")
		return
	} else if err != nil {
		log.Errorf("UpdateOrderCronjob 查询支付中订单失败, error:%s", err.Error())
		return
	}
	log.Infof("UpdateOrderCronjob 查询到支付中订单数量为:%v", len(paymentList))
	for _, payment := range paymentList {
		// 判断支付超时以及超时处理
		if IsPaymentTimeout(payment, paymentTimeOutLimitMins) {
			log.Errorf("UpdateOrderCronjob 订单支付超时, paymentId:%s, orderId:%s, createAt:%s", payment.Id, payment.OrderId, payment.CreatedAt.Format("2006-01-02 15:04:05"))
			go PaymentTimeoutHandler(payment)
			continue
		}

		// 获取YLT Token信息
		flag, gt_token, cookie := cache.GetYltUserToken(payment.Agent)
		if !flag {
			log.Errorf("UpdateOrderCronjob 获取YLT登陆Token缓存信息失败, agent:%s", payment.Agent)
			return
		}
		// 查询YLT支付状态，更新平台订单状态
		go CronYltPaymentApiToUpdateOrder(gt_token, cookie, payment)

		time.Sleep(200 * time.Millisecond)
	}
}

// @Title		判断支付是否超时
func IsPaymentTimeout(payment *router_model.Payment, timeoutMinsLimit int32) bool {
	now := time.Now()
	// 支付超时时间
	timeoutMins := payment.CreatedAt.Add(time.Duration(timeoutMinsLimit) * time.Minute)
	// 判断是否超时
	return now.After(timeoutMins)
}

// @Title		支付超时处理
func PaymentTimeoutHandler(payment *router_model.Payment) {
	var err error
	// 订单支付超时，更新支付状态
	payment.Status = router_model.PaymentStatusTimeOut
	if _, err = router_dao.UpdatePaymentByField(payment, []string{"status"}); err != nil {
		log.Errorf("PaymentTimeoutHandler 更新平台`支付状态为超时`失败, paymentId:%s, orderId:%s, gatewayId:%s, error:%s", payment.Id, payment.OrderId, payment.GatewayID, err.Error())
		// TODO 更新数据失败告警
	}
	order := &router_model.Order{
		Id:            payment.OrderId,
		PaymentStatus: router_model.OrderPaymentStatusTimeOut,
	}
	if _, err = router_dao.UpdateOrderByField(order, []string{"payment_status"}); err != nil {
		log.Errorf("PaymentTimeoutHandler 更新平台`订单支付状态为超时`失败, orderId:%s, error:%s", payment.OrderId, err.Error())
		// TODO 更新数据失败告警
	}
}

// 查询YLT付费接口付费状态，付费成功更新平台订单状态
func CronYltPaymentApiToUpdateOrder(gt_token string, cookie string, payment *router_model.Payment) {
	// 查询YLT订单状态
	payOk, err := router_handler.YltCheckOrder(gt_token, cookie, payment.GatewayID)
	if err != nil {
		log.Errorf("CronYltPaymentApiToUpdateOrder 查询YLT订单失败, gt_token:%s, cookie:%s, gatewayId:%s, error:%s", gt_token, cookie, payment.GatewayID, err.Error())
		return
	}
	if !payOk {
		log.Errorf("CronYltPaymentApiToUpdateOrder 用户未完成支付，等待下一轮查询... gt_token:%s, cookie:%s, gatewayId:%s", gt_token, cookie, payment.GatewayID)
		return
	}

	// 订单支付成功，更新支付状态
	payment.Status = router_model.PaymentStatusPayed
	payment.PurchasedAt = time.Now()
	if _, err = router_dao.UpdatePaymentByField(payment, []string{"status", "purchased_at"}); err != nil {
		log.Errorf("CronYltPaymentApiToUpdateOrder 更新平台支付状态失败, error:%s", err.Error())
		// TODO 更新数据库失败告警
		return
	}

	// 事务操作：如果数据更新失败，需要设置订单状态为支付失败（是否需要重试？）
	// TODO 新增本地消息表任务，定时重试
	func() {
		log.Infof("CronYltPaymentApiToUpdateOrder 查询到YLT订单已支付，开始更新数据库状态，paymentId:%s , orderId:%s", payment.Id, payment.OrderId)

		// 1. 订单支付成功，更新订单状态
		order, err := router_dao.GetOrderByOrderId(payment.OrderId)
		if err != nil {
			log.Errorf("CronYltPaymentApiToUpdateOrder 查询订单信息失败, error:%s", err.Error())
			// TODO 更新数据库失败告警
			return
		}
		order.PaymentStatus = router_model.OrderPaymentStatusPayed
		order, err = router_dao.UpdateOrderByField(order, []string{"payment_status"})
		if err != nil {
			log.Errorf("CronYltPaymentApiToUpdateOrder 更新平台订单状态失败, error:%s", err.Error())
			// TODO 更新数据库失败告警
			return
		}

		// 2. 获取用户下单商品信息
		items, err := router_dao.GetOrderItemsById(order.ItemId)
		if err != nil {
			log.Errorf("CronYltPaymentApiToUpdateOrder 查询订单商品信息失败, error:%s", err.Error())
			// TODO 更新数据库失败告警
			return
		}

		// 3. 创建用户商品购买记录
		for _, item := range items {
			// UserId-ProductId-PaymentId 唯一索引，防止重复创建
			_, err = router_dao.GetPurchaseHistoryByUserIdAndProductIdAndPaymentId(order.UserId, item.ProductId, payment.Id)
			if err != nil && gorm.ErrRecordNotFound == err {
				ph := &router_model.PurchaseHistory{
					UserId:      order.UserId,
					ProductId:   item.ProductId,
					Quantity:    item.Quantity,
					OrderId:     order.Id,
					PaymentId:   payment.Id,
					PurchasedAt: payment.PurchasedAt,
				}
				if _, err = router_dao.CreatePurchaseHistory(ph); err != nil {
					log.Errorf("CronYltPaymentApiToUpdateOrder 创建用户购买记录失败, error:%s", err.Error())
					// TODO 操作数据失败告警
				}
			} else if err != nil {
				log.Errorf("CronYltPaymentApiToUpdateOrder 查询用户购买记录失败, error:%s", err.Error())
				// TODO 操作数据失败告警
			} else {
				log.Warnf("CronYltPaymentApiToUpdateOrder 用户购买记录已存在，跳过创建，paymentId:%s, orderId:%s, productId:%s", payment.Id, payment.OrderId, item.ProductId)
			}
		}

		// 5. 订单支付成功，更新商品销量
		// TODO: 放入每日定时任务更新该数据

		log.Infof("CronYltPaymentApiToUpdateOrder 更新平台订单成功，处理完毕，paymentId:%s , orderId:%s", payment.Id, payment.OrderId)
	}()
}
