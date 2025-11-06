package cron

import (
	"time"

	"ecommerce-go-api/domain"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	scheduler            gocron.Scheduler
	paymentExpiryJob     *PaymentExpiryJob
	orderAutoCompleteJob *OrderAutoCompleteJob
}

func NewScheduler(orderRepo domain.OrderRepository, productRepo domain.ProductRepository) (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	paymentExpiryJob := NewPaymentExpiryJob(orderRepo, productRepo)
	orderAutoCompleteJob := NewOrderAutoCompleteJob(orderRepo)

	return &Scheduler{
		scheduler:            s,
		paymentExpiryJob:     paymentExpiryJob,
		orderAutoCompleteJob: orderAutoCompleteJob,
	}, nil
}

func (s *Scheduler) Start() error {

	_, err := s.scheduler.NewJob(
		gocron.DurationJob(10*time.Minute),
		gocron.NewTask(s.paymentExpiryJob.ProcessExpiredPayments),
	)
	if err != nil {
		return err
	}

	_, err = s.scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(s.orderAutoCompleteJob.ProcessDeliveredOrders),
	)
	if err != nil {
		return err
	}

	s.scheduler.Start()

	return nil
}

func (s *Scheduler) Stop() error {
	return s.scheduler.Shutdown()
}
