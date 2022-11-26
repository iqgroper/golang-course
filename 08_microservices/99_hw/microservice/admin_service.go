package main

type Admin struct {
	UnimplementedAdminServer
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
