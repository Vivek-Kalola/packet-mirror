package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"packet-mirror/utils"
	"packet-mirror/worker"
	"strconv"
	"time"
)

var _logger = utils.NewLogger("Bootstrap")

func main() {

	// Read Config
	configFile := "./config.json"
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		panic("config file not found")
	}

	buffer, err := os.ReadFile(configFile)

	if err == nil && len(buffer) > 0 {

		config := make(map[string]interface{})
		err = json.Unmarshal(buffer, &config)

		if err != nil {
			panic(err)
		}

		if configurations, found := config["configurations"]; found {

			printInterval := int(config["print.interval.sec"].(float64))

			for _, c := range configurations.([]interface{}) {
				mirror(c.(map[string]interface{}), printInterval)
			}
		}

		select {}

	} else {
		panic("Config file is empty")
	}
}

func mirror(config map[string]interface{}, interval int) {

	srcPort := strconv.Itoa(int(config["src.port"].(float64)))

	// Open a socket to receive incoming packets
	connection, err := net.ListenPacket(config["protocol"].(string), ":"+srcPort)

	if err != nil {
		_logger.Error(err.Error())
		return
	}

	// stats
	totalPacket := 0

	var workers []*worker.Worker

	for _, _mirror := range config["mirrors"].([]interface{}) {
		w := worker.New(_mirror.(map[string]interface{}), connection)
		_logger.Info(fmt.Sprintf("%s --> %s", connection.LocalAddr().String(), w.ToString()))
		workers = append(workers, w)
	}

	// Receive packets from the socket and send them to the destination addresses
	go func() {

		for {
			buf := make([]byte, 5*1024*1024)
			n, _, err := connection.ReadFrom(buf)
			if err != nil {
				_logger.Error(err.Error())
				continue
			}

			totalPacket++
			for _, _worker := range workers {
				_worker.Packets <- buf[:n]
			}
		}
	}()

	// stats
	go func() {

		ticker := time.NewTicker(time.Second * time.Duration(interval))

		for {
			select {

			case <-ticker.C:

				_logger.Debug(connection.LocalAddr().String() + ": " + strconv.Itoa(totalPacket))

				totalPacket = 0
			}
		}
	}()
}
