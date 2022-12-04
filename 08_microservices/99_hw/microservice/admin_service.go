package main

import (
	context "context"
	"log"
	sync "sync"
	"time"

	"google.golang.org/grpc/metadata"
)

type Admin struct {
	UnimplementedAdminServer
	ACL              map[string][]string
	Logs             chan *Event
	Stats            chan *RawStat
	EventsChan       chan *Event
	StatsChan        chan *Stat
	LoggingClientCnt int
	StatClientCnt    int
	StatByMethod     map[string]uint64
	StatByConsumer   map[string]uint64
	StatMu           *sync.RWMutex
}

func NewAdmin(acl map[string][]string, logs chan *Event) *Admin {
	newAdmin := &Admin{
		Logs:             logs,
		Stats:            make(chan *RawStat, 2),
		ACL:              acl,
		EventsChan:       make(chan *Event, 2),
		StatsChan:        make(chan *Stat, 2),
		StatByMethod:     make(map[string]uint64),
		StatByConsumer:   make(map[string]uint64),
		StatMu:           &sync.RWMutex{},
		LoggingClientCnt: 0,
	}

	newAdmin.fomringStats()
	newAdmin.gettingLogsAndStats()

	return newAdmin
}

type RawStat struct {
	Method   string
	Consumer string
}

func (admin *Admin) gettingLogsAndStats() {
	go func() {
		for event := range admin.Logs {

			statNumber := admin.StatClientCnt

			stat := &RawStat{
				Method:   event.Method,
				Consumer: event.Consumer,
			}
			for i := 0; i < statNumber; i++ {
				admin.Stats <- stat
			}

			logNumber := admin.LoggingClientCnt

			log.Println("CLIENT NUM", admin.LoggingClientCnt)

			for i := 0; i < logNumber; i++ {
				admin.EventsChan <- event
			}
		}
	}()
}

func (admin *Admin) sendNewLoggingClientLogs(ctx context.Context) {

	md, _ := metadata.FromIncomingContext(ctx)

	event := Event{
		Timestamp: time.Now().Unix(),
		Consumer:  md["consumer"][0],
		Method:    "/main.Admin/Logging",
		Host:      "127.0.0.1:",
	}
	for i := 0; i < admin.LoggingClientCnt; i++ {
		admin.EventsChan <- &event
	}
}

func (admin *Admin) Logging(nothing *Nothing, outStream Admin_LoggingServer) error {

	ctx := outStream.Context()
	admin.sendNewLoggingClientLogs(ctx)

	admin.LoggingClientCnt++
	clientId := admin.LoggingClientCnt

	for event := range admin.EventsChan {
		log.Printf(" FOR CLIENT NUM %d SENDING EVENT %v", clientId, event)
		outStream.Send(event)
	}

	return nil
}

func (admin *Admin) sendNewStatClientLogs(ctx context.Context) {

	md, _ := metadata.FromIncomingContext(ctx)

	admin.StatMu.Lock()

	admin.StatByMethod["/main.Admin/Statistics"]++
	admin.StatByConsumer[md["consumer"][0]]++

	log.Println("BY CONS", admin.StatByConsumer)

	stat := &Stat{
		ByMethod:   admin.StatByMethod,
		ByConsumer: admin.StatByConsumer,
	}

	admin.StatMu.Unlock()

	for i := 0; i < admin.StatClientCnt; i++ {
		admin.StatsChan <- stat
	}
}

func (admin *Admin) fomringStats() {
	go func() {
		for stat := range admin.Stats {
			admin.StatMu.Lock()

			admin.StatByMethod[stat.Method]++
			admin.StatByConsumer[stat.Consumer]++
			statFinal := &Stat{
				ByMethod:   admin.StatByMethod,
				ByConsumer: admin.StatByConsumer,
			}

			admin.StatMu.Unlock()
			admin.StatsChan <- statFinal
		}
	}()
}

func (admin *Admin) Statistics(interval *StatInterval, outStream Admin_StatisticsServer) error {

	ctx := outStream.Context()
	admin.sendNewStatClientLogs(ctx)

	admin.StatClientCnt++

	ticker := time.NewTicker(time.Duration(interval.IntervalSeconds))
	for range ticker.C {
		stat := <-admin.StatsChan
		outStream.Send(stat)
	}

	return nil
}
