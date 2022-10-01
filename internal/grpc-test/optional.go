package grpc_test

func Bool(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}
