package service

import (
	"chat_controller_go/repository"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Service struct {
	repository *repository.Repository

	AvailableServers map[string]bool
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository, AvailableServers: make(map[string]bool)}

	s.refreshAvailableServers()

	if err := s.repository.Kafka.RegisterSubTopic("chat"); err != nil {
		panic(err)
	} else {
		go s.loopSubKafka()
	}

	return s
}

func (s *Service) loopSubKafka() {
	for {
		ev := s.repository.Kafka.Poll(100)

		switch event := ev.(type) {
		case *kafka.Message:
			type ServerInfoEvent struct {
				IP     string
				Status bool
			}

			var decoder ServerInfoEvent

			if err := json.Unmarshal(event.Value, &decoder); err != nil {
				log.Println("Failed To Decode Event", event.Value)
			} else {
				fmt.Println(decoder)
				s.AvailableServers[decoder.IP] = decoder.Status
			}

			fmt.Println(event)
		case *kafka.Error:
			log.Println("Failed to polling Event", event.Error())
		}
	}
}

func (s *Service) refreshAvailableServers() {
	servers, err := s.repository.GetAvailableServers()
	if err != nil {
		panic(err)
	}

	for _, server := range servers {
		s.AvailableServers[server.IP] = true
	}
}

func (s *Service) GetAvailableServerIPs() []string {
	servers, err := s.repository.GetAvailableServers()
	if err != nil {
		return []string{}
	}

	ips := make([]string, 0, len(servers)) // Preallocate for efficiency
	for _, s := range servers {
		ips = append(ips, s.IP)
	}
	return ips
}
