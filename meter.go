package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"sync"
	"time"
)

func init() {
	Meter = &meter{}
	Meter.measurementsAvg10Seconds = make([]int64, 10, 10)
	Meter.measurementsAvgMinute = make([]int64, 60, 60)
	Meter.mu = new(sync.RWMutex)
	Meter.Run()
}

var Meter *meter

type meter struct {
	bytesPer10Seconds int64
	bytesPerMinute    int64

	liveBytesCounter         int64
	measurementsAvg10Seconds []int64
	measurementsAvgMinute    []int64
	mu                       *sync.RWMutex
}

func (self *meter) GetHumanReadablePer10Seconds() string {
	bytes := self.GetBytesPer10Seconds()
	return self.bytesToHumanReadable(float64(bytes))
}

func (self *meter) GetHumanReadablePerMinute() string {
	bytes := self.GetBytesPerMinute()
	return self.bytesToHumanReadable(float64(bytes))
}

func (self *meter) GetBytesPer10Seconds() int64 {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.bytesPer10Seconds
}

func (self *meter) GetBytesPerMinute() int64 {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return self.bytesPerMinute
}

func (self *meter) RegisterBytesWritten(written int64) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.liveBytesCounter += written
}

func (self *meter) Run() {
	// Take down the stats every second
	go func() {
		for {
			waiter := time.After(1 * time.Second)
			<-waiter
			self.mu.Lock()
			self.measurementsAvg10Seconds = append([]int64{self.liveBytesCounter}, self.measurementsAvg10Seconds[:len(self.measurementsAvg10Seconds)-1]...)
			self.measurementsAvgMinute = append([]int64{self.liveBytesCounter}, self.measurementsAvgMinute[:len(self.measurementsAvgMinute)-1]...)
			self.liveBytesCounter = int64(0)
			self.mu.Unlock()
		}
	}()
	// Update output fields at the specific intervals
	go func() {
		for {
			waiter := time.After(10 * time.Second)
			<-waiter
			self.mu.Lock()
			avg := int64(0)
			for _, bytes := range self.measurementsAvg10Seconds {
				avg += bytes
			}
			self.bytesPer10Seconds = int64(avg / 10)
			self.mu.Unlock()
		}
	}()
	go func() {
		for {
			waiter := time.After(1 * time.Minute)
			<-waiter
			self.mu.Lock()
			avg := int64(0)
			for _, bytes := range self.measurementsAvgMinute {
				avg += int64(bytes)
			}
			self.bytesPerMinute = int64(avg / 60)
			self.mu.Unlock()
		}
	}()
}

func (self *meter) bytesToHumanReadable(bytes float64) string {
	if bytes < math.Pow(1024, 1) {
		return fmt.Sprintf("%.f bytes", bytes)
	} else if bytes < math.Pow(1024, 2) {
		return fmt.Sprintf("%.2f KiB", bytes/math.Pow(1024, 1))
	} else if bytes < math.Pow(1024, 3) {
		return fmt.Sprintf("%.2f MiB", bytes/math.Pow(1024, 2))
	} else if bytes < math.Pow(1024, 4) {
		return fmt.Sprintf("%.2f GiB", bytes/math.Pow(1024, 3))
	}
	return fmt.Sprintf("%.2f TiB", bytes/math.Pow(1024, 4))
}

func CopyAndMeasureThroughput(writer, reader net.Conn) {
	var err error
	written := int64(0)
	buffer := make([]byte, 32*1024)
	for {
		bytesRead, readErr := reader.Read(buffer)
		if bytesRead > 0 {
			bytesWritten, writeErr := writer.Write(buffer[0:bytesRead])
			if bytesWritten > 0 {
				written += int64(bytesWritten)
			}
			if writeErr != nil {
				err = writeErr
				break
			}
			if bytesRead != bytesWritten {
				err = io.ErrShortWrite
				break
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			err = readErr
			break
		}
		Meter.RegisterBytesWritten(written)
		written = 0
	}
	if err != nil {
		log.Println("Tunneling connection produced an error", err)
		return
	}
	Meter.RegisterBytesWritten(written)
	written = 0
}
