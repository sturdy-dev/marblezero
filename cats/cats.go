package cats

import _ "embed"

//go:embed cat_1.txt
var CatDefault string

//go:embed cat_2.txt
var CatDefaultLookRight string

//go:embed cat_3.txt
var Cat3 string

//go:embed cat_4.txt
var CatAmused string
