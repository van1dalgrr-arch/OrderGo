package main

func validStatus(status string) bool {
	switch status {
	case "pending", "paid", "processing", "shipped", "completed", "cancelled":
		return true
	default:
		return false
	}
}
