package providers

type DocIdStubProvider struct{}

func (p *DocIdStubProvider) Provide(dir string) (int, error) {
	return len(dir), nil
}
