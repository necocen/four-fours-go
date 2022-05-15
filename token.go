package fourfours

// 演算子および桁トークン。
// `0xe0`より後の値は桁のために予約済みで、
// `t = 0xe0 + n (0 <= n < 10)`
// または
// `t = 0xf0 + n (0 <= n < 10)`
//
// のとき、数字`n`を表す。
// `0xeX`は先頭桁を示し、`0xfX`は残りの桁を示す。
type OperatorToken uint8
