package loadbalancer

// Define basic behavior for a single loadBalancer
// May need adapt to different types: HAProxy/Nginx/LVS...
type LoadBalancer struct {
	LBtype string
}