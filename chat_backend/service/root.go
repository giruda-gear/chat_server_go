package service

import (
	"chat_server_go/repository"
	"chat_server_go/types/schema"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Service struct {
	repository *repository.Repository
}

func Newservice(repository *repository.Repository) *Service {
	s := &Service{repository: repository}

	return s
}

func (s *Service) PublishServerStatusEvent(ip string, status bool) {
	type ServerInfoEvent struct {
		IP     string
		Status bool
	}

	e := &ServerInfoEvent{IP: ip, Status: status}
	ch := make(chan kafka.Event)

	if v, err := json.Marshal(e); err != nil {
		log.Println("Failed to marshall")
	} else if result, err := s.PublishEvent("chat", v, ch); err != nil {
		// TODO Send event to Kafka
		log.Println("Failed to send send event to Kafka err:", err)
	} else {
		log.Println("Success to send send event to Kafka", result)
	}
}

func (s *Service) PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	return s.repository.Kafka.PublishEvent(topic, value, ch)
}

func (s *Service) ServerSet(ip string, available bool) error {
	if err := s.repository.ServerSet(ip, available); err != nil {
		log.Println("Failed to serverset", "ip", ip, "available", available)
		return err
	}
	return nil
}

func (s *Service) InsertChatting(user, message, roomName string) {
	if err := s.repository.InsertChatting(user, message, roomName); err != nil {
		log.Println("Failed to chat..", "err:", err)
	}
}

func (s *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	if res, err := s.repository.GetChatList(roomName); err != nil {
		log.Println("Failed to get chat list", "err:", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) RoomList() ([]*schema.Room, error) {
	if res, err := s.repository.RoomList(); err != nil {
		log.Println("Failed to get all room list", "err:", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) MakeRoom(name string) error {
	if err := s.repository.MakeRoom(name); err != nil {
		log.Println("Failed to make a room", "err:", err.Error())
		return err
	} else {
		return nil
	}
}

func (s *Service) Room(name string) (*schema.Room, error) {
	if res, err := s.repository.Room(name); err != nil {
		log.Println("Failed to get a room", "err:", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}
