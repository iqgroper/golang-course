package main

import (
	context "context"
	"strings"
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
	LogListMu        *sync.RWMutex
	Ctx              context.Context
}

func NewAdmin(ctx context.Context, acl map[string][]string, logs chan *Event) *Admin {
	newAdmin := &Admin{
		Logs:             logs,
		Stats:            make(chan *RawStat, 2),
		ACL:              acl,
		EventClientList:  make([]chan *Event, 0),
		StatsChan:        make(chan *Stat, 2),
		StatByMethod:     make([]map[string]uint64, 1),
		StatByConsumer:   make([]map[string]uint64, 1),
		StatMu:           &sync.RWMutex{},
		LogListMu:        &sync.RWMutex{},
		LoggingClientCnt: 0,
		Ctx:              ctx,
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
		ctx := admin.Ctx
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-admin.Logs:
				stat := &RawStat{
					Method:   event.Method,
					Consumer: event.Consumer,
				}
				// стоит ли делать из этого горутину, ведь он может заблочится на попытке записать?
				if admin.StatClientCnt != 0 {
					admin.Stats <- stat
				}

				logNumber := admin.LoggingClientCnt

				for i := 0; i < logNumber; i++ {
					admin.EventClientList[i] <- event
				}
			}
		}

	}()
}

func (admin *Admin) sendNewLoggingClientLogs(ctx context.Context, clientID int) {

	md, _ := metadata.FromIncomingContext(ctx)

	newChan := make(chan *Event, 1)
	admin.EventClientList = append(admin.EventClientList, newChan)

	// host := strings.Split(md[":authority"][0], ":")[0]

	event := &Event{
		Timestamp: time.Now().Unix(),
		Consumer:  md["consumer"][0],
		Method:    "/main.Admin/Logging",
		Host:      md[":authority"][0][:strings.IndexByte(md[":authority"][0], ':')+1],
	}
	for i := 0; i < clientID; i++ {
		admin.EventClientList[i] <- event
	}
}

func (admin *Admin) Logging(nothing *Nothing, outStream Admin_LoggingServer) error {

	ctx := outStream.Context()

	admin.LogListMu.Lock()
	clientID := admin.LoggingClientCnt
	admin.LoggingClientCnt++
	admin.sendNewLoggingClientLogs(ctx, clientID)
	admin.LogListMu.Unlock()

	for event := range admin.EventClientList[clientID] {
		outStream.Send(event)
	}

	return nil
}

func (admin *Admin) fomringStats() {
	go func() {
		ctx := admin.Ctx
		for {
			select {
			case <-ctx.Done():
				return
			case stat := <-admin.Stats:
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
		}

	}()
}

func (admin *Admin) sendNewStatClientLogs(ctx context.Context) {

	md, _ := metadata.FromIncomingContext(ctx)

	admin.StatMu.Lock()

	if admin.StatClientCnt == 0 {
		admin.StatByMethod[0] = map[string]uint64{}
		admin.StatByConsumer[0] = map[string]uint64{}
		return
	}

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

	admin.StatMu.Lock()
	admin.StatClientCnt++
	admin.StatMu.Unlock()

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
