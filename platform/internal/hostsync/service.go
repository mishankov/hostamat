package hostsync

import (
	"context"
	"slices"
	"sync"
)

type webserver interface {
	GetHostsFromShorthandId(ctx context.Context, shorthandId string) ([]string, error)
	RewriteHostsInShorthandId(ctx context.Context, shorthandId string, hosts []string) error
}

type Service struct {
	mu          sync.Mutex
	webserver   webserver
	shorthandId string
}

func NewService(webserver webserver, shorthandId string) *Service {
	return &Service{
		webserver:   webserver,
		shorthandId: shorthandId,
	}
}

func (s *Service) Add(ctx context.Context, host string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	hosts, err := s.webserver.GetHostsFromShorthandId(ctx, s.shorthandId)
	if err != nil {
		return err
	}

	if slices.Contains(hosts, host) {
		return nil
	}

	hosts = append(hosts, host)

	return s.webserver.RewriteHostsInShorthandId(ctx, s.shorthandId, canonicalHosts(hosts))
}

func (s *Service) Remove(ctx context.Context, host string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	hosts, err := s.webserver.GetHostsFromShorthandId(ctx, s.shorthandId)
	if err != nil {
		return err
	}

	filtered := hosts[:0]
	for _, current := range hosts {
		if current == host {
			continue
		}

		filtered = append(filtered, current)
	}

	return s.webserver.RewriteHostsInShorthandId(ctx, s.shorthandId, canonicalHosts(filtered))
}

func (s *Service) Sync(ctx context.Context, hosts []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.webserver.RewriteHostsInShorthandId(ctx, s.shorthandId, canonicalHosts(hosts))
}

func canonicalHosts(hosts []string) []string {
	if len(hosts) == 0 {
		return []string{}
	}

	uniq := make(map[string]struct{}, len(hosts))
	result := make([]string, 0, len(hosts))
	for _, host := range hosts {
		if host == "" {
			continue
		}

		if _, ok := uniq[host]; ok {
			continue
		}

		uniq[host] = struct{}{}
		result = append(result, host)
	}

	slices.Sort(result)

	return result
}
