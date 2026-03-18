package remote

import "time"

const SSHTimeout = 2 * time.Second

func TimeoutDuration(seconds int) time.Duration {
	if seconds <= 0 {
		return SSHTimeout
	}
	return time.Duration(seconds) * time.Second
}
