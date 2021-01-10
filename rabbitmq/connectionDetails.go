package rabbitmq

type connectionDetails struct {
	hostname string
	port     int
	username string
	password string
}

func ConnectionDetails(hostname string, port int, username string, password string) *connectionDetails {
	connDetails := connectionDetails{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
	}
	// Todo some more checking

	return &connDetails
}