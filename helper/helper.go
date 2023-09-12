package helper

import "fmt"

func WriteMessageAPI(code int, message string) map[string]any {
	return map[string]any{
		"message":    message,
		"statusCode": code,
	}
}

func Required(datas ...any) error {
	for _, data := range datas {
		if data == "" {
			return fmt.Errorf("you've not fulfilled the required data")
		}
	}
	return nil
}
