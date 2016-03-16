package scanner

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/paypal/gatt"
)

// Device stores all BLE device related information collected
type Device struct {
	advertisement         *gatt.Advertisement
	peripheral            gatt.Peripheral
	rssis                 []int
	onRSSIUpdateCallbacks []OnRSSIUpdateCallback
}

// OnRSSIUpdateCallback gets called with rssi values
type OnRSSIUpdateCallback func(rssi int)

// NewDevice constructs a new Device
func NewDevice(advertisement *gatt.Advertisement, peripheral gatt.Peripheral, rssi int) *Device {
	var onRSSIUpdateCallbacks []OnRSSIUpdateCallback
	rssis := []int{rssi}
	return &Device{advertisement, peripheral, rssis, onRSSIUpdateCallbacks}
}

// AddRSSI logs the RSSI
func (device *Device) AddRSSI(rssi int) {
	device.rssis = append(device.rssis, rssi)
	for _, callback := range device.onRSSIUpdateCallbacks {
		callback(rssi)
	}
}

// Accuracy calculates the distance accuracy
// using a random formula from stack overflow
func (device *Device) Accuracy() float64 {
	rssi := device.RSSI()
	txPower := int(device.Tx())

	if rssi == 0 {
		return -1.0 // if we cannot determine accuracy, return -1.
	}

	ratio := float64(rssi) / float64(txPower)
	if ratio < 1.0 {
		return math.Pow(ratio, 10)
	}

	return (0.89976)*math.Pow(ratio, 7.7095) + 0.111
}

// Distance calculates the distance to the device
func (device *Device) Distance() int {
	rssi := device.RSSI()
	txPower := int(device.Tx())

	ratioDB := txPower - rssi
	ratio := math.Pow(10, float64(ratioDB)/10)

	return int(math.Sqrt(ratio))
}

// Major returns the major id from the device's advertisement
func (device *Device) Major() uint16 {
	data := device.advertisement.ManufacturerData
	return binary.BigEndian.Uint16(data[20:22])
}

// Minor returns the minor id from the device's advertisement
func (device *Device) Minor() uint16 {
	data := device.advertisement.ManufacturerData
	return binary.BigEndian.Uint16(data[22:24])
}

// OnRSSIUpdate registers an on rssi update listener
func (device *Device) OnRSSIUpdate(callback OnRSSIUpdateCallback) {
	device.onRSSIUpdateCallbacks = append(device.onRSSIUpdateCallbacks, callback)
}

// RSSI returns the most current RSSI value
func (device *Device) RSSI() int {
	return device.rssis[len(device.rssis)-1]
}

func (device *Device) String() string {
	return fmt.Sprintf("%v", device.RSSI())
}

// Tx returns the transmission power of the device
func (device *Device) Tx() uint8 {
	return 1
	// data := device.advertisement.ManufacturerData
	// return data[24]
}
