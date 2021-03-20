package ray

type ServerOption func(s *Server)

func WithServerIp(ip string) ServerOption {
	return func(s *Server) {
		s.ip = ip
	}
}

func WithServerPort(port int) ServerOption {
	return func(s *Server) {
		s.port = port
	}
}

func WithServerLogger(logger Logger) ServerOption {
	return func(s *Server) {
		s.logger = logger
	}
}
