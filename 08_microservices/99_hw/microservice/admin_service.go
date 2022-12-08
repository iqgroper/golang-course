package main

import (
	context "context"
	sync "sync"
	"time"

	"google.golang.org/grpc/metadata"
)

type Admin struct {
	UnimplementedAdminServer
	ACL              map[string][]string
	Logs             chan *Event
	Stats            chan *RawStat
	EventClientList  []chan *Event
	StatsChan        chan *Stat
	LoggingClientCnt int
	StatClientCnt    int
	StatByMethod     []map[string]uint64
	StatByConsumer   []map[string]uint64
	StatMu           *sync.RWMutex
}

func NewAdmin(acl map[string][]string, logs chan *Event) *Admin {
	newAdmin := &Admin{
		Logs:             logs,
		Stats:            make(chan *RawStat, 2),
		ACL:              acl,
		EventClientList:  make([]chan *Event, 0),
		StatsChan:        make(chan *Stat, 2),
		StatByMethod:     make([]map[string]uint64, 1),
		StatByConsumer:   make([]map[string]uint64, 1),
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

			stat := &RawStat{
				Method:   event.Method,
				Consumer: event.Consumer,
			}
			// стоит ли делать из этого горутину, ведь он может заблочится на попытке записать?
			admin.Stats <- stat

			logNumber := admin.LoggingClientCnt

			for i := 0; i < logNumber; i++ {
				admin.EventClientList[i] <- event
			}
		}
	}()
}

func (admin *Admin) sendNewLoggingClientLogs(ctx context.Context, clientID int) {

	md, _ := metadata.FromIncomingContext(ctx)

	newChan := make(chan *Event, 1)
	admin.EventClientList = append(admin.EventClientList, newChan)

	event := &Event{
		Timestamp: time.Now().Unix(),
		Consumer:  md["consumer"][0],
		Method:    "/main.Admin/Logging",
		Host:      "127.0.0.1:",
	}
	for i := 0; i < clientID; i++ {
		admin.EventClientList[i] <- event
	}
}

func (admin *Admin) Logging(nothing *Nothing, outStream Admin_LoggingServer) error {

	clientID := admin.LoggingClientCnt
	admin.LoggingClientCnt++

	ctx := outStream.Context()
	admin.sendNewLoggingClientLogs(ctx, clientID)

	for event := range admin.EventClientList[clientID] {
		outStream.Send(event)
	}

	return nil
}

func (admin *Admin) fomringStats() {
	go func() {
		for stat := range admin.Stats {

			if admin.StatClientCnt == 0 {
				continue
			}

			admin.StatMu.Lock()

			for _, statMap := range admin.StatByMethod {
				statMap[stat.Method]++
			}

			for _, statMap := range admin.StatByConsumer {
				statMap[stat.Consumer]++
			}

			admin.StatMu.Unlock()
		}
	}()
}

func (admin *Admin) sendNewStatClientLogs(ctx context.Context) {

	md, _ := metadata.FromIncomingContext(ctx)

	if admin.StatClientCnt == 0 {
		admin.StatByMethod[0] = map[string]uint64{}
		admin.StatByConsumer[0] = map[string]uint64{}
		return
	}

	admin.StatMu.Lock()

	for _, statMap := range admin.StatByMethod {
		statMap["/main.Admin/Statistics"]++
	}

	for _, statMap := range admin.StatByConsumer {
		statMap[md["consumer"][0]]++
	}

	admin.StatByMethod = append(admin.StatByMethod, map[string]uint64{})
	admin.StatByConsumer = append(admin.StatByConsumer, map[string]uint64{})

	admin.StatMu.Unlock()
}

func (admin *Admin) Statistics(interval *StatInterval, outStream Admin_StatisticsServer) error {

	ctx := outStream.Context()
	admin.sendNewStatClientLogs(ctx)

	statClientId := admin.StatClientCnt
	admin.StatClientCnt++

	ticker := time.NewTicker(time.Duration(interval.IntervalSeconds) * time.Second)
	for range ticker.C {
		newStat := &Stat{
			ByMethod:   admin.StatByMethod[statClientId],
			ByConsumer: admin.StatByConsumer[statClientId],
		}
		outStream.Send(newStat)

		for key := range admin.StatByMethod[statClientId] {
			delete(admin.StatByMethod[statClientId], key)
		}
		for key := range admin.StatByConsumer[statClientId] {
			delete(admin.StatByConsumer[statClientId], key)
		}
	}

	return nil
}
