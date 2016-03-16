package scanner

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/paypal/gatt"
	BleOptions "github.com/paypal/gatt/examples/option"
)

// Scanner scans for new Devices
type Scanner struct {
	device                      gatt.Device
	filter                      string
	onErrorCallbacks            []OnErrorCallback
	onNewDeviceScannedCallbacks []OnNewDeviceScannedCallback
	scannedDevices              map[string]*Device
}

// OnErrorCallback is called with any errors
// that happen on the scanner
type OnErrorCallback func(error)

// OnNewDeviceScannedCallback is called with the id of the
// device that has been scanned
type OnNewDeviceScannedCallback func(device *Device)

// New returns a new Scanner or an error. Filter is a
// hex string that is used to filter peripherals
// by ManufacturerData
func New(filter string) (*Scanner, error) {
	var onErrorCallbacks []OnErrorCallback
	var onNewDeviceScannedCallbacks []OnNewDeviceScannedCallback
	scannedDevices := make(map[string]*Device)

	device, err := gatt.NewDevice(BleOptions.DefaultClientOptions...)
	if err != nil {
		return nil, err
	}

	return &Scanner{device, filter, onErrorCallbacks, onNewDeviceScannedCallbacks, scannedDevices}, nil
}

// OnError calls the callback with an error whenever it happens.
// if no error callbacks are registered, the scanner will panic
func (scanner *Scanner) OnError(callback OnErrorCallback) {
	scanner.onErrorCallbacks = append(scanner.onErrorCallbacks, callback)
}

// OnNewDeviceScanned registers a function that is
// called whenever a new device is scanned
func (scanner *Scanner) OnNewDeviceScanned(callback OnNewDeviceScannedCallback) {
	scanner.onNewDeviceScannedCallbacks = append(scanner.onNewDeviceScannedCallbacks, callback)
}

// Scan starts the scanning
func (scanner *Scanner) Scan() {
	device := scanner.device
	device.Handle(gatt.PeripheralDiscovered(scanner.onPeripheralDiscovered))
	device.Init(scanner.onStateChanged)
}

func (scanner *Scanner) emitError(err error) {
	if len(scanner.onErrorCallbacks) == 0 {
		log.Panicln("No error callbacks registered, but error occured: ", err.Error())
	}

	for _, callback := range scanner.onErrorCallbacks {
		callback(err)
	}
}

func (scanner *Scanner) deviceMatchesFilter(data []byte) bool {
	hexData := hex.EncodeToString(data)
	filter := strings.ToLower(scanner.filter)
	return strings.Contains(hexData, filter)
}

func (scanner *Scanner) onPeripheralDiscovered(peripheral gatt.Peripheral, advertisement *gatt.Advertisement, rssi int) {
	if !scanner.deviceMatchesFilter(advertisement.ManufacturerData) {
		return
	}

	id := peripheral.ID()
	if device, ok := scanner.scannedDevices[id]; ok {
		device.AddRSSI(rssi)
		return
	}

	device := NewDevice(advertisement, peripheral, rssi)
	scanner.scannedDevices[id] = device

	for _, callback := range scanner.onNewDeviceScannedCallbacks {
		callback(device)
	}
}

func (scanner *Scanner) onStateChanged(device gatt.Device, state gatt.State) {
	switch state {
	case gatt.StatePoweredOn:
		// fmt.Println("Scanning...")
		device.Scan([]gatt.UUID{}, false)
		return
	case gatt.StatePoweredOff:
		scanner.emitError(fmt.Errorf("Bluetooth is powered off."))
		return
	default:
		device.StopScanning()
	}
}
