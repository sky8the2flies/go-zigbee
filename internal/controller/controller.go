package controller

type Controller struct {
	Port Port
}

func NewController(port Port) *Controller {
	return &Controller{Port: port}
}

func (c *Controller) Start() error {
	err := c.Port.Open()
	if err != nil {
		return err
	}
	return nil
}
