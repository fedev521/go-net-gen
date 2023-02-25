package app

type Drawer interface {
	Draw() error
	Render() error
}
