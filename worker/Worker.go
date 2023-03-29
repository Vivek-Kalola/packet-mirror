package worker

import (
	"fmt"
	"net"
	"packet-mirror/utils"
	"time"
)

type Worker struct {
	Packets chan []byte
	dstIP   string
	dstPort int
}

func New(config map[string]interface{}, connection net.PacketConn) *Worker {

	packets := make(chan []byte, 1024*1024)

	totalPacket := 0

	dstIP := config["dst.ip"].(string)

	dstPort := int(config["dst.port"].(float64))

	interval := int(config["print.interval.sec"].(float64))

	_logger := utils.NewLogger(dstIP)

	destAddress := &net.UDPAddr{IP: net.ParseIP(dstIP), Port: dstPort}

	_logger.Info(connection.LocalAddr().String() + " --> " + destAddress.String())

	go func() {

		for {

			select {

			case packet := <-packets:

				totalPacket++

				_, err := connection.WriteTo(packet, destAddress)

				if err != nil {
					_logger.Error(err.Error())
				}
			}
		}
	}()

	go func() {

		ticker := time.NewTicker(time.Second * time.Duration(interval))

		for {
			select {

			case <-ticker.C:

				_logger.Debug(fmt.Sprintf("%s --> %s : %d", connection.LocalAddr().String(), destAddress.String(), totalPacket))

				totalPacket = 0
			}
		}
	}()

	return &Worker{
		Packets: packets,
		dstIP:   dstIP,
		dstPort: dstPort,
	}
}

func (worker *Worker) ToString() string {
	return fmt.Sprintf("%s:%d", worker.dstIP, worker.dstPort)
}
