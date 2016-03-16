package trilateration

import "github.com/octoblu/find-yourself/scanner"

// Swarm is a swarm of devices
type Swarm struct {
	devices                   []*scanner.Device
	onLocationUpdateCallbacks []OnLocationUpdateCallback
}

// OnLocationUpdateCallback callbacks are called whenever
// the location updates
type OnLocationUpdateCallback func()

// NewSwarm constructs a Swarm instance
func NewSwarm() *Swarm {
	var onLocationUpdateCallbacks []OnLocationUpdateCallback
	var devices []*scanner.Device
	return &Swarm{devices, onLocationUpdateCallbacks}
}

// Accuracies returns an array of ints
func (swarm *Swarm) Accuracies() []float64 {
	strenths := make([]float64, swarm.DeviceCount())
	for i, device := range swarm.devices {
		strenths[i] = device.Accuracy()
	}
	return strenths
}

// AddDevice adds a device to the swarm
func (swarm *Swarm) AddDevice(device *scanner.Device) {
	device.OnRSSIUpdate(func(int) {
		swarm.emitLocationUpdate()
	})
	swarm.devices = append(swarm.devices, device)
}

// DeviceCount returns the number of devices in the swarm
func (swarm *Swarm) DeviceCount() int {
	return len(swarm.devices)
}

// Distances returns an array of ints
func (swarm *Swarm) Distances() []int {
	strenths := make([]int, swarm.DeviceCount())
	for i, device := range swarm.devices {
		strenths[i] = device.Distance()
	}
	return strenths
}

// OnLocationUpdate registers an OnLocationUpdateCallback listener
// with the swarm
func (swarm *Swarm) OnLocationUpdate(callback OnLocationUpdateCallback) {
	swarm.onLocationUpdateCallbacks = append(swarm.onLocationUpdateCallbacks, callback)
}

// SignalStrenths returns an array of ints
func (swarm *Swarm) SignalStrenths() []int {
	strenths := make([]int, swarm.DeviceCount())
	for i, device := range swarm.devices {
		strenths[i] = device.RSSI()
	}
	return strenths
}

func (swarm *Swarm) emitLocationUpdate() {
	for _, callback := range swarm.onLocationUpdateCallbacks {
		callback()
	}
}
