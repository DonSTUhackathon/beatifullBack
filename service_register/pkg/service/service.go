package service

import (
	"errors"
	"log"
	"os/exec"
	"runtime"
)

type Service struct {
	url string
}

func (s *Service) GetUrl() string {
	return s.url
}

func (s *Service) SetUrl(url string) {
	s.url = url
}

func (s *Service) RunApp() error {
	var cmd string

	ops := runtime.GOOS
	switch ops {
	case "windows":
		cmd = "start "
	case "darwin":
		cmd = "open "
	case "linux":
		cmd = "xdg-open "
	default:
		return errors.New("error in OS switch")
	}

	log.Print(cmd, s.url)
	c := exec.Command("cmd", "/C", cmd, s.url)
	err := c.Start()
	if err != nil {
		return err
	}

	return nil
}

func NewService(ipAddr, port string) *Service {
	return &Service{"http://" + ipAddr + ":" + port + "/auth"}
}
