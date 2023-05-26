package tools

type Md5 struct {
}

// @provider({"id":"md5"})
// @inject({"param":"name", "id":"name"})
// @inject({"param":"class", "id":"class"})
// @inject({"param":"tools", "id":"tools.tools"})
// @inject({"param":"tools2", "id":"tools.tools2"})
func NewMd5(name, class string, tools Tools, tools2 Tools) *Md5 {
	return &Md5{}
}
