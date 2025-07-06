package stackable

import "iter"

type HeadersContainer struct {
    headers map[string][]string
}

func NewHeadersContainer() HeadersContainer {
    return HeadersContainer{
        headers: make(map[string][]string),
    }
}

func (s *HeadersContainer) Contains(key string) bool {
    _, present := s.headers[key]

    return present
}

func (s *HeadersContainer) Set(key string, value string) {
    s.headers[key] = []string{value}
}

func (s *HeadersContainer) Add(key string, value string) {
    if !s.Contains(key) {
        s.Set(key, value)
        return
    }

    s.headers[key] = append(s.headers[key], value)
}

func (s *HeadersContainer) Get(key string) []string {
    return s.headers[key]
}

func (s *HeadersContainer) Delete(key string) []string {
    deleted := s.headers[key]

    delete(s.headers, key)

    return deleted
}

func (s *HeadersContainer) Entries() iter.Seq2[string, []string] {
    return func (yield func(string, []string) bool) {
        for k, v := range s.headers {
            if !yield(k, v) {
                return
            }
        }
    }
}
