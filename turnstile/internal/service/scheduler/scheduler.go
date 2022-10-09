package scheduler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"time"
	"turnstile/internal/service"
	"turnstile/pkg/logging"
)

type Scheduler struct {
	logger    logging.Logger
	gripper   service.DataGripper
	logSender service.LogSender

	unscheduledUpdateChan chan bool
	unscheduledLogChan    chan bool

	dailyUpdatingTime time.Duration
	logTime           time.Duration
	retryTime         time.Duration
}

func New(logger logging.Logger, dg service.DataGripper, logSender service.LogSender, dailyUpdatingTime time.Duration, logTime time.Duration, retryTime time.Duration) *Scheduler {
	return &Scheduler{
		logger:                logger,
		gripper:               dg,
		logSender:             logSender,
		unscheduledUpdateChan: make(chan bool, 1),
		unscheduledLogChan:    make(chan bool, 1),
		dailyUpdatingTime:     dailyUpdatingTime,
		logTime:               logTime,
		retryTime:             retryTime,
	}
}

func (s *Scheduler) Start() error {
	// First Load
	start := time.Now()
	err := s.gripper.LoadData()
	duration := time.Since(start)
	fmt.Println(">>>", duration)

	if err != nil {
		return err
	}
	go s.updateTicker(s.unscheduledUpdateChan)
	go s.logTicker(s.unscheduledLogChan)
	return nil
}

func (s *Scheduler) UnscheduledUpdate() {
	if len(s.unscheduledUpdateChan) == 1 {
		return
	}
	s.unscheduledUpdateChan <- true
}

func (s *Scheduler) UnscheduledSendLogs() {
	if len(s.unscheduledLogChan) == 1 {
		return
	}
	s.unscheduledLogChan <- true
}

func (s *Scheduler) updateTicker(unscheduledChan <-chan bool) {
	dailyLoadingTicker := time.Tick(s.dailyUpdatingTime)
	retryTicker := time.Tick(s.retryTime)
	retry := false
	for {
		select {
		case <-dailyLoadingTicker:
			s.logger.Info("Trying DailyUpdate...")
			s.update(&retry)
		case <-unscheduledChan:
			s.logger.Info("Trying UnscheduledUpdate...")
			s.update(&retry)
		case <-retryTicker:
			if retry {
				s.logger.Warn("Trying RetryUpdate...")
				s.update(&retry)
			}
		}
	}
}

func (s *Scheduler) logTicker(logChan <-chan bool) {
	logTicker := time.Tick(s.logTime)
	retryTicker := time.Tick(s.retryTime)
	retry := false
	for {
		select {
		case <-logTicker:
			s.logger.Info("Trying Daily LogSending...")
			s.sendLogs(&retry)
		case <-logChan:
			s.logger.Info("Trying Unscheduled LogSending...")
			s.sendLogs(&retry)
		case <-retryTicker:
			if retry {
				s.logger.Warn("Trying Retry LogSending...")
				s.sendLogs(&retry)
			}
		}
	}
}

func (s *Scheduler) update(retry *bool) {
	if isOnline() {
		*retry = false

		err := s.gripper.LoadData()
		if err != nil {
			s.logger.Error("Error in Update: ", err)
			*retry = true
			return
		}
		s.logger.Info("Successfully update")
	} else {
		*retry = true
	}
}

func (s *Scheduler) sendLogs(retry *bool) {
	if isOnline() {
		*retry = false

		err := s.logSender.SendLogs()
		if err != nil {
			s.logger.Error("Error in LogSending: ", err)
			*retry = true
			return
		}
		s.logger.Info("Logs have been sent")
	} else {
		*retry = true
	}
}

func isOnline() bool {
	const (
		protocol = "udp"
		dns1     = "8.8.8.8:80"      // google
		dns2     = "5.255.255.70:80" // yandex
		dns3     = "1.1.1.1:80"      // cloudflare
	)
	conns := []string{dns1, dns2, dns3}
	for i := range conns {
		conn, err := net.Dial(protocol, conns[i])
		if err != nil {
			logrus.Debugf("conn to %s failed", conns[i])
			continue
		}
		_ = conn.Close()
		return true
	}
	return false
}
