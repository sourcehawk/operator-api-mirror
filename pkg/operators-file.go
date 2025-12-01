package pkg

type OperatorsFile struct {
	Operators []OperatorConfig `yaml:"operators"`
}

func (op *OperatorsFile) Process(mirrorRootPath string, moduleRoot string, target string) error {
	for _, op := range op.Operators {
		if target == "" || target == op.Slug {
			err := op.Mirror(mirrorRootPath, moduleRoot)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
