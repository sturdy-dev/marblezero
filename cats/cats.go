package cats

import _ "embed"

//go:embed cat_normal_straight.txt
var CatNormalStraight string

//go:embed cat_normal_straight_raised_tail.txt
var CatNormalStraightRaisedTail string

//go:embed cat_normal_straight_folded_left_ear.txt
var CatNormalStraightFoldedLeftEar string

//go:embed cat_normal_straight_folded_right_ear.txt
var CatNormalStraightFoldedRightEar string

//go:embed cat_normal_left.txt
var CatNormalLeft string

//go:embed cat_normal_right.txt
var CatNormalRight string

//go:embed cat_amused.txt
var CatAmused string

//go:embed cat_curious.txt
var CatCurious string
