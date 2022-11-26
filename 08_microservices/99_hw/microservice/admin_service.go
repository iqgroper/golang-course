package main

type Admin struct {
	UnimplementedAdminServer
	// rpc Logging (Nothing) returns (stream Event) {}
	// rpc Statistics (StatInterval) returns (stream Stat) {}
}

func NewAdmin() *Admin {
	return &Admin{}
}

func (admin *Admin) Logging(nothing *Nothing, outStream Admin_LoggingServer) error {
	return nil
}

func (admin *Admin) Statistics(interval *StatInterval, outStream Admin_StatisticsServer) error {
	return nil
}
