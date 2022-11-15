package main

import "emperror.dev/errors"

func checkAndFillVecMatrix(mat map[string][]string) error {
	if len(mat) == 0 {
		return errors.NewPlain("*_vec metric should specify matrix")
	}
	for k, v := range mat {
		if len(v) == 0 {
			mat[k] = gPredefinedMatrix[k]
			if len(mat[k]) == 0 {
				return errors.Errorf("matrix key %s does not contains value")
			}
		}
	}
	return nil
}
