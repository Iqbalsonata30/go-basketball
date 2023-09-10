package helper

func CheckError(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func WriteMessageAPI(code int, message string) map[string]any {
	return map[string]any{
		"message":    message,
		"statusCode": code,
	}
}
