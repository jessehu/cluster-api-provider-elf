package util

func TowerBool(v bool) *bool {
	return &v
}

func TowerFloat64(v int) *float64 {
	val := float64(v)

	return &val
}

func TowerString(v string) *string {
	return &v
}

func TowerCPU(v int32) float64 {
	return float64(v)
}

func TowerMemory(memoryMiB int64) float64 {
	memory := float64(memoryMiB)
	memory = memory * 1024 * 1024

	return memory
}

func TowerDisk(diskGiB int32) *float64 {
	disk := float64(diskGiB)
	disk = disk * 1024 * 1024 * 1024

	return &disk
}
