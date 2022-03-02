package source

import "github.com/andriiyaremenko/mg/dto"

type Source func() (*dto.Target, bool, error)
