package atlas

import (
	"fmt"
	"net/url"
)

type GroupWhiteListService interface {
	Create(w GroupWhiteList) (*GroupWhiteList, *Response, error)
	Delete(groupID string, identifier string) (*Response, error)
	List(groupID string) ([]GroupWhiteList, *Response, error)
	Get(groupID string, identifier string) (*GroupWhiteList, *Response, error)
}

type GroupWhiteListServiceOp struct {
	client *Client
}

var _ = &GroupWhiteListServiceOp{}

type GroupWhiteLists []GroupWhiteList

type GroupWhiteList struct {
	Links     []Link `json:"links,omitempty"`
	CidrBlock string `json:"cidrBlock"`
	GroupID   string `json:"groupId,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
}

type GroupWhiteListCollection struct {
	Results []GroupWhiteList
}

type Link struct {
	Href string `json:"href,omitempty"`
	Rel  string `json:"rel,omitempty"`
}

func (s *GroupWhiteListServiceOp) List(groupID string) ([]GroupWhiteList, *Response, error) {
	path := fmt.Sprintf("groups/%s/whitelist", groupID)
	req, err := s.client.NewRequest("GET", path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(GroupWhiteListCollection)
	resp, err := s.client.Do(req, root)

	if err != nil {
		return nil, resp, err
	}

	return root.Results, resp, err
}

func (s *GroupWhiteListServiceOp) Get(groupID string, identifier string) (*GroupWhiteList, *Response, error) {
	path := fmt.Sprintf("groups/%s/whitelist/%s", groupID, url.QueryEscape(identifier))
	req, err := s.client.NewRequest("GET", path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(GroupWhiteList)
	resp, err := s.client.Do(req, root)

	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *GroupWhiteListServiceOp) Create(w GroupWhiteList) (*GroupWhiteList, *Response, error) {
	reqBody := &GroupWhiteLists{w}
	path := fmt.Sprintf("groups/%s/whitelist", w.GroupID)
	req, err := s.client.NewRequest("POST", path, reqBody)

	if err != nil {
		return nil, nil, err
	}

	root := new(GroupWhiteList)
	resp, err := s.client.Do(req, root)

	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *GroupWhiteListServiceOp) Delete(groupID string, identifier string) (*Response, error) {
	path := fmt.Sprintf("groups/%s/whitelist/%s", groupID, url.QueryEscape(identifier))
	req, err := s.client.NewRequest("DELETE", path, nil)

	if err != nil {
		return nil, err
	}

	root := new(GroupWhiteList)
	resp, err := s.client.Do(req, root)

	if err != nil {
		return resp, err
	}

	return resp, err
}
