package main

func JoinStringSlices(params ...[]string) []string {
	var result = []string{}
	for _, p := range params {
		for _, v := range p {
			result = append(result, v)
		}
	}
	return result
}

func GetValue(value interface{}, err interface{}) interface{} {
	if err != nil {
        panic(err)
    }
	return value
}