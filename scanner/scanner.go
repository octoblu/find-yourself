package scanner

import (
	"github.com/paypal/gatt"
	BleOptions "github.com/paypal/gatt/examples/option"
)

// Scanner scans for new Devices
type Scanner struct {
	device                      gatt.Device
	onNewDeviceScannedCallbacks []func()
}

// New returns a new Scanner or an error
func New() (*Scanner, error) {
	var onNewDeviceScannedCallbacks []func()

	device, err := gatt.NewDevice(BleOptions.DefaultClientOptions...)
	if err != nil {
		return nil, err
	}

	return &Scanner{device, onNewDeviceScannedCallbacks}, nil
}

// OnNewDeviceScanned registers a function that is
// called whenever a new device is scanned
func (scanner *Scanner) OnNewDeviceScanned(onNewDeviceScannedCallback func()) {
	scanner.onNewDeviceScannedCallbacks = append(scanner.onNewDeviceScannedCallbacks, onNewDeviceScannedCallback)
}

// Scan starts the scanning
func (scanner *Scanner) Scan() {
	device := scanner.device
	device.Handle(gatt.PeripheralDiscovered(scanner.onPeripheralDiscovered))
	device.Init(scanner.onStateChanged)
}

func (scanner *Scanner) onPeripheralDiscovered(peripheral gatt.Peripheral, advertisement *gatt.Advertisement, rssi int) {
	for _, callback := range scanner.onNewDeviceScannedCallbacks {
		callback()
	}
	// fmt.Println("onPeripheralDiscovered")
	// fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", peripheral.ID(), peripheral.Name())
	// fmt.Println("  Local Name        =", advertisement.LocalName)
	// fmt.Println("  TX Power Level    =", advertisement.TxPowerLevel)
	// fmt.Println("  Manufacturer Data =", advertisement.ManufacturerData)
	// fmt.Println("  Service Data      =", advertisement.ServiceData)
}

func (scanner *Scanner) onStateChanged(device gatt.Device, state gatt.State) {
	// fmt.Println("onStateChanged: ", state)

	switch state {
	case gatt.StatePoweredOn:
		// fmt.Println("Scanning...")
		device.Scan([]gatt.UUID{}, false)
		return
	default:
		device.StopScanning()
	}
}
