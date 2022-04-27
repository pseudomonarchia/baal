package controller

import "baal/service"

// Controllers represents a global controllers struct
type Controllers struct {
	Index *Index
	OAuth *OAuth
}

// New return all controller
func New(s *service.Services) *Controllers {
	return &Controllers{
		Index: &Index{s},
		OAuth: &OAuth{s},
	}
}