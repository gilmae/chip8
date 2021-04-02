package chip8

type renderer interface {
	Close() error
	Render(d *display) error
}

type nullRenderer struct{}

func NewNullRenderer() *nullRenderer {
	return &nullRenderer{}
}

func (n *nullRenderer) Close() {

}

func (n *nullRenderer) Render(d *display) error {
	return nil
}
