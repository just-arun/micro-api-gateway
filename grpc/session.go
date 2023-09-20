package grpc

type session struct {}

func Session() session {
	return session{}	
}

func (st session) VerifySession() {
	
}
